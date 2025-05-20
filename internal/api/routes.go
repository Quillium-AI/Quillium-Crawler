package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	// Health and system endpoints
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/reload-config", reloadConfigHandler)
	mux.HandleFunc("/version", versionHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
	log.Println("Health check request received")
}

// TODO: Implement reload config handler
func reloadConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
	log.Println("Reload config request received")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"version": "0.1.0",
	}
	json.NewEncoder(w).Encode(response)
}
