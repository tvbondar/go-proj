package main

import "fmt"

func main() {
	Slice := make([]int, 5, 10)
	for j := range Slice {
		Slice[j] = (j + 1) * 5
	}
	fmt.Println(Slice)
	i := 3
	copy(Slice[i:], Slice[i+1:])
	Slice = Slice[:len(Slice)-1]
	fmt.Println(Slice)
}
