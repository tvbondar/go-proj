package main

import "fmt"

func CheckType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer. %d\n", v)
	case string:
		fmt.Printf("String. %s\n", v)
	case bool:
		fmt.Printf("Bool. %t\n", v)
	case chan int, chan string, chan bool, chan interface{}:
		fmt.Printf("Channel. %v\n", v)
	default:
		fmt.Printf("Unknown.\n")
	}
}

func main() {
	CheckType(5)
	CheckType("Hi")
	CheckType(false)
	CheckType(3.345)
	CheckType(make(chan int))
	CheckType(make(chan string))
	CheckType(make(chan bool))
	CheckType(make(chan interface{}))
	CheckType(1.2344)
}
