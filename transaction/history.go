package transaction

import (
	"encoding/json"
	"net/http"
)

var history []string

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from context
	customerID, ok := r.Context().Value("customerID").(string)
	if !ok {
		http.Error(w, "Invalid customer ID in request context", http.StatusInternalServerError)
		return
	}

	// Filter history by customer ID
	var filteredHistory []string
	for _, item := range history {
		var transaction map[string]interface{}
		err := json.Unmarshal([]byte(item), &transaction)
		if err == nil && transaction["customerID"] == customerID {
			filteredHistory = append(filteredHistory, item)
		}
	}

	// Return filtered history
	json.NewEncoder(w).Encode(filteredHistory)
}

func AddToHistory(transaction map[string]interface{}) {
	transactionJSON, _ := json.Marshal(transaction)
	history = append(history, string(transactionJSON))
}
