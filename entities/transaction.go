package entities

import "time"

type Transaction struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Merchant   string    `json:"merchant"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
