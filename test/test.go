package test

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func RunLoadTest(url string, totalRequests, concurrency int) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var statusCount = make(map[int]int)
	var totalRequestsProcessed int

	start := time.Now()

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
	fmt.Printf("Load Test Report:\n")
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
