package main

import (
	"fmt"
	"math"
)

// СТРУКТУРА
type Point struct {
	x float64
	y float64
}

// КОНСТРУКТОР
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

// МЕТОД
func (p Point) Distance(other Point) float64 {
	dx := p.x - other.x
	dy := p.y - other.y
	return math.Sqrt(dx*dx + dy*dy)
}

func main() {
	point1 := NewPoint(2.4, 6.7)
	point2 := NewPoint(4.6, 1.5)
	distance := point1.Distance(point2)
	fmt.Printf("Расстояние между точками: %f", distance)
}
