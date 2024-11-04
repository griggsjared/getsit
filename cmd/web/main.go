package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/griggsjared/getsit/internal/repository"
	"github.com/griggsjared/getsit/internal/service"
)

type app struct {
	service *service.Service
	logger  *log.Logger
	session *sessions.CookieStore
}

func main() {
	godotenv.Load()

	ctx := context.Background()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		fmt.Println("DATABASE_DSN is not set")
		os.Exit(1)
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		fmt.Println("SESSION_SECRET is required")
		os.Exit(1)
	}

	csrfSecret := os.Getenv("CSRF_TOKEN")
	if csrfSecret == "" {
		fmt.Println("CSRF_TOKEN is not set")
		os.Exit(1)
	}

	app := &app{
		service: service.New(repository.NewPGXUrlEntryRepository(db)),
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		session: sessions.NewCookieStore([]byte(sessionSecret)),
	}

	csrfMiddleware := csrf.Protect(
		[]byte(csrfSecret),
		csrf.Secure(false),
		csrf.CookieName("CSRF-TOKEN"),
		csrf.RequestHeader("X-CSRF-TOKEN"),
		csrf.FieldName("csrf_token"),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.middlewareStackFunc(app.homepageHandler, csrfMiddleware))
	mux.HandleFunc("POST /create", app.middlewareStackFunc(app.createHandler, csrfMiddleware))
	mux.HandleFunc("GET /i/{token}", app.infoHandler)
	mux.HandleFunc("GET /{token}", app.redirectHandler)
	mux.HandleFunc("/", app.notFoundHandler)

	fileServer := http.FileServer(http.Dir("web/public"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	serverAddr := ":" + port
	if host != "" {
		serverAddr = host + serverAddr
	}

	http.ListenAndServe(serverAddr, app.middlewareStack(mux, app.loggerMiddleware, app.gzipMiddleware))

	fmt.Println("Server started on", serverAddr)
}
