package api

import (
	"github.com/gorilla/mux"
)

func NewRouter(api *API)(*mux.Router){
	r := mux.NewRouter();
	r.HandleFunc("/job" , api.HandleCreateJob);


	return r;
}