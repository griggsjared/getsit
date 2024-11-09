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
	"github.com/griggsjared/getsit/web"
)

type app struct {
	service *service.Service
	logger  *log.Logger
	session *sessions.CookieStore
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

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		fmt.Println("SESSION_SECRET is required")
		os.Exit(1)
	}

	csrfSecret := os.Getenv("CSRF_SECRET")
	if csrfSecret == "" {
		fmt.Println("CSRF_SECRET is not set")
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
		csrf.ErrorHandler(http.HandlerFunc(app.tokenMismatchHandler)),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.middlewareStackFunc(app.homepageHandler, csrfMiddleware, app.templateColorMiddleware))
	mux.HandleFunc("POST /create", app.middlewareStackFunc(app.createHandler, csrfMiddleware))
	mux.HandleFunc("GET /i/{token}", app.middlewareStackFunc(app.infoHandler, app.templateColorMiddleware))
	mux.HandleFunc("GET /{token}", app.redirectHandler)
	mux.HandleFunc("GET /healthz", app.healthzHandler)
	mux.HandleFunc("/", app.notFoundHandler)

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(web.AssetsFS())))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	serverAddr := ":" + port
	if host != "" {
		serverAddr = host + serverAddr
	}

	fmt.Println("Starting server on", serverAddr)

	http.ListenAndServe(serverAddr, app.middlewareStack(mux, app.loggerMiddleware))
}
