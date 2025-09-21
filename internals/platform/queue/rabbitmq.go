package queue

import (
	"djq/internals/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type QueueStore struct{
	Channel *amqp.Channel
	queueName string
}

func Init_queue(rabbitURL string)(*QueueStore, error){
	// fmt.Println(rabbitURL)
	conn , err := amqp.Dial(rabbitURL);
	
	if(err != nil){
		return nil, fmt.Errorf("error while initilization of queue: %w" , err);
	}
	ch , err := conn.Channel();
	
	if(err != nil){
		log.Fatal("error while declaring queue ",err);
	}
	_ , err = ch.QueueDeclare(
		"jobQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil{
		log.Println("error while getting the queue channel" , err);
		return nil , fmt.Errorf("error while getting the queue channel: %w" , err);
	}
	return &QueueStore{ch , "jobQueue"}, nil;
}

func (broker *QueueStore) Publish (job *models.Job){
	ch:= broker.Channel;
	
	body , err := json.Marshal(job);
	if(err != nil){
		log.Println("Error: Failed to marshal job")
	}
	if err := ch.Publish(
		"",
		broker.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
		},
	);err != nil{
		log.Println("error while pushing data into queue ", err);
	}	
}

func (broker *QueueStore) Consume()(<-chan amqp.Delivery , error){
	err := broker.Channel.Qos(
		1,
		0,
		false,
	)
	if err != nil{
		log.Println("error while Qos: " , err);
		return nil, fmt.Errorf("error while Qos: %w" , err);
	}

	msgs , err := broker.Channel.Consume(
		broker.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if  err != nil {
		log.Println("error while fetching msg from the queue: " , err);
		return nil, fmt.Errorf("error while fetching msg from the queue: %w" , err);
	}
	return msgs, nil;
}

func (broker * QueueStore) Close(){
	if(broker.Channel!=nil){
		broker.Channel.Close()
	}
}