package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"url-shortener/src/database"
	"url-shortener/src/links"
)

func main() {
	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/short", links.CreateLink).Methods("POST")
	r.HandleFunc("/{code}", links.RedirectLink).Methods("GET")

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println("failed to start server")
	}
}
