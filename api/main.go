package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/3d0c/sample-api/api/handlers"
	"github.com/3d0c/sample-api/api/models"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var (
		listenOn string
	)

	flag.StringVar(&listenOn, "listen-on", ":5560", "listen on")
	flag.Parse()

	if err := models.ConnectDatabase(); err != nil {
		log.Fatalf("Error connecting to database - %s\n", err)
	}

	router := handlers.SetupRouter()

	log.Printf("API handler is listening on %s\n", listenOn)

	log.Fatalln(
		http.ListenAndServe(listenOn, router),
	)
}
