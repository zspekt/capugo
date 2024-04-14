package main

import (
	"log"

	"github.com/zspekt/capugo/src/srv"
)

func main() {
	srv := srv.ReturnServer()

	log.Fatal(srv.ListenAndServe())
}
