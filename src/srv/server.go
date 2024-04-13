package srv

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env %v\n", err)
		return
	}

	port = os.Getenv("PORT")
	log.Printf("port var has been set %v\n", port)

	address = os.Getenv("ADDRESS")
	log.Printf("address var has been set %v\n", address)

	log.Println("Server init func executed without any errors...")
}

func ReturnServer() *http.Server {
	router := http.NewServeMux()

	// do note that ONLY ONE SPACE is allowed between the http method
	// and the endpoint.  ↓
	router.HandleFunc("GET /health", handlers.HealthCheck)

	return &http.Server{
		Addr:    address + port,
		Handler: router,
	}
}
