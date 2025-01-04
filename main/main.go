package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// Keep record of points
var pointsRecord = map[string]int{}

// Create ReceiptResponse struct; this is the response structure for when you create a new receipt
type ReceiptResponse struct {
	Id string `json:"id"`
}

// Create Points Struct; this is the response for when you retrieve the points of a receipt with a given Id
type PointsResponse struct {
	Points float64 `json:"points"`
}

// Define a handler function for root endpoint
func okHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	// Send an "OK" message
	fmt.Fprint(w, "OK")
}

// Calculates one point for every alphanumeric character in the retailer name
func alphanumericCalculator(retailer string, id string) {
	var totalPoints = len(retailer)
	fmt.Println("Points: ", totalPoints)
	for _, char := range retailer {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			totalPoints = totalPoints - 1
		}
	}
	pointsRecord[id] = pointsRecord[id] + totalPoints
	fmt.Println("RULE 1 TOTAL POINTS: ", pointsRecord[id])
}

func calculatePoints(receipt Receipt, id string) {
	alphanumericCalculator(receipt.Retailer, id)
}

// Define a handler function for POST /receipt/process endpoint
func postReceiptProccessHandler(w http.ResponseWriter, r *http.Request) {
	// Validate for post request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
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
	// Initialize Id and zero points to the points record system and do calculations
	pointsRecord[id] = 0
	calculatePoints(receipt, id)
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
		http.Error(w, "Invalid request method:", http.StatusBadRequest)
		return
	}
	// Get id from url
	vars := mux.Vars(r)
	id := vars["id"]
	//Check for ID in the points record
	value, exists := pointsRecord[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	// Set value to pointsResponse struct for response
	var pointsResponse PointsResponse
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
	router.HandleFunc("/", okHandler)
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
