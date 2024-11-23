package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/griggsjared/getsit/internal/url"
	"github.com/griggsjared/getsit/internal/url/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type app struct {
	urlService *url.Service
	logger     *slog.Logger
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
		urlService: url.NewService(repository.NewPGXUrlEntryRepository(db)),
		logger:     slog.Default().With(slog.String("service", "getsit-api")),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /url-entries", app.createUrlEntryHandler)
	mux.HandleFunc("GET /url-entries/{token}", app.getUrlEntryHandler)
	mux.HandleFunc("GET /healthz", app.healthzHandler)

	fmt.Printf("Starting server on %s\n", serverAddr)
	http.ListenAndServe(serverAddr, app.middlewareStack(mux, app.loggerMiddleware))
}
