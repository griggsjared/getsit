package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/griggsjared/getsit/internal/repository"
	"github.com/griggsjared/getsit/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type app struct {
	service *service.Service
	logger  *log.Logger
}

func main() {

	godotenv.Load()

	ctx := context.Background()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		fmt.Println("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	serverAddr := ":" + port
	if host != "" {
		serverAddr = host + serverAddr
	}

	app := &app{
		service: service.New(repository.NewPGXUrlEntryRepository(db)),
		logger:  log.New(os.Stdout, "", log.LstdFlags),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /url-entries", app.createUrlEntryHandler)
	mux.HandleFunc("GET /url-entries/{token}", app.getUrlEntryHandler)
	mux.HandleFunc("GET /healthz", app.healthzHandler)

	fmt.Printf("Starting server on %s\n", serverAddr)

	http.ListenAndServe(serverAddr, app.middlewareStack(mux, app.loggerMiddleware))
}
