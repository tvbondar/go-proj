package main

import (
	"fmt"
	"strings"
)

func reverseString(s string) string {
	phrase := strings.Split(s, " ")
	reversedPhrase := make([]string, len(phrase))
	for i := range phrase {
		reversedPhrase[len(phrase)-1-i] = phrase[i]
	}
	return strings.Join(reversedPhrase, " ")

}

func main() {
	var String string
	_, err := fmt.Scanln(&String)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	fmt.Println(reverseString(String))

}
