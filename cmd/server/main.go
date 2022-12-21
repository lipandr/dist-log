package main

import (
	"log"

	"github.com/lipandr/dist-log/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":5000")
	log.Fatal(srv.ListenAndServe())
}
