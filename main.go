package main

import "urlShortenerService/internal/transport/http"

func main() {
	// Initialize the HTTP router
	router := http.NewRouter()

	// Start the service on the port 8080
	router.Run(":8080")
}
