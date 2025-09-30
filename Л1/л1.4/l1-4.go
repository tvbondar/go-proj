package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		fmt.Println("Number of workers is not provided")
		cancel()
	}
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Println("Provide a positive number of workers")
		cancel()
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
		<-sigChan
		fmt.Println("\nStoping program.")
		cancel()
	}()

	workerPool(ctx, numWorkers)
}

func workerPool(ctx context.Context, numWorkers int) {

	wg := &sync.WaitGroup{}
	dataChan := make(chan int)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			worker(workerID, dataChan, ctx)
		}(i + 1)
	}

	go func() {
		defer close(dataChan)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				dataChan <- counter
				counter++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	fmt.Printf("Started %d workers\n ", numWorkers)
	<-ctx.Done()
	wg.Wait()
	fmt.Println("\nProgram stopped")

}

func worker(workerID int, dataChan <-chan int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-dataChan:
			if !ok {
				return
			}
			fmt.Printf("\nWorker %d: %d\n", workerID, data)
		}
	}
}
