package main

import "fmt"

func Group(values []float64, step float64) map[int][]float64 {
	groups := make(map[int][]float64)
	for _, value := range values {
		groupKey := int(value/step) * int(step)
		groups[groupKey] = append(groups[groupKey], value)
	}
	return groups
}

func main() {
	data := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	step := 10.0
	groups := Group(data, step)
	for key, values := range groups {
		fmt.Printf("%d:{%v}\n", key, values)
	}
}
