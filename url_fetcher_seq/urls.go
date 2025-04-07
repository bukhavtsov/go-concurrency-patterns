package main

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomURLs() []string {
	var urls []string
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// List of sample domains for random URL generation
	domains := []string{"example.com", "example.org", "example.net", "example.edu", "test.com"}

	for i := 0; i < 1000; i++ {
		// Randomly choose a protocol (http or https)
		protocol := "http"
		if rand.Intn(2) == 1 {
			protocol = "https"
		}

		// Randomly select a domain from the list
		domain := domains[rand.Intn(len(domains))]

		// Randomly create a path
		path := fmt.Sprintf("/%d/%d/%d", rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))

		// Construct the URL and append to the slice
		url := fmt.Sprintf("%s://%s%s", protocol, domain, path)
		urls = append(urls, url)
	}
	return urls
}
