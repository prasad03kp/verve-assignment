package main

import (
	"github.com/gorilla/mux"
	utilities "github.com/prasad03kp/verve-assignment/utilities"
	verve "github.com/prasad03kp/verve-assignment/verve"
	"net/http"
	"os"
	"log"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		utilities.GetVersion(w, r)
	}).Methods("GET")

	r.HandleFunc("/api/verve/accept", func(w http.ResponseWriter, r *http.Request) {
		verve.Accept(w, r)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	log.Println("Server started on port " + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println("Cannot start server:", err)
	}
}

