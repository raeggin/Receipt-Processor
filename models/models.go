package models

// Item represents an individual item in a receipt
type Items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt represents the full receipt data
type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        string  `json:"total"`
	Items        []Items `json:"items"`
}

// ReceiptResponse is the response after posting a receipt
type ReceiptResponse struct {
	ID string `json:"id"`
}

// PointsResponse is the response for the number of points earned
type PointsResponse struct {
	Points float64 `json:"points"`
}
