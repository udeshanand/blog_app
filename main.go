package main

import (
	"blog-app/cache"
	"blog-app/config"
	"blog-app/handlers"
	models "blog-app/model"
	"blog-app/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := config.Connect()
	fmt.Println("Database connected successfully!")

	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	// 🔥 Create in-memory cache instance
	inMemoryCache := cache.New()

	// Pass cache to handlers
	postHandler := &handlers.PostHandler{
		DB:    db,
		Cache: inMemoryCache,
	}

	authHandler := &handlers.AuthHandler{
		DB: db,
	}

	// Pass handlers into routes
	mux := routes.SetupRoutes(authHandler, postHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
