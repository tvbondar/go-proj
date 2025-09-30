package main

import (
	"fmt"
	"math/rand"
)

func CreateSlice() []int {
	sizeOfSlice := 20
	slice := make([]int, sizeOfSlice)
	for i := 0; i < sizeOfSlice; i++ {
		slice[i] = rand.Intn(sizeOfSlice)
	}
	return slice
}

func quickSort(slice []int) []int {
	if len(slice) < 2 {
		return slice
	}
	pivot := slice[0]
	var less, greater []int
	for _, num := range slice[1:] {
		if num <= pivot {
			less = append(less, num)
		} else {
			greater = append(greater, num)
		}
	}
	sortedSlice := append(quickSort(less), pivot)
	sortedSlice = append(sortedSlice, quickSort(greater)...)
	return sortedSlice
}

func main() {
	newSlice := CreateSlice()
	fmt.Println(newSlice)
	sortedSlice := quickSort(newSlice)
	fmt.Println(sortedSlice)
}
