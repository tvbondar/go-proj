package main

import (
	"fmt"
)

func setBit(num int64, i uint, bit int) int64 {
	if bit == 1 {
		// Установка бита в 1, используем OR
		return num | (1 << i)
	} else {
		// Установка бита в 0, используем AND NOT
		return num &^ (1 << i)
	}
}

func main() {
	var num int64
	var i uint
	var bit int

	fmt.Printf("Ввод числа: ")
	fmt.Scan(&num)

	fmt.Printf("Число %d в двоичной системе: %b\n", num, num)

	fmt.Print("Ввод номера бита: ")
	fmt.Scan(&i)

	fmt.Print("Установить в 0 или 1? ")
	fmt.Scan(&bit)

	result := setBit(num, i, bit)

	fmt.Printf("Результат операции: %d (%064b)\n", result, uint64(result))
}
