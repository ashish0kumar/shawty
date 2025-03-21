package main

import (
	"ashish0kumar/shawty/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file. Using system environment variables instead.")
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
		tmpl.Execute(writer, nil)
	})

	// Shortens the provided URL, store it and return it to our UI
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {

		// Get the URL from the request
		url := req.FormValue("url")
		fmt.Println("Payload: ", url)

		// Shorten the URL
		shortURL := utils.GetShortCode()
		fullShortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortURL)

		fmt.Printf("Shortened URL: %s\n", shortURL)

		// Set the key in Redis
		utils.SetKey(&ctx, dbClient, shortURL, url, 0)

		// Return the response to the UI rendered with HTMX
		fmt.Fprintf(writer, `<p class="mt-4 text-blue-600"><a href="/r/%s" class="underline">%s</a></p>`, shortURL, fullShortURL)
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

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
