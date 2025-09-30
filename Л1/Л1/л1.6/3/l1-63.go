// остановка через контекст
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout((context.Background()), 2*time.Second)
	defer cancel()
	go worker(ctx, 1)
	go worker(ctx, 2)
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)

}

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopped\n", id)
			return
		default:
			fmt.Printf("Worker %d working\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
