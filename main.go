package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/griggsjared/getsit/internal"
	"github.com/griggsjared/getsit/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	ctx := context.Background()

	dbUri := os.Getenv("MONGODB_URI")
	if dbUri == "" {
		fmt.Println("MONGODB_URI is not set")
		os.Exit(1)
	}

	client, err := setupMongoDB(ctx, dbUri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Setting up services")

	r := repository.NewMongoDBUrlEntryStore(client)

	handler := &appHandler{
		service: internal.NewService(r),
	}

	mux := http.NewServeMux()

	handler.setup(mux)

	mux.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(
				http.Dir("web/public"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	serveAddr := ":" + port
	if host != "" {
		serveAddr = host + serveAddr
	}

	server := &http.Server{
		Addr:    serveAddr,
		Handler: mux,
	}

	server.ListenAndServe()

	fmt.Println("Server started on", serveAddr)
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
