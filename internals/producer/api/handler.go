package api

import (
	"djq/internals/producer"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)


type API struct{
	producer *producer.Service;
}

type CreateJobRequest struct{
	JobType string `json:"type"`
	JobPayload  json.RawMessage `json:"payload"`
}

func Init_API( producer *producer.Service)(*API){
	return &API{producer: producer};
}

func (apiHandler *API) HandleCreateJob(w http.ResponseWriter , r *http.Request){
	if(r.Method != http.MethodPost){
		http.Error(w , "Method not allowed" , http.StatusMethodNotAllowed);
	}

	body,err := io.ReadAll(r.Body);
	if(err != nil){
		log.Println("error while parsing the body: " , err);
		http.Error(w , "failed reading the body" , http.StatusInternalServerError);
	}
	var res  CreateJobRequest; 
	if err := json.Unmarshal(body , &res); err != nil{
		fmt.Println(res);
		log.Println("error while parsing the body: " , err);
		http.Error(w , "failed reading the body" , http.StatusInternalServerError);
	}
	fmt.Println(res.JobPayload)
	job, err := apiHandler.producer.CreateJob(res.JobType, res.JobPayload);
	if(err != nil){
		log.Println("Error while creating the job: " , err);
		http.Error(w , "Failed to create the job" , http.StatusInternalServerError);
	}
	w.Header().Set("Content-type" , "application/json");
	w.WriteHeader(http.StatusAccepted);
	if err := json.NewEncoder(w).Encode(job); err != nil{
		log.Println("Failed to write response: " , err);
	}
}
