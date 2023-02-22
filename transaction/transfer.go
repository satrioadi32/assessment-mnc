package transaction

import (
	"assessment-mnc/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

var customers []entities.Customer
var transactions []entities.Transaction

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get the authorization token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Extract customer ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}
	customerID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid customer ID in authorization token", http.StatusUnauthorized)
		return
	}
	// Check if customer exists
	var customer *entities.Customer
	for _, existingCustomer := range customers {
		if strconv.Itoa(existingCustomer.ID) == customerID {
			customer = &existingCustomer
			break
		}
	}
	if customer == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	// Check if the request is a transfer request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var transferRequest entities.TransferRequest
	err = json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the transfer amount is valid
	if transferRequest.Amount <= 0 {
		http.Error(w, "Invalid transfer amount", http.StatusBadRequest)
		return
	}
	// Check if the transfer is within the customer's balance
	if customer.Balance < transferRequest.Amount {
		http.Error(w, "Insufficient balance", http.StatusForbidden)
		return
	}
	// Perform the transfer
	customer.Balance -= transferRequest.Amount
	for _, existingCustomer := range customers {
		if existingCustomer.ID == transferRequest.RecipientID {
			existingCustomer.Balance += transferRequest.Amount
			break
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transfer successful",
	})

	transaction := map[string]interface{}{
		"customerID": customerID,
		"action":     "transfer",
	}
	AddToHistory(transaction)

	// Return Response
	response := map[string]interface{}{
		"message": "Transfer successful",
	}
	json.NewEncoder(w).Encode(response)
}
