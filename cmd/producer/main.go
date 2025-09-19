package main

import (
	"djq/internals/platform/database"
	"djq/internals/platform/queue"
	"djq/internals/producer"
	"djq/internals/producer/api"
	"fmt"
	"log"
	"net/http"
	"os"
)

// db initilize
// queue initilize
// router and handler

func main(){

	listenAddr := ":8080";

	db_host := os.Getenv("DBHOST")
	db_port := os.Getenv("DBPORT")
	db_user := os.Getenv("DBUSER")
	db_password := os.Getenv("DBPASSWORD")
	db_name := os.Getenv("DBNAME")

	postgresDSN := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    db_host, db_port, db_user, db_password, db_name)

	dbStore , err := database.Init_db(postgresDSN);
 	if(err != nil){
		fmt.Println(err);
		return
	}
	queue_url := os.Getenv("RABBIT_URL");

	broker , err:= queue.Init_queue(queue_url);

	if(err != nil){
		fmt.Println(err);
	}

	producerService := producer.Producer(dbStore , broker);

	apiHandler := api.Init_API(producerService);
		
	router := api.NewRouter(apiHandler);

	log.Println("Server listening on: ", listenAddr);
	if err := http.ListenAndServe(listenAddr,router); err != nil{
		log.Fatal("Server Crash: ",err);
		return;
	}
}