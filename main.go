package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	cache := NewKVS(10 * time.Second)

	app := NewServer(cache)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
