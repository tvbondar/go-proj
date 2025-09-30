package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("input error: No arguments ")
		return
	}
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Println("Input error: Argument is equal or lesser than 0")
		return
	}
	workerPool(numWorkers)
}

func workerPool(numWorkers int) {

	wg := &sync.WaitGroup{}
	dataChan := make(chan int)
	stopChan := make(chan struct{})

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, dataChan, stopChan)
		}(i + 1)
	}

	go func() {
		counter := 1
		for {
			select {
			case <-stopChan:
				return
			default:
				dataChan <- counter
				counter++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	fmt.Printf("Started %d workers ", numWorkers)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
	<-sigChan

	close(stopChan)
	wg.Wait()
	fmt.Println("\nProgram stopped")

}

func worker(workerID int, dataChan <-chan int, stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		case data, ok := <-dataChan:
			if !ok {
				return
			}
			fmt.Printf("Worker %d: %d\n", workerID, data)
		}
	}
}
