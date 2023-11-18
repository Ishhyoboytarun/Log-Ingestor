package main

import (
	"Log-Ingestor/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	handler.Init()
	handlers := handler.Handler
	router := mux.NewRouter()
	router.HandleFunc("/injest-logs", handlers.Injest.InjestLogs).Methods(http.MethodGet)

	err := http.ListenAndServe("localhost:3000", router)
	if err != nil {
		panic(err)
	}
}
