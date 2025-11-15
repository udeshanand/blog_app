package auth

import (
	models "blog-app/model"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

// UserRegistration handles /register
func (a *AuthHandler) UserRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)

	if input.Username == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var count int64
	a.DB.Model(&models.User{}).Where("email = ?", input.Email).Or("username = ?", input.Username).Count(&count)
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashed),
	}

	if err := a.DB.Create(&user).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// UserLogin handles /login
func (a *AuthHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := a.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	session, _ := GetSession(r)
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	SaveSession(w, r, session)

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// Logout handler
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := GetSession(r)
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // delete session
	SaveSession(w, r, session)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}
