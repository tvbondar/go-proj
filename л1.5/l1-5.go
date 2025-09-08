package main

import (
	"fmt"
	"time"
)

func main() {
	var N int
	fmt.Scan(&N)
	fmt.Printf("Программа будет работать %d секунд...\n", N)

	dataChan := make(chan int)

	go func() {
		counter := 0
		for {
			dataChan <- counter
			counter++
			time.Sleep(300 * time.Millisecond)
		}
	}()

	timeout := time.After(time.Duration(N) * time.Second)

	for {
		select {
		case data := <-dataChan:
			fmt.Printf("Получено: %d\n", data)

		case <-timeout:
			fmt.Printf("Прошло %d секунд. Завершение работы.\n", N)
			return
		}
	}
}
