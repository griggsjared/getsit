package main

import (
	"fmt"
	"time"

	"github.com/griggsjared/getsit/internal/entity"
	"github.com/griggsjared/getsit/internal/repository"
)

func main() {
	seedUrlEntries(10_000_000)
}

// seed will generate a number of tokens and check for duplicates
func seedUrlEntries(tCount int) {

	genCount := 0
	dupCount := 0

	store := repository.NewInMemoryUrlEntryStore()

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
				percent = (float64(genCount) + float64(dupCount)) / float64(tCount) * 100
			}
			fmt.Printf("\rEntries Seeded: %d, Duplicates (Already Existed): %d, Percent %.2f%%", genCount, dupCount, percent)
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
		url := entity.Url("https://example.com/" + string(entity.NewUrlToken()))
		_, new, err := store.Save(url)

		if err != nil {
			fmt.Println(err)
			break
		}

		if !new {
			dupCount++
		} else {
			genCount++
		}
	}
}
