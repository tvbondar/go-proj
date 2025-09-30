package main

import "fmt"

func main() {
	var String string
	_, err := fmt.Scanln(&String)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	runeSlice := []rune(String)
	fmt.Printf("%c\n", runeSlice)
	for i := len(runeSlice) - 1; i >= 0; i-- {
		fmt.Printf("%c ", runeSlice[i])
	}
	fmt.Print("\n")

}
