package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	r := repository.NewMongoDBUrlEntryStore(client)

	router := &appRouter{
		repo: r,
	}

	mux := http.NewServeMux()

	router.setup(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
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
