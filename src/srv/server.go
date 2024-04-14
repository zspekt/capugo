package srv

import (
	"log"
	"net/http"
	"os"

	"github.com/zspekt/capugo/src/handlers"
)

// GTG again but..

// TODO: replace godotenv with standard library os.Getenv

// for server
var (
	port    string
	address string
)

func init() {
	log.Println("Running server init func...")

	// loads env vars. make sure you have the .env file in the dir you're running
	// the server from
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
		log.Println("port env vas was empty. Using default (8080)...")
	}
	log.Printf("port var has been set %v\n", port)

	address = os.Getenv("ADDRESS")
	if len(port) == 0 {
		address = "localhost"
		log.Println("address env var was empty. Using default (localhost)")
	}
	log.Printf("address var has been set %v\n", address)

	log.Println("Server init func executed without any errors...")
}

func ReturnServer() *http.Server {
	router := http.NewServeMux()

	// do note that ONLY ONE SPACE is allowed between the http method
	// and the endpoint.  ↓
	router.HandleFunc("GET /health", handlers.HealthCheck)

	apiv1 := http.NewServeMux()
	// any request that hits on the /api/v1/ path, will simply get redirected,
	// stripping the /api/v1/
	apiv1.Handle("/api/v1/", http.StripPrefix("/api/v1/", router))

	return &http.Server{
		Addr:    address + ":" + port,
		Handler: apiv1,
	}
}
