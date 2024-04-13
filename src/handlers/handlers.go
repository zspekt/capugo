package handlers

import (
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("healthCheck handler called...")

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
