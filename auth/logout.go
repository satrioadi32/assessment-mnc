package auth

import (
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Delete token from cookie or Authorization header
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	})
	w.Header().Set("Authorization", "")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
