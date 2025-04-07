package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	url  string
	resp *http.Response
	err  error
}

func fetchURL(ctx context.Context, url string) Result {
	resp, err := http.Get(url)
	if err != nil {
		return Result{
			url: url,
			err: fmt.Errorf("error fetching %s: %v", url, err),
		}
	}
	defer resp.Body.Close()
	return Result{
		url:  url,
		resp: resp,
	}
}

func main() {
	start := time.Now()
	urls := GenerateRandomURLs()

	ctx := context.Background()

	for _, url := range urls {
		result := fetchURL(ctx, url)
		if result.err != nil {
			fmt.Printf("Error: %s\n", result.err)
		} else {
			fmt.Printf("Successfully fetched: %s, Status Code: %d\n", result.url, result.resp.StatusCode)
		}
	}

	end := time.Since(start)
	fmt.Println("it took: ", end) // it took: 3m12.897654125s for 1k URLs
}
