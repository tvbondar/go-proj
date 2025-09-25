package main

import "fmt"

func RmDuplicates(slice []string) []string {
	checked := make(map[string]bool)
	result := []string{}

	for _, word := range slice {
		if !checked[word] {
			checked[word] = true
			result = append(result, word)
		}
	}
	return result
}

func main() {
	slice := []string{"cat", "cat", "dog", "cat", "tree"}
	new := RmDuplicates(slice)
	fmt.Println(new)
}
