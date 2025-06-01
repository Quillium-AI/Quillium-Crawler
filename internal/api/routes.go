package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CorsMiddleware adds CORS headers to allow all origins (*)
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// WrapWithCors wraps an http.HandlerFunc with CORS middleware
func WrapWithCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the original handler
		h(w, r)
	}
}

func SetupRoutes(mux *http.ServeMux) {
	// Health and system endpoints
	mux.HandleFunc("/livez", WrapWithCors(livezHandler))
	mux.HandleFunc("/readyz", WrapWithCors(readyzHandler))
	mux.HandleFunc("/version", WrapWithCors(versionHandler))

	// Metrics endpoint - using CorsMiddleware for non-HandlerFunc
	mux.Handle("/metrics", CorsMiddleware(promhttp.Handler()))
}

func livezHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
	log.Println("Live check request received")
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
	log.Println("Ready check request received")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"version": "0.1.0",
	}
	json.NewEncoder(w).Encode(response)
}
