package routes

import (
	"net/http"

	"blog-app/auth"

	"blog-app/handlers"

	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := &handlers.AuthHandler{DB: db}

	postHandler := &handlers.PostHandler{DB: db}

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.RegisterPage(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.RegisterUser(w, r)
		}
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.LoginPage(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.LoginUser(w, r)
		}
	})

	mux.HandleFunc("/", authHandler.HomePage)

	mux.HandleFunc("/logout", authHandler.LogoutUser)

	mux.HandleFunc("/dashboard", auth.AuthMiddleware(postHandler.ListPosts))

	mux.HandleFunc("/create-post", auth.AuthMiddleware(postHandler.CreatePostPage))
	mux.HandleFunc("/submit-post", auth.AuthMiddleware(postHandler.CreatePost))
	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
