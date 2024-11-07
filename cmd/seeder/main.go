package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/griggsjared/getsit/internal/repository"
	"github.com/griggsjared/getsit/internal/service"
)

func main() {

	var tCount int
	var wCount int
	var fresh bool

	flag.IntVar(&tCount, "n", 1000, "number of url entries to seed")
	flag.IntVar(&wCount, "w", 4, "number of workers")
	flag.BoolVar(&fresh, "f", false, "truncate url_entries table before seeding")
	flag.Parse()

	ctx := context.Background()

	godotenv.Load()

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

	r := repository.NewPGXUrlEntryRepository(db)

	service := service.New(r)

	if fresh {
		db.Exec(ctx, "TRUNCATE url_entries")
		fmt.Println("Truncated url_entries table")
	}

	SeedUrlEntries(ctx, tCount, wCount, service)
}

// seed will generate a number of tokens and check for duplicates
func SeedUrlEntries(ctx context.Context, tCount int, wCount int, s *service.Service) {

	genCount := 0

	timeStart := time.Now()

	fmt.Println("Seeding", tCount, "url entries with random tokens", "using", wCount, "workers")

	fChan := make(chan bool)

	go func() {
		for {
			var percent float64
			var finished bool
			select {
			case finished = <-fChan:
				if finished {
					genCount = tCount
					percent = 100.0
				}
			default:
				percent = float64(genCount) / float64(tCount) * 100
			}
			fmt.Printf("\rEntries Seeded: %d, Percent %.2f%%", genCount, percent)
			if finished {
				fmt.Println("\nTime taken: ", time.Since(timeStart))
				close(fChan)
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	fChan <- false
	defer func() {
		fChan <- true
	}()

	wg := sync.WaitGroup{}

	perWorker := tCount / wCount

	for i := 0; i < wCount; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < perWorker; j++ {

				rand, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
				if err != nil {
					fmt.Println(err)
					break
				}

				s.SaveUrl(ctx, &service.SaveUrlInput{
					Url: "https://example.com/" + fmt.Sprintf("%d", rand),
				})

				genCount++
			}
		}()
	}

	wg.Wait()
}
