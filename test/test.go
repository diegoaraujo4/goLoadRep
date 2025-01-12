package test

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func RunLoadTest(url string, totalRequests, concurrency int) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var statusCount = make(map[int]int)
	var totalRequestsProcessed int
	if url == "" {
		return fmt.Errorf("URL must be provided")
	}
	if totalRequests <= 0 {
		return fmt.Errorf("total requests must be greater than 0")
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	start := time.Now()

	bar := progressbar.New(totalRequests)

	makeRequest := func() {
		defer wg.Done()
		resp, err := http.Get(url)
		if err != nil {
			log.Default().Println("Error making request:", err)
			return
		}
		mu.Lock()
		statusCount[resp.StatusCode]++
		totalRequestsProcessed++
		mu.Unlock()
		bar.Add(1)
	}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go makeRequest()
		if (i+1)%concurrency == 0 {
			wg.Wait()
		}
	}

	wg.Wait()

	elapsedTime := time.Since(start)
	fmt.Printf("\nLoad Test Report:\n")
	fmt.Printf("Total execution time: %s\n", elapsedTime)
	fmt.Printf("Total requests made: %d\n", totalRequestsProcessed)
	fmt.Printf("Requests with HTTP status 200: %d\n", statusCount[200])
	for status, count := range statusCount {
		if status != 200 {
			fmt.Printf("Requests with HTTP status %d: %d\n", status, count)
		}
	}

	return nil
}
