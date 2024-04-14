package main

import (
	"github.com/joho/godotenv"
	"github.com/zspekt/capugo/src/srv"
	"log"
)

func main() {
	srv := srv.ReturnServer()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.ListenAndServe())
}
