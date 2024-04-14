package srv

import (
	"net/http"

	"github.com/zspekt/capugo/src/handlers"
)

// TODO:
//       - subrouting

func setRoutes(router *http.ServeMux) {
	// do note that ONLY ONE SPACE is allowed between the http method
	// and the endpoint.  ↓
	router.HandleFunc("GET /api/v1/health", handlers.HealthCheck)
	router.HandleFunc("POST /api/v1/auth/token/generate", handlers.GenerateToken)
	router.HandleFunc("POST /api/v1/auth/token/validate", handlers.ValidateToken)
}
