package routers

import (
	"car-comparison-service/service/handlers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/ping", handlers.PingHandler()).Methods("GET")

	return router
}
