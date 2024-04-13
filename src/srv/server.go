package srv

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	port    string
	address string
)

func init() {
	log.Println("Running server nit func...")

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

	// router.HandleFunc("GET /health", handler func(http.ResponseWriter, *http.Request))

	return &http.Server{
		Addr:    address + port,
		Handler: router,
	}
}
