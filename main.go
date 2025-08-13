package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/skip2/go-qrcode"
)

// URLStore represents our in-memory storage
type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

// ShortenRequest represents the JSON request for shortening a URL
type ShortenRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
}

// ShortenResponse represents the JSON response for shortening a URL
type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var store = &URLStore{
	urls: make(map[string]string),
}

const (
	baseURL    = "http://localhost:8080"
	codeLength = 6
)

func main() {
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/qr/", handleQRCode)
	http.HandleFunc("/favicon.ico", handleFavicon)
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("URL Shortener running on %s\n", baseURL)
	fmt.Println("Usage:")
	fmt.Println("  POST /shorten - Shorten a URL")
	fmt.Println("  GET /{code}   - Redirect to original URL")
	fmt.Println("  GET /qr/{code} - Get QR code for short URL")
	fmt.Println("  GET /favicon.ico - Favicon")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleShorten handles POST requests to shorten URLs
func handleShorten(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST requests are supported")
		return
	}

	// Validate content type
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid content type", "Content-Type must be application/json")
		return
	}

	// Limit request body size (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req ShortenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", "Request body must be valid JSON with 'url' field")
		return
	}

	// Validate URL
	if !isValidURL(req.URL) {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid URL", "URL must be a valid HTTP or HTTPS URL")
		return
	}

	// Sanitize URL
	sanitizedURL := sanitizeURL(req.URL)

	var shortCode string
	var err error

	// Use custom code if provided, otherwise generate random code
	if req.CustomCode != "" {
		// Validate custom code
		if !isValidCustomCode(req.CustomCode) {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid custom code", "Custom code must be 3-20 characters, alphanumeric and hyphens only")
			return
		}

		// Check if custom code is reserved
		if isReservedCode(req.CustomCode) {
			sendErrorResponse(w, http.StatusBadRequest, "Reserved code", "This custom code is reserved and cannot be used")
			return
		}

		// Check if custom code already exists
		store.mu.RLock()
		_, exists := store.urls[req.CustomCode]
		store.mu.RUnlock()

		if exists {
			sendErrorResponse(w, http.StatusConflict, "Code already exists", "This custom code is already in use")
			return
		}

		shortCode = req.CustomCode
	} else {
		// Generate random short code
		shortCode, err = generateShortCode()
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Generation failed", "Failed to generate short code")
			return
		}
	}

	// Store the mapping
	store.mu.Lock()
	store.urls[shortCode] = sanitizedURL
	store.mu.Unlock()

	// Create response
	response := ShortenResponse{
		ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortCode),
		OriginalURL: sanitizedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Printf("Shortened URL: %s -> %s", sanitizedURL, response.ShortURL)
}

// handleQRCode handles GET requests to generate QR codes for short URLs
func handleQRCode(w http.ResponseWriter, r *http.Request) {
	// Extract short code from path
	shortCode := strings.TrimPrefix(r.URL.Path, "/qr/")

	// Validate short code format
	if !isValidShortCode(shortCode) {
		http.NotFound(w, r)
		return
	}

	// Look up original URL
	store.mu.RLock()
	_, exists := store.urls[shortCode]
	store.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		log.Printf("Short code not found for QR: %s", shortCode)
		return
	}

	// Generate QR code for the short URL
	shortURL := fmt.Sprintf("%s/%s", baseURL, shortCode)
	qrCode, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("QR code generation failed: %v", err)
		return
	}

	// Set headers for PNG image
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(qrCode)))
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Write QR code image
	w.Write(qrCode)
	log.Printf("QR code generated for: %s", shortCode)
}

// handleFavicon handles favicon requests to prevent 404 errors
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	// Return a simple 1x1 transparent PNG
	favicon := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	w.WriteHeader(http.StatusOK)
	w.Write(favicon)
}

// handleRedirect handles GET requests to redirect short URLs
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Extract short code from path
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	// Handle root path
	if shortCode == "" {
		sendHTMLResponse(w, http.StatusOK, getHomePage())
		return
	}

	// Validate short code format
	if !isValidShortCode(shortCode) {
		http.NotFound(w, r)
		return
	}

	// Look up original URL
	store.mu.RLock()
	originalURL, exists := store.urls[shortCode]
	store.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		log.Printf("Short code not found: %s", shortCode)
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	log.Printf("Redirected: %s -> %s", shortCode, originalURL)
}

// isValidURL validates if a string is a valid HTTP/HTTPS URL
func isValidURL(str string) bool {
	if str == "" {
		return false
	}

	// Parse URL
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// Check scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Check host
	if u.Host == "" {
		return false
	}

	return true
}

// sanitizeURL cleans and normalizes a URL
func sanitizeURL(rawURL string) string {
	// Parse and reconstruct URL to normalize it
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Ensure scheme is lowercase
	u.Scheme = strings.ToLower(u.Scheme)
	
	// Ensure host is lowercase
	u.Host = strings.ToLower(u.Host)

	return u.String()
}

// isValidShortCode validates the format of a short code
func isValidShortCode(code string) bool {
	// For custom codes, allow variable length
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	// Check if code contains only alphanumeric characters and hyphens
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') ||
			 char == '-') {
			return false
		}
	}

	return true
}

// isValidCustomCode validates the format of a custom short code
func isValidCustomCode(code string) bool {
	// Length check
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	// Cannot start or end with hyphen
	if strings.HasPrefix(code, "-") || strings.HasSuffix(code, "-") {
		return false
	}

	// Cannot have consecutive hyphens
	if strings.Contains(code, "--") {
		return false
	}

	// Check if code contains only alphanumeric characters and hyphens
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') ||
			 char == '-') {
			return false
		}
	}

	return true
}

// isReservedCode checks if a code is reserved and cannot be used
func isReservedCode(code string) bool {
	reserved := []string{
		"admin", "api", "www", "ftp", "mail", "email", "support", "help",
		"about", "contact", "terms", "privacy", "legal", "blog", "news",
		"docs", "documentation", "download", "downloads", "upload", "uploads",
		"static", "assets", "css", "js", "img", "images", "favicon",
		"robots", "sitemap", "feed", "rss", "atom", "xml", "json",
		"login", "logout", "signin", "signup", "register", "auth",
		"dashboard", "profile", "account", "settings", "config",
		"test", "testing", "dev", "development", "staging", "prod", "production",
		"qr", "shorten", "short", "url", "link", "redirect",
	}

	lowerCode := strings.ToLower(code)
	for _, reservedWord := range reserved {
		if lowerCode == reservedWord {
			return true
		}
	}

	return false
}

// generateShortCode generates a random short code
func generateShortCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	for attempts := 0; attempts < 10; attempts++ {
		code := make([]byte, codeLength)
		
		// Generate random bytes
		if _, err := rand.Read(code); err != nil {
			return "", err
		}

		// Convert to charset
		for i := range code {
			code[i] = charset[int(code[i])%len(charset)]
		}

		shortCode := string(code)

		// Check for collision
		store.mu.RLock()
		_, exists := store.urls[shortCode]
		store.mu.RUnlock()

		if !exists {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after 10 attempts")
}

// sendErrorResponse sends a JSON error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
	log.Printf("Error response: %d - %s: %s", statusCode, error, message)
}

// sendHTMLResponse sends an HTML response
func sendHTMLResponse(w http.ResponseWriter, statusCode int, html string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(html))
}

// getHomePage returns a modern HTML page for testing
func getHomePage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>QuickLink - URL Shortener</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            margin: 0;
            padding: 0;
        }
        
        .page-wrapper {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        
        .navbar {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            padding: 15px 0;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .nav-content {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        
        .nav-logo {
            display: flex;
            align-items: center;
            gap: 10px;
            color: white;
            text-decoration: none;
            font-weight: 600;
            font-size: 1.2rem;
        }
        
        .nav-logo i {
            font-size: 1.5rem;
        }
        
        .nav-links {
            display: flex;
            gap: 20px;
            align-items: center;
        }
        
        .nav-link {
            color: rgba(255, 255, 255, 0.8);
            text-decoration: none;
            font-size: 0.9rem;
            transition: color 0.3s ease;
        }
        
        .nav-link:hover {
            color: white;
        }
        
        .main-content {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 40px 20px;
        }
        
        .container {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 50px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 100%;
            position: relative;
        }
        
        .footer-section {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            padding: 30px 0;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            margin-top: auto;
        }
        
        .footer-content {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
            text-align: center;
        }
        
        .footer-links {
            display: flex;
            justify-content: center;
            gap: 30px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .footer-link {
            color: rgba(255, 255, 255, 0.8);
            text-decoration: none;
            font-size: 0.9rem;
            transition: color 0.3s ease;
        }
        
        .footer-link:hover {
            color: white;
        }
        
        .footer-info {
            color: rgba(255, 255, 255, 0.6);
            font-size: 0.8rem;
        }
        
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .logo {
            display: inline-flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 15px;
        }
        
        .logo i {
            font-size: 2.5rem;
            background: linear-gradient(135deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        h1 {
            font-size: 2.5rem;
            font-weight: 700;
            background: linear-gradient(135deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin: 0;
        }
        
        .subtitle {
            color: #6b7280;
            font-size: 1.1rem;
            font-weight: 400;
            margin-bottom: 10px;
        }
        
        .description {
            color: #9ca3af;
            font-size: 0.95rem;
            line-height: 1.5;
            margin-bottom: 20px;
        }
        
        .form-container {
            margin-bottom: 30px;
        }
        
        .input-group {
            position: relative;
            margin-bottom: 20px;
        }
        
        .input-icon {
            position: absolute;
            left: 15px;
            top: 50%;
            transform: translateY(-50%);
            color: #9ca3af;
            font-size: 1.1rem;
        }
        
        input[type="url"] {
            width: 100%;
            padding: 15px 15px 15px 45px;
            border: 2px solid #e5e7eb;
            border-radius: 12px;
            font-size: 1rem;
            transition: all 0.3s ease;
            background: #f9fafb;
        }
        
        input[type="url"]:focus,
        input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
            background: white;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        input[type="text"] {
            width: 100%;
            padding: 15px 15px 15px 45px;
            border: 2px solid #e5e7eb;
            border-radius: 12px;
            font-size: 1rem;
            transition: all 0.3s ease;
            background: #f9fafb;
        }
        
        .input-hint {
            display: block;
            color: #9ca3af;
            font-size: 0.8rem;
            margin-top: 5px;
            margin-left: 45px;
        }
        
        .shorten-btn {
            width: 100%;
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 15px;
            border: none;
            border-radius: 12px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
        }
        
        .shorten-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.3);
        }
        
        .shorten-btn:active {
            transform: translateY(0);
        }
        
        .result {
            margin-top: 20px;
            padding: 20px;
            background: linear-gradient(135deg, #d1fae5, #a7f3d0);
            border-radius: 12px;
            border-left: 4px solid #10b981;
            animation: slideIn 0.3s ease;
        }
        
        .result strong {
            color: #065f46;
            display: block;
            margin-bottom: 8px;
        }
        
        .result a {
            color: #059669;
            text-decoration: none;
            font-weight: 500;
            word-break: break-all;
        }
        
        .result a:hover {
            text-decoration: underline;
        }
        
        .error {
            margin-top: 20px;
            padding: 20px;
            background: linear-gradient(135deg, #fee2e2, #fecaca);
            border-radius: 12px;
            border-left: 4px solid #ef4444;
            animation: slideIn 0.3s ease;
        }
        
        .error strong {
            color: #991b1b;
        }
        
        .footer {
            margin-top: 30px;
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid #e5e7eb;
        }
        
        .github-link {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            color: #6b7280;
            text-decoration: none;
            font-size: 0.9rem;
            transition: color 0.3s ease;
            margin-bottom: 10px;
        }
        
        .github-link:hover {
            color: #374151;
        }
        
        .github-link i {
            font-size: 1.1rem;
        }
        
        .trademark {
            color: #9ca3af;
            font-size: 0.8rem;
            margin-top: 10px;
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        .loading {
            opacity: 0.7;
            pointer-events: none;
        }
        
        .loading .shorten-btn {
            background: #9ca3af;
        }
        
        .url-container {
            display: flex;
            align-items: center;
            gap: 10px;
            margin: 10px 0;
            padding: 12px;
            background: rgba(255, 255, 255, 0.8);
            border-radius: 8px;
            border: 1px solid #e5e7eb;
        }
        
        .short-url {
            flex: 1;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.9rem;
            padding: 8px 12px;
            background: #f8fafc;
            border-radius: 6px;
            border: 1px solid #e2e8f0;
        }
        
        .copy-btn {
            background: #667eea;
            color: white;
            border: none;
            padding: 8px 12px;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s ease;
            font-size: 0.9rem;
        }
        
        .copy-btn:hover {
            background: #5a67d8;
            transform: translateY(-1px);
        }
        
        .qr-section {
            margin-top: 20px;
            padding: 20px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 12px;
            border: 1px solid #e5e7eb;
        }
        
        .qr-header {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 15px;
            color: #374151;
            font-weight: 500;
        }
        
        .qr-header i {
            color: #667eea;
        }
        
        .qr-container {
            text-align: center;
        }
        
        .qr-code {
            max-width: 160px;
            width: 100%;
            height: auto;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            animation: fadeIn 0.5s ease;
            padding: 8px;
            background: white;
        }
        
        .qr-actions {
            margin-top: 12px;
        }
        
        .qr-download-btn {
            background: #10b981;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s ease;
            font-size: 0.85rem;
            font-weight: 500;
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
        
        .qr-download-btn:hover {
            background: #059669;
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
        }
        
        .result-container {
            animation: slideInUp 0.4s ease-out;
        }
        
        @keyframes slideInUp {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: scale(0.95);
            }
            to {
                opacity: 1;
                transform: scale(1);
            }
        }
        
        .copy-success {
            background: #10b981 !important;
            transform: scale(0.95);
        }
        
        .copy-success i {
            animation: checkPulse 0.3s ease;
        }
        
        @keyframes checkPulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.2); }
            100% { transform: scale(1); }
        }
        
        @media (max-width: 480px) {
            .container {
                padding: 30px 20px;
                margin: 10px;
            }
            
            h1 {
                font-size: 2rem;
            }
            
            .subtitle {
                font-size: 1rem;
            }
            
            .url-container {
                flex-direction: column;
                gap: 8px;
            }
            
            .short-url {
                width: 100%;
                text-align: center;
            }
            
            .qr-code {
                max-width: 150px;
            }
        }
    </style>
</head>
<body>
    <div class="page-wrapper">
        <!-- Navigation Bar -->
        <nav class="navbar">
            <div class="nav-content">
                <a href="/" class="nav-logo">
                    <i class="fas fa-link"></i>
                    QuickLink
                </a>
                <div class="nav-links">
                    <a href="#features" class="nav-link">Features</a>
                    <a href="#api" class="nav-link">API</a>
                    <a href="https://github.com/Neorex80/Quick-Link" target="_blank" class="nav-link">
                        <i class="fab fa-github"></i> GitHub
                    </a>
                </div>
            </div>
        </nav>

        <!-- Main Content -->
        <main class="main-content">
            <div class="container">
                <div class="header">
                    <div class="logo">
                        <i class="fas fa-link"></i>
                        <h1>QuickLink</h1>
                    </div>
                    <p class="subtitle">Transform long URLs into short, shareable links</p>
                    <p class="description">
                        A fast, secure, and reliable URL shortener built with Go. 
                        Perfect for social media, emails, and anywhere you need clean, compact links.
                    </p>
                </div>
                
                <div class="form-container">
                    <form id="shortenForm">
                        <div class="input-group">
                            <i class="fas fa-globe input-icon"></i>
                            <input type="url" id="urlInput" placeholder="https://example.com/your/very/long/url/here" required>
                        </div>
                        <div class="input-group">
                            <i class="fas fa-edit input-icon"></i>
                            <input type="text" id="customCodeInput" placeholder="my-custom-code (optional)" maxlength="20">
                            <small class="input-hint">3-20 characters, letters, numbers, and hyphens only</small>
                        </div>
                        <button type="submit" class="shorten-btn">
                            <i class="fas fa-magic"></i>
                            Shorten URL
                        </button>
                    </form>
                </div>
                
                <div id="result"></div>
                
                <div class="footer">
                    <a href="https://github.com/Neorex80" target="_blank" class="github-link">
                        <i class="fab fa-github"></i>
                        Created by Neorex80
                    </a>
                    <div class="trademark">
                        © 2025 QuickLink. Built with ❤️ using Go
                    </div>
                </div>
            </div>
        </main>

        <!-- Footer -->
        <footer class="footer-section">
            <div class="footer-content">
                <div class="footer-links">
                    <a href="#about" class="footer-link">About</a>
                    <a href="#privacy" class="footer-link">Privacy</a>
                    <a href="#terms" class="footer-link">Terms</a>
                    <a href="#contact" class="footer-link">Contact</a>
                    <a href="https://github.com/Neorex80/Quick-Link" target="_blank" class="footer-link">GitHub</a>
                </div>
                <div class="footer-info">
                    © 2025 QuickLink. Built with ❤️ using Go. Fast, secure, and reliable URL shortening service.
                </div>
            </div>
        </footer>
    </div>

    <script>
        document.getElementById('shortenForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const url = document.getElementById('urlInput').value;
            const customCode = document.getElementById('customCodeInput').value.trim();
            const resultDiv = document.getElementById('result');
            const container = document.querySelector('.container');
            const submitBtn = document.querySelector('.shorten-btn');
            
            // Add loading state
            container.classList.add('loading');
            submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Shortening...';
            
            // Prepare request body
            const requestBody = { url: url };
            if (customCode) {
                requestBody.custom_code = customCode;
            }
            
            try {
                const response = await fetch('/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(requestBody)
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    // Extract short code from URL
                    const shortCode = data.short_url.split('/').pop();
                    const qrUrl = '/qr/' + shortCode;
                    
                    resultDiv.innerHTML = 
                        '<div class="result-container">' +
                            '<div class="result">' +
                                '<strong><i class="fas fa-check-circle"></i> Success! Your shortened URL:</strong>' +
                                '<div class="url-container">' +
                                    '<a href="' + data.short_url + '" target="_blank" class="short-url">' + data.short_url + '</a>' +
                                    '<button class="copy-btn" onclick="copyToClipboard(\'' + data.short_url + '\')">' +
                                        '<i class="fas fa-copy"></i>' +
                                    '</button>' +
                                '</div>' +
                                '<div class="qr-section">' +
                                    '<div class="qr-header">' +
                                        '<i class="fas fa-qrcode"></i>' +
                                        '<span>QR Code for easy sharing</span>' +
                                    '</div>' +
                                    '<div class="qr-container">' +
                                        '<img src="' + qrUrl + '" alt="QR Code" class="qr-code" loading="lazy">' +
                                        '<div class="qr-actions">' +
                                            '<button class="qr-download-btn" onclick="downloadQR(\'' + qrUrl + '\', \'' + shortCode + '\')">' +
                                                '<i class="fas fa-download"></i> Download' +
                                            '</button>' +
                                        '</div>' +
                                    '</div>' +
                                '</div>' +
                            '</div>' +
                        '</div>';
                } else {
                    resultDiv.innerHTML = '<div class="error"><strong><i class="fas fa-exclamation-triangle"></i> Error:</strong> ' + data.message + '</div>';
                }
            } catch (error) {
                resultDiv.innerHTML = '<div class="error"><strong><i class="fas fa-exclamation-triangle"></i> Error:</strong> Failed to shorten URL. Please try again.</div>';
            } finally {
                // Remove loading state
                container.classList.remove('loading');
                submitBtn.innerHTML = '<i class="fas fa-magic"></i> Shorten URL';
            }
        });
        
        // Copy to clipboard function
        function copyToClipboard(text) {
            if (navigator.clipboard && navigator.clipboard.writeText) {
                navigator.clipboard.writeText(text).then(function() {
                    showCopySuccess();
                }).catch(function() {
                    fallbackCopyToClipboard(text);
                });
            } else {
                fallbackCopyToClipboard(text);
            }
        }
        
        // Fallback copy function for older browsers
        function fallbackCopyToClipboard(text) {
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            
            try {
                document.execCommand('copy');
                showCopySuccess();
            } catch (err) {
                alert('Failed to copy to clipboard. Please copy manually: ' + text);
            }
            
            document.body.removeChild(textArea);
        }
        
        // Show copy success feedback
        function showCopySuccess() {
            const copyBtn = document.querySelector('.copy-btn');
            if (copyBtn) {
                const originalHTML = copyBtn.innerHTML;
                const originalBg = copyBtn.style.background;
                
                copyBtn.innerHTML = '<i class="fas fa-check"></i>';
                copyBtn.classList.add('copy-success');
                
                setTimeout(() => {
                    copyBtn.innerHTML = originalHTML;
                    copyBtn.classList.remove('copy-success');
                    copyBtn.style.background = originalBg;
                }, 2000);
            }
        }
        
        // Download QR code function
        function downloadQR(qrUrl, shortCode) {
            const link = document.createElement('a');
            link.href = qrUrl;
            link.download = 'qr-code-' + shortCode + '.png';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }
        
        // Clear results when typing
        document.getElementById('urlInput').addEventListener('input', function() {
            const resultDiv = document.getElementById('result');
            if (resultDiv.innerHTML) {
                resultDiv.innerHTML = '';
            }
        });
        
        document.getElementById('customCodeInput').addEventListener('input', function() {
            const resultDiv = document.getElementById('result');
            if (resultDiv.innerHTML) {
                resultDiv.innerHTML = '';
            }
        });
    </script>
</body>
</html>`
}
