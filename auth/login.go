package auth

import (
	"assessment-mnc/entities"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer entities.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, existingCustomer := range customers {
		if existingCustomer.Email == existingCustomer.Email && existingCustomer.Password == existingCustomer.Password {
			// Create token with customer ID
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id": customer.ID,
			})
			// Sign the token with the secret key
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return the token as response
			json.NewEncoder(w).Encode(map[string]string{
				"token": tokenString,
			})
			return
		}
	}
	http.Error(w, "Invalid Email or Password", http.StatusUnauthorized)
}
