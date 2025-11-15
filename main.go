package main

import (
	"blog-app/config"
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

	mux := routes.SetupRoutes(db)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
