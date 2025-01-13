package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/points"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Define a handler function for POST /receipt/process endpoint
func postReceiptProccessHandler(w http.ResponseWriter, r *http.Request) {
	// Validate for post request
	if r.Method != http.MethodPost {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}
	// Parse the JSON data from the request body
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// Create UUID and assign it to a ReceiptResponse Struct
	id := uuid.New().String()
	var receiptResponse models.ReceiptResponse
	receiptResponse.ID = id
	// Initialize Id and zero points to the points record system and do calculations
	points.PointsRecord[id] = 0
	points.CalculatePoints(receipt, id)
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	// Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(receiptResponse)

}

// Define a handler function for GET /receipts/{id}/points endpoint
func getReceiptProccesHandler(w http.ResponseWriter, r *http.Request) {
	// Validate for post request
	if r.Method != http.MethodGet {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}
	// Get id from url
	vars := mux.Vars(r)
	id := vars["id"]
	//Check for ID in the points record
	value, exists := points.PointsRecord[id]
	if !exists {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}
	// Set value to pointsResponse struct for response
	var pointsResponse models.PointsResponse
	pointsResponse.Points = float64(value)
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	// Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pointsResponse)
}

func main() {
	// Create new router
	router := mux.NewRouter()

	// Register handler functions for routes
	router.HandleFunc("/receipts/process", postReceiptProccessHandler)
	router.HandleFunc("/receipts/{id}/points", getReceiptProccesHandler).Methods("GET")

	// Set the address and port to listen on
	port := ":8080"
	fmt.Println("Starting server on port:", port)

	// Start the server and log errors if any
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
