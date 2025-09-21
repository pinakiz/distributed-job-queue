package producer

import (
	"djq/internals/models"
	"djq/internals/platform/database"
	"djq/internals/platform/queue"
	"encoding/json"
	"log"
)

type Service struct{
	db *database.DBstore
	broker *queue.QueueStore
}


func Producer(db *database.DBstore , rabbit *queue.QueueStore)(*Service){
	return &Service{db: db , broker: rabbit};
}

func (s *Service) CreateJob(jobType string, jobPayload json.RawMessage)(*models.Job , error){
	db := s.db;
	job , err := db.InsertJob(jobType , jobPayload);
	if(err != nil){
		log.Println("error while job creation: ",err);
		return &models.Job{}, err;
	}
	s.broker.Publish(job);
	return job , nil;
}