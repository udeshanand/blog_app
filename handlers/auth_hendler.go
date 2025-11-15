package handlers

import (
	"blog-app/config"
	"html/template"
	"log"
	"net/http"

	"blog-app/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	models "blog-app/model"
)

type AuthHandler struct {
	DB *gorm.DB
}

// Render HTML template
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}
	t.Execute(w, data)
}

// Register Page (GET)
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

// Register User (POST)
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{Username: username, Email: email, PasswordHash: string(hash)}
	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Page (GET)
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

// Login User (POST)
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user models.User
	result := h.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare password with hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	session, _ := auth.GetSession(r)
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username

	if err := auth.SaveSession(w, r, session); err != nil {
		log.Println("Session save error:", err)
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout (GET)
func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.GetSession(r)
	delete(session.Values, "user_id")
	session.Options.MaxAge = -1 // clear cookie
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Home Page
func (h *AuthHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.GetSession(r)
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user models.User
	h.DB.First(&user, userID)
	renderTemplate(w, "home.html", user)
}

func (h *AuthHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.GetSession(r)

	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading dashboard template", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	tmpl.Execute(w, data)
}

// Initialize database for this handler
func InitDB() *gorm.DB {
	db := config.Connect()
	db.AutoMigrate(&models.User{})
	return db
}
