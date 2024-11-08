package main

import (
	"fmt"
	"net/http"
	"os"
)

type app struct{}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	serverAddr := ":" + port
	if host != "" {
		serverAddr = host + serverAddr
	}

	app := &app{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.rootHandler)
	mux.HandleFunc("GET /healthz", app.healthzHandler)

	fmt.Printf("Starting server on %s\n", serverAddr)

	http.ListenAndServe(serverAddr, app.middlewareStack(mux))
}
