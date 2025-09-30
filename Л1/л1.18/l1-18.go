package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	counter int
	mu      sync.Mutex
}

func (count *Counter) Incremention() {
	count.mu.Lock()
	defer count.mu.Unlock()
	count.counter++

}

func (count *Counter) ReturnValue() int {
	count.mu.Lock()
	defer count.mu.Unlock()
	return count.counter
}

func main() {
	counter := Counter{}
	var wg sync.WaitGroup
	numOfGoroutrines := 1000
	wg.Add(numOfGoroutrines)

	for i := 0; i < numOfGoroutrines; i++ {
		go func() {
			defer wg.Done()
			counter.Incremention()
		}()
	}
	wg.Wait()
	fmt.Println(counter.ReturnValue())
}
