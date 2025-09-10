// выход через закрытие канала (если горутина читает из него)
package main

import (
	"fmt"
	"time"
)

func main() {
	tasks := make(chan int, 5)
	done := make(chan bool)
	go worker(tasks, done)
	for i := 1; i <= 3; i++ {
		tasks <- i
	}
	close(tasks)
	<-done
}

func worker(tasks <-chan int, done chan bool) {
	for task := range tasks {
		fmt.Printf("Processing task %d\n", task)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("All tasks processed")
	done <- true
}
