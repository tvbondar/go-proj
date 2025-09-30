//Алгоритм бинарного поиска работает только на отсортированных данных.
// Поэтому для начала создадим массив  отсортируем его.
//Для простоты используем код из л1.16

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

func binarySearch(num int, slice []int) int {
	low := 0
	high := len(slice) - 1
	for low <= high {
		mid := (low + high) / 2
		if slice[mid] < num {
			low = mid + 1
		} else if slice[mid] > num {
			high = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

func main() {
	newSlice := CreateSlice()
	fmt.Println(newSlice)
	sortedSlice := quickSort(newSlice)
	fmt.Println(sortedSlice)
	fmt.Println(binarySearch(5, sortedSlice))
}
