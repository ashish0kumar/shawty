package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/ashish0kumar/shawty/utils"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

// GetEnv retrieves an environment variable or returns a default value if not found
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file. Using system environment variables instead.")
	}

	// Get base URL from environment variable or use localhost as default
	baseURL := GetEnv("BASE_URL", "http://localhost:8080")

	// Create the redisDB connection
	dbClient := utils.NewRedisClient()
	if dbClient == nil {
		fmt.Println("Error connecting to Redis")
		return
	}

	// Serves the UI
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(writer, nil)
	})

	// Shortens the provided URL, store it and return it to our UI
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {
		// Get the URL from the request
		url := req.FormValue("url")
		fmt.Println("Payload: ", url)

		// Shorten the URL
		shortURL := utils.GetShortCode()
		fullShortURL := fmt.Sprintf("%s/r/%s", baseURL, shortURL)
		fmt.Printf("Shortened URL: %s\n", shortURL)

		// Set the key in Redis
		utils.SetKey(&ctx, dbClient, shortURL, url, 0)

		// Return the response to the UI rendered with HTMX
		resultHTML := `
            <div class="bg-green-50 border border-green-200 rounded-lg p-4">
                <div class="flex items-center mb-2">
                    <i class="fas fa-check-circle text-green-500 mr-2"></i>
                    <p class="font-medium text-green-800">URL Shortened Successfully!</p>
                </div>
                <div class="flex items-center justify-between bg-white rounded border p-3 mt-2">
                    <a href="/r/%s" class="text-blue-600 hover:text-blue-800 truncate max-w-[70%%]" target="_blank">%s</a>
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
