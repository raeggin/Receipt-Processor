package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
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
func awardPointsForAlphanumeric(retailer string, id string) {
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

// Adds 50 points if the total is a round dollar amount with no cents
func awardPointsForRoundTotal(total string, id string) {
	var cents = (strings.Split(total, "."))[1]
	if cents == "00" {
		pointsRecord[id] = pointsRecord[id] + 50
	}
	fmt.Println("RULE 2 ROUND TOTAL: ", pointsRecord[id])
}

// Adds 25 points if the total is a multiple of 0.25
func awardPointsForMultipleOf(total string, id string) {
	floatValue, err := strconv.ParseFloat(total, 64)
	if err != nil {
		fmt.Println("Total can not be converted to float64 :", err)
		return
	}
	if int(floatValue*4)%4 == 0 {
		pointsRecord[id] = pointsRecord[id] + 25
	}
	fmt.Println("RULE 3 ROUND TOTAL: ", pointsRecord[id])
}

// Adds 5 points for every two items on the receipt
func awardPointsForEveryPair(items []Items, id string) {
	var pairs = len(items) / 2
	pointsRecord[id] = pointsRecord[id] + (pairs * 5)
	fmt.Println("RULE 4 ROUND TOTAL: ", pointsRecord[id])
}

func awardPointsForDescriptions(items []Items, id string) {
	for _, item := range items {
		if (len(strings.TrimSpace(item.ShortDescription))*3)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points := math.Ceil(price * 0.2)
			pointsRecord[id] = pointsRecord[id] + int(points)
		}
	}
	fmt.Println("RULE 5 ROUND TOTAL: ", pointsRecord[id])
}

// Adds 6 points if the day in the purchase date is odd
func awardPointsForOddDay(date string, id string) {
	day, err := strconv.Atoi(((strings.Split(date, "-"))[2]))
	if err != nil {
		fmt.Println("Error getting day: ", err)
	}
	if day%2 != 0 {
		pointsRecord[id] = pointsRecord[id] + 6
	}
	fmt.Println("RULE 6 ROUND TOTAL: ", pointsRecord[id])
}

// Adds 10 points if the time of purchase is after 2:00pm and before 4:00pm
func awardPointsForTime(purchaseTime string, id string) {
	hour, err := strconv.Atoi(((strings.Split(purchaseTime, ":"))[0]))
	if err != nil {
		fmt.Println("Error getting day: ", err)
	}
	if hour >= 14 && hour <= 16 {
		pointsRecord[id] = pointsRecord[id] + 10
	}
	fmt.Println("RULE 7 ROUND TOTAL: ", pointsRecord[id])
}

// Runs all rules for awarded points
func calculatePoints(receipt Receipt, id string) {
	awardPointsForAlphanumeric(receipt.Retailer, id)
	awardPointsForRoundTotal(receipt.Total, id)
	awardPointsForMultipleOf(receipt.Total, id)
	awardPointsForEveryPair(receipt.Items, id)
	awardPointsForOddDay(receipt.PurchaseDate, id)
	awardPointsForTime(receipt.PurchaseTime, id)
	awardPointsForDescriptions(receipt.Items, id)
}

// Define a handler function for POST /receipt/process endpoint
func postReceiptProccessHandler(w http.ResponseWriter, r *http.Request) {
	// Validate for post request
	if r.Method != http.MethodPost {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
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
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
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
