package points

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"receipt-processor/models"
)

var PointsRecord = map[string]int{}

// Award points for alphanumeric characters in retailer name
func AwardPointsForAlphanumeric(retailer string, id string) {
	var totalPoints = len(retailer)
	for _, char := range retailer {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			totalPoints = totalPoints - 1
		}
	}
	PointsRecord[id] = PointsRecord[id] + totalPoints
}

// Award points for round total (no cents)
func AwardPointsForRoundTotal(total string, id string) {
	var cents = (strings.Split(total, "."))[1]
	if cents == "00" {
		PointsRecord[id] = PointsRecord[id] + 50
	}
}

// Award points if the total is a multiple of 0.25
func AwardPointsForMultipleOf(total string, id string) {
	floatValue, err := strconv.ParseFloat(total, 64)
	if err != nil {
		fmt.Println("Total cannot be converted to float64:", err)
		return
	}
	if int(floatValue*4)%4 == 0 {
		PointsRecord[id] = PointsRecord[id] + 25
	}
}

// Award points for every pair of items
func AwardPointsForEveryPair(items []models.Items, id string) {
	var pairs = len(items) / 2
	PointsRecord[id] = PointsRecord[id] + (pairs * 5)
}

// Award points based on item descriptions
func AwardPointsForDescriptions(items []models.Items, id string) {
	for _, item := range items {
		if (len(strings.TrimSpace(item.ShortDescription)))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points := math.Ceil(price * 0.2)
			PointsRecord[id] = PointsRecord[id] + int(points)
		}
	}
}

// Award points if receipt total is greater than 10
func AwardPointsForLLM(total string, id string) {
	receiptTotal, _ := strconv.ParseFloat(total, 64)
	if receiptTotal > 10.00 {
		PointsRecord[id] = PointsRecord[id] + 5
	}
}

// Award points if the purchase day is odd
func AwardPointsForOddDay(date string, id string) {
	day, err := strconv.Atoi(((strings.Split(date, "-"))[2]))
	if err != nil {
		fmt.Println("Error getting day:", err)
	}
	if day%2 != 0 {
		PointsRecord[id] = PointsRecord[id] + 6
	}
}

// Award points for time between 2:00pm and 4:00pm
func AwardPointsForTime(purchaseTime string, id string) {
	hour, err := strconv.Atoi(((strings.Split(purchaseTime, ":"))[0]))
	if err != nil {
		fmt.Println("Error getting hour:", err)
	}
	if hour >= 14 && hour <= 16 {
		PointsRecord[id] = PointsRecord[id] + 10
	}
}

// Calculate all points for a given receipt
func CalculatePoints(receipt models.Receipt, id string) {
	AwardPointsForAlphanumeric(receipt.Retailer, id)
	AwardPointsForRoundTotal(receipt.Total, id)
	AwardPointsForMultipleOf(receipt.Total, id)
	AwardPointsForEveryPair(receipt.Items, id)
	AwardPointsForDescriptions(receipt.Items, id)
	AwardPointsForLLM(receipt.Total, id)
	AwardPointsForOddDay(receipt.PurchaseDate, id)
	AwardPointsForTime(receipt.PurchaseTime, id)
}
