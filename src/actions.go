package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/ GET Hello")
	fmt.Fprintln(w, "Hello!")
}

func getTransactionsRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/transaction GET Transactions")
		transactions:=getTransactions(db)

		json.NewEncoder(w).Encode(transactions)
		
	}
}

func createTransactionsRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/Transactions POST Transaction")
		var transaction Transaction
		err := json.NewDecoder(r.Body).Decode(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Transaction values",transaction)
		// Create transaction
		return
		
	}
}

