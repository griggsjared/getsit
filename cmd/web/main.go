package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/griggsjared/getsit/internal"
	"github.com/griggsjared/getsit/internal/repository"
)

type app struct {
	service *internal.Service
	logger  *log.Logger
}

const database = "POSTGRES"

func main() {
	godotenv.Load()

	ctx := context.Background()

	var r internal.UrlEntryRepository

	if database == "MONGO" {
		dsn := os.Getenv("MONGODB_DSN")
		if dsn == "" {
			fmt.Println("MONGODB_DSN is not set")
			os.Exit(1)
		}

		client, err := setupMongoDB(ctx, dsn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				fmt.Println(err)
			}
		}()
		r = repository.NewMongoDBUrlEntryStore(client)

	} else {

		dsn := os.Getenv("POSTGRES_DSN")
		if dsn == "" {
			fmt.Println("POSTGRES_DSN is not set")
			os.Exit(1)
		}

		db, err := pgxpool.New(ctx, dsn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer db.Close()

		r = repository.NewPGXUrlEntryStore(db)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := &app{
		service: internal.NewService(r),
		logger:  logger,
	}

	token := os.Getenv("CSRF_TOKEN")
	if token == "" {
		fmt.Println("CSRF_TOKEN is not set")
		os.Exit(1)
	}

	csrf := csrf.Protect(
		[]byte(token),
		csrf.Secure(false),
		csrf.CookieName("CSRF-TOKEN"),
		csrf.RequestHeader("X-CSRF-TOKEN"),
		csrf.FieldName("csrf_token"),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.middlewareStackFunc(app.homepageHandler, csrf))
	mux.HandleFunc("POST /create", app.middlewareStackFunc(app.createHandler, csrf))
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

	http.ListenAndServe(serverAddr, app.middlewareStack(mux, app.loggerMiddleware))

	fmt.Println("Server started on", serverAddr)
}

// setupMongoDB will setup the connection to the MongoDB database
func setupMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {

	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
