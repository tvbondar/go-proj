// выход из горутины через runtime.Goexit() (завершает только текущую горутину)
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	go worker()
	time.Sleep(2 * time.Second)
}

func worker() {
	defer fmt.Println("Worker killed")
	fmt.Println("Worker started")
	time.Sleep(1 * time.Second)
	fmt.Println("About to exit")
	runtime.Goexit()
}
