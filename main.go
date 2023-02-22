package main

import (
	"assessment-mnc/auth"
	"assessment-mnc/transaction"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/transaction", auth.AuthTokenMiddleware(http.HandlerFunc(transaction.TransactionHandler)))
	http.HandleFunc("/logout", auth.LogoutHandler)

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
