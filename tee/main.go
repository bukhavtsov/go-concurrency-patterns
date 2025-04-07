/*
The fan-out pattern splits one input channel into multiple channels for parallel
processing, while the tee pattern creates two identical channels with the same
data as the original.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Amount float64
}

func Tee[T any](source <-chan T, sinks ...chan<- T) { // Sink: A component that receives data from a stream and processes or stores it.
	var waitSinks sync.WaitGroup

	for data := range source {
		waitSinks.Add(1)
		go func(data T) { // Pass data explicitly to avoid race conditions
			defer waitSinks.Done()
			for _, sink := range sinks {
				sink <- data
			}
		}(data)
	}

	waitSinks.Wait() // Wait for sinks to process all sources
	for _, sink := range sinks {
		close(sink)
	}
}

func logOrder(orderChan <-chan Order) {
	for order := range orderChan {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500))) // Simulate processing delay
		fmt.Printf("[Logs] Order %+v logged\n", order)
	}
}

func trackOrder(orderChan <-chan Order) {
	for order := range orderChan {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		fmt.Printf("[Metrics] Order %+v tracked \n", order)
	}
}

func share(orderChan <-chan Order) {
	for order := range orderChan {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		fmt.Printf("[Analytics] Order %+v shared with analytics team\n", order)
	}
}

func main() {
	orderSource := make(chan Order)

	loggerChan := make(chan Order)
	metricsChan := make(chan Order)
	distributorChan := make(chan Order)

	go Tee[Order](orderSource, loggerChan, metricsChan, distributorChan)

	// Start order processing in separate goroutines
	go logOrder(loggerChan)
	go trackOrder(metricsChan)
	go share(distributorChan)

	go func() {
		for i := 1; i <= 5; i++ {
			order := Order{ID: i, Amount: float64(rand.Intn(100) + 1)}
			fmt.Printf("[New Order] Received Order %d\n", order.ID)
			orderSource <- order // generating orders
			time.Sleep(time.Millisecond * 300)
		}
	}()

	time.Sleep(time.Second * 10)
}
