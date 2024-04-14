package srv

import (
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zspekt/capugo/docs"
	"github.com/zspekt/capugo/src/handlers"
	"net/http"
)

func setRoutes(router *http.ServeMux) {
	// do note that ONLY ONE SPACE is allowed between the http method
	// and the endpoint.  ↓
	router.HandleFunc("GET /api/v1/health", handlers.HealthCheck)
	router.HandleFunc("POST /api/v1/login", handlers.LoginHandler)
	router.HandleFunc("GET /documentation/", httpSwagger.WrapHandler)
}
