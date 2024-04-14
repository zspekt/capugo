package main

import (
	"github.com/zspekt/capugo/src/srv"
	"log"
)

// @title CapuGO API Documentation
// @description Swagger API Documentation
// @BasePath /api/v1
func main() {
	srv := srv.ReturnServer()

	log.Fatal(srv.ListenAndServe())
}
