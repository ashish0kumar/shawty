package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/ashish0kumar/shawty/utils"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var ctx = context.Background()

// retrieves an env variable or returns a default value if not found
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// structure to pass data to templates
type TemplateData struct {
	BaseURL string
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file. Using system environment variables instead.")
	}

	// Initialize safe browsing
	if err := utils.InitSafeBrowsing(os.Getenv("SAFE_BROWSING_API_KEY")); err != nil {
		log.Printf("Failed to initialize Safe Browsing: %v", err)
	}

	// Get base URL from env variable or use localhost as default
	baseURL := GetEnv("BASE_URL", "http://localhost:8080")

	// Create template data
	templateData := TemplateData{
		BaseURL: baseURL,
	}

	// Create the redisDB connection
	dbClient := utils.NewRedisClient()
	if dbClient == nil {
		fmt.Println("Error connecting to Redis")
		return
	}

	// Serves the UI
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(writer, templateData)
	})

	var limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 req/sec

	// Shortens the provided URL, store it and return it to our UI
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {

		// Rate limiting
		if !limiter.Allow() {
			writer.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(writer, `
				<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
					<div class="flex items-center mb-2">
						<i class="fas fa-exclamation-triangle text-yellow-500 mr-2"></i>
						<p class="font-medium text-yellow-800">Too Many Requests</p>
					</div>
					<p class="text-gray-600 text-sm mt-2">Please wait a moment before shortening another URL.</p>
				</div>
    		`)
			return
		}

		// Get the URL from the request
		url := req.FormValue("url")
		fmt.Println("Payload: ", url)

		// Validate the URL
		if err := utils.ValidateURL(url); err != nil {
			errorHTML := `
				<div class="bg-red-50 border border-red-200 rounded-lg p-4">
					<div class="flex items-center mb-2">
						<i class="fas fa-exclamation-circle text-red-500 mr-2"></i>
						<p class="font-medium text-red-800">Invalid URL</p>
					</div>
					<p class="text-gray-600 text-sm mt-2">%s</p>
				</div>
			`
			fmt.Fprintf(writer, errorHTML, err.Error())
			return
		}

		// Check if URL is safe
		isSafe, err := utils.CheckURLSafety(url)
		if err != nil {
			errorHTML := `
				<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
					<div class="flex items-center mb-2">
						<i class="fas fa-exclamation-triangle text-yellow-500 mr-2"></i>
						<p class="font-medium text-yellow-800">Security Check Failed</p>
					</div>
					<p class="text-gray-600 text-sm mt-2">Could not verify URL safety: %s</p>
				</div>
			`
			fmt.Fprintf(writer, errorHTML, err.Error())
			return
		}

		if !isSafe {
			errorHTML := `
				<div class="bg-red-50 border border-red-200 rounded-lg p-4">
					<div class="flex items-center mb-2">
						<i class="fas fa-ban text-red-500 mr-2"></i>
						<p class="font-medium text-red-800">Malicious URL Detected</p>
					</div>
					<p class="text-gray-600 text-sm mt-2">This URL has been identified as malicious by Google Safe Browsing.</p>
				</div>
			`
			fmt.Fprint(writer, errorHTML)
			return
		}

		// Shortening logic
		var shortURL string

		// Check if this URL has already been shortened
		existingShortURL, err := utils.GetExistingShortURL(&ctx, dbClient, url)
		if err != nil {
			fmt.Printf("Error checking existing URL: %v\n", err)
		}

		if existingShortURL != "" {
			// Use the existing short URL
			shortURL = existingShortURL
			fmt.Printf("Using existing short URL: %s\n", shortURL)
		} else {
			// Create a new short URL
			shortURL = utils.GetShortCode()
			fmt.Printf("Generated new short URL: %s\n", shortURL)

			// Store both mappings
			utils.SetKey(&ctx, dbClient, shortURL, url, 24)
		}

		fullShortURL := fmt.Sprintf("%s/r/%s", baseURL, shortURL)

		// Return the response to the UI rendered with HTMX
		resultHTML := `
			<div class="bg-green-50 border border-green-200 rounded-lg p-4">
				<div class="flex items-center mb-2">
					<i class="fas fa-check-circle text-green-500 mr-2"></i>
					<p class="font-medium text-green-800">URL Shortened Successfully!</p>
				</div>
				<div class="flex items-center justify-between bg-white rounded border p-3 mt-2">
					<a href="/r/%s" class="text-blue-600 text-sm hover:text-blue-800 truncate max-w-[90%%]" target="_blank">%s</a>
					<button class="copy-btn bg-gray-100 hover:bg-gray-200 text-gray-800 px-3 py-1 rounded text-sm transition" data-clipboard="%s">
						<i class="fa-regular fa-copy text-gray-600"></i>
					</button>
				</div>
				<p class="text-gray-600 text-sm mt-3">Click the link to visit or copy the link to clipboard!</p>
			</div>
		`
		fmt.Fprintf(writer, resultHTML, shortURL, fullShortURL, fullShortURL)
	})

	// Redirects to the long URL based on the short url
	http.HandleFunc("/r/{code}", func(writer http.ResponseWriter, req *http.Request) {
		key := req.PathValue("code")
		if key == "" {
			http.Error(writer, "Invalid URL", http.StatusBadRequest)
			return
		}
		longURL, err := utils.GetLongURL(&ctx, dbClient, key)
		if err != nil {
			http.Error(writer, "Short URL not found", http.StatusNotFound)
			return
		}
		http.Redirect(writer, req, longURL, http.StatusPermanentRedirect)
	})

	// Get port from environment variable or use 8080 as default
	port := GetEnv("PORT", "8080")

	fmt.Printf("Server running on %s\n", baseURL)
	http.ListenAndServe(":"+port, nil)
}
