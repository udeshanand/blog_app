package auth

import (
	"net/http"
)

// AuthMiddleware ensures that a user is logged in before accessing protected routes
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := GetSession(r)

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized. Please login first.", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}
