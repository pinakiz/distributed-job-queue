package queue

import (
	"djq/internals/models"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type QueueStore struct{
	Conn *amqp.Connection;
}

func Init_queue(rabbitURL string)(*QueueStore, error){
	conn , err := amqp.Dial(rabbitURL);
	if(err != nil){
		return nil, fmt.Errorf("error while initilization of queue");
	}
	return &QueueStore{conn}, nil;
}

func (broker *QueueStore) Publish (job *models.Job){
	ch, err := broker.Conn.Channel();
	if(err != nil){
		log.Println("error while opening the channel (queue):", err);
	}
	defer ch.Close();

	queue , err := ch.QueueDeclare(
		"jobQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if(err != nil){
		log.Fatal("error while declaring queue ",err);
	}
	if err := ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(job.Payload),
		},
	);err != nil{
		log.Println("error while pushing data into queue ", err);
	}	
}
