package database

import (
	"djq/internals/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBstore struct{
	Store *sqlx.DB
}


func Init_db(postgresDSN string) (*DBstore , error){
	db , err := sqlx.Open("postgres",postgresDSN);
	if(err != nil){
		return nil, fmt.Errorf("error while initiating DB");
	}else{
		return &DBstore{db}, nil;
	}
}


func (db *DBstore) InsertJob(jobType string,jobPayload json.RawMessage)(*models.Job , error){
	var job models.Job;
	job.Id = uuid.NewString(); 
	job.Type = jobType;
	job.Payload = string(jobPayload);
	job.Status = models.StatusQueued;
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt;

	query := `Insert INTO jobs (id, type, status, payload, created_at , updated_at) 
			Values ($1 , $2 , $3 , $4 , $5 , %6)`
	_,err := db.Store.Exec(query,job.Id,job.Type,job.Status,job.Payload,job.CreatedAt,job.UpdatedAt)
	if(err != nil){
		log.Println("error while inserting job into table");
		return nil, err;
	}
	return &job,err;	
}

func (db *DBstore) GetJob(id string) (*models.Job , error){
	query := `SELECT * FROM jobs WHERE id=$1`;
	var job models.Job
	err := db.Store.Get(&job, query,id);

	if(err != nil){
		log.Println("error while getting job: ",err);
		return nil, err;
	}
	return &job,err;

}

func (db *DBstore) UpdateJobStatus(id string , status string)(error){
	query := `UPDATE jobs SET status=$1 , updated_at = $2  WHERE id = $3`;
	_ , err := db.Store.Exec(query , status , time.Now().UTC() , id);
	return err;
}