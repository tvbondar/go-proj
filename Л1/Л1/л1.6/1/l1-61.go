// завершение по условию
package main

import (
	"fmt"
	"time"
)

//в этом задании рассматриваются способы завершения
//горутин, поэтому можем сделать имитацию работы

func main() {
	done := make(chan bool)
	go worker(1, done)
	<-done
}

func worker(id int, done chan bool) {
	fmt.Printf("Worker %d started\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d completed\n", id)
	done <- true
}
