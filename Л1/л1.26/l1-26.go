package main

import (
	"fmt"
	"strings"
)

func IsUnique(phrase string) bool {
	phrase = strings.ToLower(phrase)
	var seen = make(map[rune]bool)
	for _, char := range phrase {
		if seen[char] {
			return false
		}
		seen[char] = true
	}
	return true
}

func main() {
	var Phrase string
	fmt.Println("Введите строку: ")
	fmt.Scanln(&Phrase)
	fmt.Println("false - есть повторы, true - нет повторов.")
	fmt.Println(IsUnique(Phrase))
}
