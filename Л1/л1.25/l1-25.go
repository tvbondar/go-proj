package main

import (
	"fmt"
	"time"
)

func NewSleep(duration time.Duration) {
	<-time.After(duration)

}

func main() {
	fmt.Println("Sleep начинает работу:", time.Now().Format("15:04:05.00"))
	NewSleep(10 * time.Second)
	fmt.Println("Sleep заканчивает работу: ", time.Now().Format("15:04:05.000"))

}
