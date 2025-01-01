package main

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	//"net/http"
)

// Define a handler function for root endpoint
func okHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	// Send an "OK" message
	fmt.Fprint(w, "OK")
}
func main() {
	// Register handler functions for routes
	http.HandleFunc("/", okHandler)

	// Set the address and port to listen on
	port := ":8080"
	fmt.Println("Starting server on port:", port)

	// Start the server and log errors if any
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
