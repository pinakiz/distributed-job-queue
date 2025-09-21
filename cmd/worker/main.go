package worker

import (
	"djq/internals/platform/database"
	"djq/internals/platform/queue"
	"djq/internals/worker"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// db initilize
// queue initilize
func main(){
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

	broker , err:= queue.Init_queue(queue_url);

	if(err != nil){
		fmt.Println(err);
	}

	workerService := worker.Init_worker(dbStore, broker);

	go workerService.Run()	

	quit := make(chan os.Signal, 1);
	signal.Notify(quit , syscall.SIGINT , syscall.SIGTERM);

	<- quit;
	log.Println("Shutting down the worker service ...");
}