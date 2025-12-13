package main

import (
	"fmt"

	"l2.11/internal/anagram"
)

func main() {
	words := []string{"пятак", "ПятКа", "тяпка", "листок", "слиток", "столик", "стол", "пятак"}
	res := anagram.FindAnagrams(words)
	fmt.Println(res)

	otherWords := []string{
		"кот", "ток", "окт",
		"арбуз",
		"мир", "рим",
		"дом", "мод",
		"соло"}
	otherRes := anagram.FindAnagrams(otherWords)
	fmt.Println(otherRes)
}
