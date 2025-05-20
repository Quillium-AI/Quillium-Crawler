package api

import (
	"log"
	"net/http"
)

func StartServer(addr string) error {

	mux := http.NewServeMux()
	SetupRoutes(mux)

	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, mux)
}
