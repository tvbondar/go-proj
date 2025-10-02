package main

import (
	"fmt"
	"strings"
)

func reverseString(s string) string {
	phrase := strings.Fields(s)
	for i := 0; i < len(phrase)/2; i++ {
		j := len(phrase) - 1 - i
		phrase[i], phrase[j] = phrase[j], phrase[i]
	}
	return strings.Join(phrase, " ")
}

func main() {
	String := "snow dog sun"
	fmt.Println(String)
	fmt.Println(reverseString(String))

	String2 := "new long short string bool"
	fmt.Println(String2)
	fmt.Println(reverseString(String2))

}
