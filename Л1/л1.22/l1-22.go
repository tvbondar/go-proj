package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewInt(2345678995869689880)
	b := big.NewInt(13948586493902198)
	fmt.Printf("Число a: %d\nЧисло b: %d\n", a, b)
	result := new(big.Int)

	result.Add(a, b)
	fmt.Printf("Сложение: %v\n", result)

	result.Sub(a, b)
	fmt.Printf("Вычитание: %v\n", result)

	result.Mul(a, b)
	fmt.Printf("Умножение: %v\n", result)

	result.Div(a, b)
	fmt.Printf("Деление: %v\n", result)
}
