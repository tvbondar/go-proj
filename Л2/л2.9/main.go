package main

import (
	"fmt"

	"github.com/tvbondar/go-proj/Л2/л2.9/pkg/unpacker"
)

func main() {
	var packedString string
	fmt.Println("Введите строку для распаковки: \n")
	_, err := fmt.Scanf("%s", &packedString)
	if err != nil {
		fmt.Println(err)
	} else {
		unpackedString := unpacker.Unpack(packedString)
		fmt.Println("Распакованная строка: ", unpackedString)
	}
}
