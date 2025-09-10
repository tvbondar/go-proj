// остановка через канал уведомлений
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := make(chan bool)
	go worker(stop)
	time.Sleep(2 * time.Second)
	stop <- true
	time.Sleep(100 * time.Millisecond)
}

func worker(stop chan bool) {
	for {
		select {
		case <-stop:
			fmt.Printf("Worker stopped\n")
			return
		default:
			fmt.Printf("Working\n")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
