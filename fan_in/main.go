/*
A function can read from multiple inputs and proceed until all are closed by
multiplexing the input channels onto a single channel that’s closed when all the
inputs are closed. This is called fan-in.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generateTickets(airline string) <-chan string { // Lexical confinement - restrict a return channel to be read-only
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 0; i < 10; i++ {
			ticketNumber := rand.Intn(1000) + 1
			out <- fmt.Sprintf("%s Ticket #%d", airline, ticketNumber)
			time.Sleep(time.Second)
		}
	}()
	return out
}

func MergeChannels[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup // waiting unit all input channels are closed

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for val := range ch {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out) // closing when all input channels are closed
	}()

	return out
}

func main() {
	for ticket := range MergeChannels(
		generateTickets("Ryanair"),
		generateTickets("Lufthansa"),
		generateTickets("EasyJet"),
	) {
		fmt.Println(ticket)
	}
}
