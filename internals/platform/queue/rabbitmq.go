package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

type QueueStore struct{
	*amqp.Connection;
}

func Init_queue(rabbitURL string)(*QueueStore, error){
	conn , err := amqp.Dial(rabbitURL);
	if(err != nil){
		return nil, fmt.Errorf("error while initilization of queue");
	}
	return &QueueStore{conn}, nil;
}
