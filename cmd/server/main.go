package main

import (
	"github.com/lipandr/dist-log/internal/server"
	"log"
)

func main() {
	srv := server.NewHTTPServer(":5000")
	log.Fatal(srv.ListenAndServe())
}
