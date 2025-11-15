package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Global session store
var Store = sessions.NewCookieStore([]byte("super-secret-session-key"))

// initialize store options
func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// GetSession retrieves the current session
func GetSession(r *http.Request) (*sessions.Session, error) {
	return Store.Get(r, "blog_session")
}

// SetSessionValue sets a key-value pair in the session
func SetSessionValue(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, _ := GetSession(r)
	session.Values[key] = value
	return session.Save(r, w)
}

// GetSessionValue retrieves a value by key
func GetSessionValue(r *http.Request, key string) interface{} {
	session, _ := GetSession(r)
	return session.Values[key]
}

// ClearSession deletes all session data (used on logout)
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := GetSession(r)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
func SaveSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) error {
	return session.Save(r, w)
}
