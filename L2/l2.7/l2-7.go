package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Функция создает канал и горутину, которая отправляет числа со случайными задержками
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v // Отправляем значения в выходной канал
				} else {
					a = nil //закрываем канал
				}
			case v, ok := <-b:
				if ok {
					c <- v // отправляем значения
				} else {
					b = nil // закрываем канал
				}
			}
			//проверяем, что a и b закрыты
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

// запись из двух каналов в один одновременно, с задержками
func main() {
	//Читаются данные из двух каналов и записываются в канал c
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}
