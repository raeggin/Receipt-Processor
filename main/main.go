package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Create the Items struct
type Items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Create the Recipt struct
type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        string  `json:"total"`
	Items        []Items `json:"items"`
}

// Create ReceiptResponse struct; this is the response structure for when you create a new receipt
type ReceiptResponse struct {
	Id string `json:"id"`
}

// Define a handler function for root endpoint
func okHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	// Send an "OK" message
	fmt.Fprint(w, "OK")
}

// Define a handler function for POST /receipt/process endpoint
func postReceiptProccessHandler(w http.ResponseWriter, r *http.Request) {
	// Validate for post request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON data from the request body
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Create UUID and assign it to a ReceiptResponse Struct
	id := uuid.New().String()
	var receiptResponse ReceiptResponse
	receiptResponse.Id = id

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(receiptResponse)

}

// Define a handler function for GET /receipts/{id}/points endpoint
func getReceiptProccesHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	//Send message
	fmt.Fprint(w, "This is my GET request: ")
}

func main() {
	// Register handler functions for routes
	http.HandleFunc("/", okHandler)
	http.HandleFunc("/receipts/process", postReceiptProccessHandler)
	http.HandleFunc("/receipts/{id}/points", getReceiptProccesHandler)

	// Set the address and port to listen on
	port := ":8080"
	fmt.Println("Starting server on port:", port)

	// Start the server and log errors if any
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
