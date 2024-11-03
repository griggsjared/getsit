package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/griggsjared/getsit/internal/entity"
	"github.com/griggsjared/getsit/internal/repository"
	"github.com/griggsjared/getsit/internal/service"
)

func main() {

	godotenv.Load()

	ctx := context.Background()

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

	r := repository.NewPGXUrlEntryRepository(db)

	service := service.New(r)

	db.Exec(ctx, "TRUNCATE url_entries")

	SeedUrlEntries(ctx, 100000, service)
}

// seed will generate a number of tokens and check for duplicates
func SeedUrlEntries(ctx context.Context, tCount int, s *service.Service) {

	genCount := 0

	timeStart := time.Now()

	fmt.Println("Seeding", tCount, "url entries with random tokens")

	fChan := make(chan bool)

	go func() {
		for {
			var percent float64
			var finished bool
			select {
			case finished = <-fChan:
				if finished {
					percent = 100.0
				}
			default:
				percent = float64(genCount) / float64(tCount) * 100
			}
			fmt.Printf("\rEntries Seeded: %d, Percent %.2f%%", genCount, percent)
			time.Sleep(10 * time.Millisecond)
			if finished {
				fmt.Println("\nTime taken: ", time.Since(timeStart))
				close(fChan)
				break
			}
		}
	}()

	fChan <- false
	defer func() {
		fChan <- true
	}()

	for i := 0; i < tCount; i++ {

		e, err := s.SaveUrl(ctx, &service.SaveUrlInput{
			Url: "https://example.com/" + entity.NewUrlToken().String(),
		})

		if err != nil {
			fmt.Println(err)
			break
		}

		genCount++

		s.VisitUrl(ctx, &service.VisitUrlInput{
			Token: e.Token.String(),
		})
	}
}
