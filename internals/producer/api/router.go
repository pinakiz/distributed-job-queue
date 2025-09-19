package api

import (
	"github.com/gorilla/mux"
)



type Router struct{
	router *mux.Router
}

func NewRouter(api *API)(*mux.Router){
	r := mux.NewRouter();
	r.HandleFunc("/job" , api.HandleCreateJob);


	return r;
}