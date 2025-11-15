package handlers

import (
	"blog-app/auth"
	"fmt"

	models "blog-app/model"
	"html/template"
	"net/http"

	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

// Show form to create a post
func (h *PostHandler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/create_post.html")
	tmpl.Execute(w, nil)
}

// Handle post creation (form submission)
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.GetSession(r)
	userID := session.Values["user_id"]

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/create-post", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	post := models.Post{
		Title:   title,
		Content: content,
		UserID:  userID.(uint64),
	}
	fmt.Println("Received post:", title, content, userID)
	result := h.DB.Create(&post)
	fmt.Println("Insert result:", result.Error)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Display all posts
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	h.DB.Preload("User").Find(&posts)
	fmt.Println("Found posts:", len(posts))
	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(w, posts)
}
