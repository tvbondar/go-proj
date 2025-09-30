package main

import "fmt"

func main() {
	a := 4
	b := 5
	c := 10
	d := 12

	//Способ 1
	fmt.Println(a, b)
	a = a + b
	b = a - b
	a = a - b
	fmt.Println(a, b)

	//Способ 2
	fmt.Println(c, d)
	c = c ^ d
	d = c ^ d
	c = c ^ d
	fmt.Println(c, d)

}
