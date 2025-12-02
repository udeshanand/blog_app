package handlers

import (
	"blog-app/auth"
	"blog-app/cache"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	models "blog-app/model"

	"gorm.io/gorm"
)

type PostHandler struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

// Show form to create a post
func (h *PostHandler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create_post.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}
	tmpl.Execute(w, nil)
}

// Handle post creation (POST)
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

	// Insert into DB
	if err := h.DB.Create(&post).Error; err != nil {
		http.Error(w, "Failed to create post", 500)
		return
	}

	// Clear cache when new post created
	h.Cache.Flush()

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// LIST POSTS (WITH CACHE + PAGINATION)
// new
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("posts_page_%d_limit_%d", page, limit)

	// --- TOTAL COUNT (ALWAYS REQUIRED) ---
	var total int64
	h.DB.Model(&models.Post{}).Count(&total)

	lastPage := int((total + int64(limit) - 1) / int64(limit))

	hasPrev := page > 1
	hasNext := page < lastPage

	// --- CACHE CHECK ---
	if cached, found := h.Cache.Get(cacheKey); found {
		fmt.Println("CACHE HIT")

		if posts, ok := cached.([]models.Post); ok {
			tmpl, _ := template.ParseFiles("templates/dashboard.html")
			tmpl.Execute(w, map[string]interface{}{
				"Posts":    posts,
				"Page":     page,
				"Limit":    limit,
				"LastPage": lastPage,
				"HasPrev":  hasPrev,
				"HasNext":  hasNext,
				"PrevPage": page - 1,
				"NextPage": page + 1,
			})
			return
		}

		h.Cache.Delete(cacheKey) // wrong type → remove
	}

	fmt.Println("CACHE MISS")

	var posts []models.Post
	err := h.DB.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	h.Cache.Set(cacheKey, posts, cache.DefaultExpiration)

	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(w, map[string]interface{}{
		"Posts":    posts,
		"Page":     page,
		"Limit":    limit,
		"LastPage": lastPage,
		"HasPrev":  hasPrev,
		"HasNext":  hasNext,
		"PrevPage": page - 1,
		"NextPage": page + 1,
	})
}

// new
func (h *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	// Get post ID from query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Fetch post with user
	var post models.Post
	if err := h.DB.Preload("User").First(&post, id).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Debug
	fmt.Printf("POST DEBUG: %+v\n", post)

	// Parse template with safeHTML function
	tmpl := template.Must(template.New("post.html").
		Funcs(template.FuncMap{
			"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		}).
		ParseFiles("templates/post.html"))

	// Execute template
	if err := tmpl.Execute(w, post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
