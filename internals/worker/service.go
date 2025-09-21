package worker

import (
	"djq/internals/models"
	"djq/internals/platform/database"
	"djq/internals/platform/queue"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)



type Service struct{
	db *database.DBstore;
	broker *queue.QueueStore;
}



func Init_worker(db *database.DBstore , broker *queue.QueueStore)(*Service){
	return &Service{db: db , broker:  broker};
}

func (s *Service) Run (){
	msg , err := s.broker.Consume();
	if err != nil{
		log.Println("error while consuming data from the queue");
	}
	
	forever := make(chan bool);

	log.Println("Worker is running and waiting for the jobs...");

	for d := range msg{
		go func (del amqp.Delivery){
			var job models.Job;
			if err := json.Unmarshal(del.Body , &job); err != nil{
				log.Printf("Error: couldn't decode job json: %v , Message will be discarded" , err);
				del.Nack(false, false)
				if er := s.db.UpdateJobStatus(job.Id , models.StatusFailed); er != nil{
					log.Printf("Error: Status update failes: %v" , er);
				};
				return;
			}
			log.Printf("Received job: %s. Starting processing..." , job.Id);
			if err := s.db.UpdateJobStatus(job.Id , models.StatusProcessing);err != nil{
					log.Printf("Error: Status update failes: %v" , err);
			}
			
			if err := s.handleJob(&job);err != nil{
				log.Printf("Error: failed to process the job: %v" , job.Id);
				s.db.UpdateJobStatus(job.Id, models.StatusFailed);
			}else{
				log.Printf("Successfully process the job: %v" , job.Id);
				s.db.UpdateJobStatus(job.Id , models.StatusCompleted);
				del.Ack(false);
			}
		}(d)
	}
	<- forever;
}

func (s *Service) handleJob(job *models.Job) error {
		switch job.Type {
	case "send_email":
		log.Printf("Simulating sending email for job %s with payload: %s", job.Id, job.Payload)
		time.Sleep(10 * time.Second)  
	case "generate_report":
		log.Printf("Simulating report generation for job %s...", job.Id)
		time.Sleep(15 * time.Second) 
	default:
		log.Printf("WARNING: Unknown job type '%s' for job %s. Marking as failed.", job.Type, job.Id)
	}
	return nil

}