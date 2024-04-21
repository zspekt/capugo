package srv

import (
	"log/slog"
	"net/http"
	"os"
)

// for server
var (
	port    string
	address string
)

func init() {
	slog.Info("running server init func...")

	// loads env vars. make sure you have the .env file in the dir you're running
	// the server from
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
		slog.Info("port env vas was empty. using default", "defPort", port)
	}
	slog.Info("port var has been set", "port", port)

	address = os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = "localhost"
		slog.Info("address env var was empty. using default", "defAddr", address)
	}
	slog.Info("address var has been set", "address", address)

	slog.Info("server init func executed without any errors...")
}

func ReturnServer() *http.Server {
	router := http.NewServeMux()

	setRoutes(router)

	return &http.Server{
		Addr:    address + ":" + port,
		Handler: router,
	}
}
