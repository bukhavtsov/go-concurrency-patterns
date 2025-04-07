package main

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomURLs() []string {
	var urls []string
	rand.New(rand.NewSource(time.Now().UnixNano()))
	domains := []string{"example.com", "example.org", "example.net", "example.edu", "test.com"}

	for i := 0; i < 1000; i++ {
		protocol := "http"
		if rand.Intn(2) == 1 {
			protocol = "https"
		}

		domain := domains[rand.Intn(len(domains))]

		path := fmt.Sprintf("/%d/%d/%d", rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))

		url := fmt.Sprintf("%s://%s%s", protocol, domain, path)
		urls = append(urls, url)
	}
	return urls
}
