package main

import (
	"fmt"
)

func Numbers(arr []int) <-chan int {
	c := make(chan int)
	go func() {
		for _, num := range arr {
			c <- num
		}
		close(c)
	}()
	return c
}

func Mul(c <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range c {
			out <- v * 2
		}
		close(out)
	}()
	return out
}

func main() {
	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numChan := Numbers(num)
	mulChan := Mul(numChan)
	for result := range mulChan {
		fmt.Println(result)
	}
}
