package auth

import (
	"assessment-mnc/entities"
	"encoding/json"
	"net/http"
)

var customers []entities.Customer

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer entities.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	customers = append(customers, customer)
	json.NewEncoder(w).Encode(struct {
		Message string            `json:"message"`
		Data    entities.Customer `json:"data"`
	}{
		Message: "Customer registered successfully",
		Data:    customer,
	})
}
