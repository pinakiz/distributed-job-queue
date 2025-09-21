package main

import (
	"djq/internals/platform/database"
	"djq/internals/platform/queue"
	"djq/internals/worker"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

// db initilize
// queue initilize
func main(){
	
	err := godotenv.Load();
	if(err != nil){
		log.Fatal("Error: failed env loading");
	}
	db_host := os.Getenv("DBHOST")
	db_port := os.Getenv("DBPORT")
	db_user := os.Getenv("DBUSER")
	db_password := os.Getenv("DBPASSWORD")
	db_name := os.Getenv("DBNAME")

	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    db_host, db_port, db_user, db_password, db_name)

	dbStore , err := database.Init_db(postgresDSN);
 	if(err != nil){
		fmt.Println(err);
		return
	}
	queue_url := os.Getenv("RABBIT_URL");
	fmt.Println(db_host)
	broker , err:= queue.Init_queue(queue_url);

	if(err != nil){
		fmt.Println(err);
	}

	workerService := worker.Init_worker(dbStore, broker);
	defer broker.Close();
	go workerService.Run()	

	quit := make(chan os.Signal, 1);
	signal.Notify(quit , syscall.SIGINT , syscall.SIGTERM);

	<- quit;
	log.Println("Shutting down the worker service ...");
}