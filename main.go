package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt-processor/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Create new router
	router := mux.NewRouter()

	// Register handler functions for routes
	router.HandleFunc("/receipts/process", handlers.PostReceiptProccessHandler)
	router.HandleFunc("/receipts/{id}/points", handlers.GetReceiptProccesHandler).Methods("GET")

	// Set the address and port to listen on
	port := ":8080"
	fmt.Println("Starting server on port:", port)

	// Start the server and log errors if any
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
