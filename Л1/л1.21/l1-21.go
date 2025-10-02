package main

import "fmt"

//определяет целевую сигнатуру метода, который мы адаптируем из сторонних классов
type Target interface {
	Operation()
}

//адаптируемый класс, его методы надо вызвать в другом месте
type Adaptee struct {
}

// метод адаптируемого класса
func (adaptee *Adaptee) AdaptedOperation() {
	fmt.Println("Adapted Operation")
}

//класс конкретного адаптера
type OneAdapter struct {
	*Adaptee
}

//реализация метода интерфейса, который вызывает адаптируемый класс
func (adapter *OneAdapter) Operation() {
	adapter.AdaptedOperation()
}

//конструктор нового адаптера
func NewAdapter(adaptee *Adaptee) Target {
	return &OneAdapter{adaptee}
}

//применение
func main() {
	adapter := NewAdapter(&Adaptee{})
	adapter.Operation()
}
