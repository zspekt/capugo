package handlers

import (
	"log"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
  log.Println("healthCheck handler called...")

  w.Header().Set(key string, value string)
  // I gotta leave. 
}
