/*
Multiple functions can read from the same channel until that channel is closed;
this is called fan-out. This provides a way to distribute work amongst a group
of workers to parallelize CPU use and I/O.
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Result struct {
	workerID int // for example only, usually we don't care about the workerID
	url      string
	resp     *http.Response
	err      error
}

func generateURLsChan(ctx context.Context, urls []string) <-chan string {
	urlsChan := make(chan string)
	go func() {
		defer close(urlsChan) // only the writer can close the channel, otherwise with an attempt to write to the closed channel - you can get a panic
		for _, url := range urls {
			select {
			case <-ctx.Done():
				fmt.Println("the generate urls channel is closed")
				return
			case urlsChan <- url:
			}
		}
	}()
	return urlsChan
}

func fetchURLs(ctx context.Context, urlsChan <-chan string, workerID int) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("fetch URLs channel is done")
				return
			case url, ok := <-urlsChan:
				if !ok {
					fmt.Println("urlsChan is closed")
					return
				}
				resp, err := http.Get(url)
				if err != nil {
					out <- Result{
						workerID: workerID,
						err:      fmt.Errorf("worker %d: error fetching %s, err: %w", workerID, url, err),
					}
					continue
				}
				out <- Result{
					workerID: workerID,
					url:      url,
					resp:     resp,
				}
				resp.Body.Close()
			}
		}
	}()

	return out
}

func fetchURLsWithWorkers(ctx context.Context, urls <-chan string, numWorkers int) <-chan Result {
	out := make(chan Result, numWorkers)
	var wg sync.WaitGroup

	for workerID := 0; workerID < numWorkers; workerID++ {
		wg.Add(1)
		go func(id int) { // avoid race conditions
			defer wg.Done()
			for fetchResult := range fetchURLs(ctx, urls, id) {
				out <- fetchResult
			}
		}(workerID)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	start := time.Now()

	urls := GenerateRandomURLs() // generates 1k random urls
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// input
	urlsChan := generateURLsChan(ctx, urls)
	numWorkers := runtime.NumCPU()

	// fan-out
	resultStream := fetchURLsWithWorkers(ctx, urlsChan, numWorkers)

	for result := range resultStream {
		if result.err != nil {
			fmt.Printf("received an error: %s, workerID: %d\n", result.err.Error(), result.workerID)
			continue
		}
		fmt.Printf("received result: %+v\n", result)
	}

	end := time.Since(start)
	fmt.Println("it took: ", end) // it took: 29.589958959s for 1k URLs ---> 6.5x faster than sequential
}
