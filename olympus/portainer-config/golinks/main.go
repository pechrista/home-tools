package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Link struct {
	Slug      string    `json:"slug"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type AddLinkRequest struct {
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

type RemoveLinkRequest struct {
	Slug string `json:"slug"`
}

var (
	db        *sql.DB
	adminUser string
	adminPass string
)

func main() {
	// Get configuration from environment
	dbPath := getEnv("DB_PATH", "./data/links.db")
	listenAddr := getEnv("LISTEN_ADDR", "0.0.0.0:8080")
	adminUser = os.Getenv("ADMIN_USER")
	adminPass = os.Getenv("ADMIN_PASS")

	// Initialize database
	if err := initDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup routes
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/admin/add", basicAuth(handleAdminAdd))
	http.HandleFunc("/admin/remove", basicAuth(handleAdminRemove))

	// Start server
	log.Printf("Starting golinks server on %s", listenAddr)
	log.Printf("Database: %s", dbPath)
	if adminUser != "" {
		log.Printf("Admin authentication enabled")
	}

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initDB(dbPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		slug TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Root path - list all links
	if path == "" {
		handleListLinks(w, r)
		return
	}

	// Slug lookup
	slug := path
	link, err := getLink(slug)
	if err != nil {
		log.Printf("404 - Slug not found: %s (from %s)", slug, r.RemoteAddr)
		http.NotFound(w, r)
		return
	}

	log.Printf("302 - Redirecting %s -> %s (from %s)", slug, link.URL, r.RemoteAddr)
	http.Redirect(w, r, link.URL, http.StatusFound)
}

func handleListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := getAllLinks()
	if err != nil {
		log.Printf("Error fetching links: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Go Links</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			padding: 2rem;
		}
		.container {
			max-width: 900px;
			margin: 0 auto;
			background: white;
			border-radius: 12px;
			box-shadow: 0 20px 60px rgba(0,0,0,0.3);
			padding: 2rem;
		}
		h1 {
			color: #333;
			margin-bottom: 0.5rem;
			font-size: 2rem;
		}
		.subtitle {
			color: #666;
			margin-bottom: 2rem;
			font-size: 0.95rem;
		}
		.empty {
			text-align: center;
			padding: 3rem;
			color: #999;
		}
		.link-list {
			list-style: none;
		}
		.link-item {
			border-bottom: 1px solid #eee;
			padding: 1rem 0;
			transition: background 0.2s;
		}
		.link-item:last-child {
			border-bottom: none;
		}
		.link-item:hover {
			background: #f8f9fa;
			margin: 0 -1rem;
			padding: 1rem;
			border-radius: 6px;
		}
		.link-slug {
			font-weight: 600;
			color: #667eea;
			text-decoration: none;
			font-size: 1.1rem;
			display: inline-block;
			margin-bottom: 0.25rem;
		}
		.link-slug:hover {
			color: #764ba2;
			text-decoration: underline;
		}
		.link-url {
			color: #666;
			font-size: 0.9rem;
			word-break: break-all;
			display: block;
		}
		.link-date {
			color: #999;
			font-size: 0.85rem;
			margin-top: 0.25rem;
		}
		.count {
			background: #667eea;
			color: white;
			padding: 0.25rem 0.75rem;
			border-radius: 20px;
			font-size: 0.85rem;
			display: inline-block;
			margin-left: 0.5rem;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>🔗 Go Links <span class="count">{{.Count}}</span></h1>
		<p class="subtitle">Internal URL Shortener</p>
		{{if .Links}}
			<ul class="link-list">
			{{range .Links}}
				<li class="link-item">
					<a href="/{{.Slug}}" class="link-slug">go/{{.Slug}}</a>
					<span class="link-url">→ {{.URL}}</span>
					<div class="link-date">Created {{.CreatedAt.Format "Jan 02, 2006 15:04"}}</div>
				</li>
			{{end}}
			</ul>
		{{else}}
			<div class="empty">
				<p>No links yet. Add one via POST /admin/add</p>
			</div>
		{{end}}
	</div>
</body>
</html>`

	t, err := template.New("links").Parse(tmpl)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Links []Link
		Count int
	}{
		Links: links,
		Count: len(links),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleAdminAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate slug
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" || req.Slug == "admin" {
		http.Error(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	// Validate URL
	req.URL = strings.TrimSpace(req.URL)
	if !isValidURL(req.URL) {
		http.Error(w, "Invalid URL - must start with http:// or https://", http.StatusBadRequest)
		return
	}

	// Insert link
	if err := addLink(req.Slug, req.URL); err != nil {
		log.Printf("Error adding link: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			http.Error(w, "Slug already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Link added: %s -> %s (by %s)", req.Slug, req.URL, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"slug":   req.Slug,
		"url":    req.URL,
	})
}

func handleAdminRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.Slug = strings.TrimSpace(req.Slug)
	if req.Slug == "" || req.Slug == "admin" {
		http.Error(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	if err := removeLink(req.Slug); err != nil {
		log.Printf("Error removing link: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Slug not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Link removed: %s (by %s)", req.Slug, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "removed",
		"slug":   req.Slug,
	})
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If admin credentials not set, allow access
		if adminUser == "" || adminPass == "" {
			log.Printf("Warning: Admin endpoint accessed without authentication configured")
			next(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != adminUser || pass != adminPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Unauthorized admin access attempt from %s", r.RemoteAddr)
			return
		}

		next(w, r)
	}
}

func getLink(slug string) (*Link, error) {
	var link Link
	err := db.QueryRow("SELECT slug, url, created_at FROM links WHERE slug = ?", slug).
		Scan(&link.Slug, &link.URL, &link.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("link not found")
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func getAllLinks() ([]Link, error) {
	rows, err := db.Query("SELECT slug, url, created_at FROM links ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.Slug, &link.URL, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

func addLink(slug, url string) error {
	_, err := db.Exec("INSERT INTO links (slug, url) VALUES (?, ?)", slug, url)
	return err
}

func removeLink(slug string) error {
	res, err := db.Exec("DELETE FROM links WHERE slug = ?", slug)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func isValidURL(urlStr string) bool {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
