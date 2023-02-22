package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization token from the request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}
		tokenParts := strings.Split(tokenString, "Bearer ")
		if len(tokenParts) != 2 {
			http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
			return
		}
		tokenString = tokenParts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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

		customerIDFloat64, ok := claims["id"].(float64)
		if !ok {
			http.Error(w, "Invalid customer ID in authorization token", http.StatusUnauthorized)
			return
		}
		customerID := strconv.Itoa(int(customerIDFloat64))

		// Pass the customer ID to the next handler
		ctx := context.WithValue(r.Context(), "customerID", customerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
