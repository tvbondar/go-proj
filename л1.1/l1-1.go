package main

import "fmt"

type Human struct {
	Run   string
	Jump  string
	Swim  string
	Learn string
}

func (h Human) printRun() {
	fmt.Println(h.Run)
}

func (h Human) printJump() {
	fmt.Println(h.Jump)
}

func (h Human) printSwim() {
	fmt.Println(h.Swim)
}

func (h Human) printLearn() {
	fmt.Println(h.Learn)
}

type Action struct {
	Human
}

func main() {
	action := Action{
		Human: Human{
			Run:   "Бежит",
			Jump:  "Прыгает",
			Swim:  "Плывет",
			Learn: "Учится",
		},
	}
	action.printRun()
	action.printJump()
	action.printSwim()
	action.printLearn()
}
