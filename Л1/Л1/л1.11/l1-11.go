package main

import "fmt"

func Intersection(slice1, slice2 []int) []int {
	set := make(map[int]bool)
	var result []int

	for _, item := range slice1 {
		set[item] = true //ключ со значением true
	}

	for _, item := range slice2 {
		if set[item] {
			result = append(result, item)
			set[item] = false //защита от дубликатов
		}
	}
	return result
}

func main() {
	slice1 := []int{1, 3, 3, 5, 6, 7}
	slice2 := []int{2, 4, 5, 3, 6, 7}
	slice3 := Intersection(slice1, slice2)
	fmt.Println(slice3)
}
