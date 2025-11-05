package main

import (
	"fmt"
	"os"

	"github.com/tvbondar/go-proj/Л2/л2.9/pkg/unpacker"
)

func main() {
	var packedString string
	fmt.Println("Введите строку для распаковки: \n")
	_, err := fmt.Scanf("%s", &packedString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении ввода: ", err)
		os.Exit(1)
	}

	unpackedString, err := unpacker.Unpack(packedString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при распаковке строки: ", err)
		os.Exit(1)
	}
	fmt.Println("Распакованная строка: ", unpackedString)
}
