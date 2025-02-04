package main

import (	
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	db := ConfigureDB()
	//defer db.Close()
	//SeedDB(db)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", Hello).Methods(http.MethodGet)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Use(APImiddleware)
	s := r.PathPrefix("/v1").Subrouter()
	s.Use(loggingMiddleware)
	
	//// Authentication Routes
	// Login
	r.HandleFunc("/authentication/login", Login(db)).Methods(http.MethodPost)

	//// Routes with active "session"
	s.HandleFunc("/transaction", getTransactionsRequest(db)).Methods(http.MethodGet)
	s.HandleFunc("/transaction", createTransactionsRequest(db)).Methods(http.MethodPost)

	fmt.Println("The server is on: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
