package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(k float64) {
	r.Width *= k
	r.Height *= k
}

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	r := Rectangle{
		Width:  3,
		Height: 4,
	}

	c := Circle{
		Radius: 2,
	}

	r.Scale(2)

	shapes := []Shape{
		r,
		&c,
	}

	for _, shape := range shapes {
		fmt.Printf(
			"Area: %.2f, Perimeter: %.2f\n",
			shape.Area(),
			shape.Perimeter(),
		)
	}

	// shapes = append(shapes, Circle{Radius: 1}) not work, because circle not implement shape but *circle impelement shape

	shapes = append(shapes, &Circle{Radius: 1}) 

	for _, shape := range shapes {
		fmt.Printf(
			"Area: %.2f, Perimeter: %.2f\n",
			shape.Area(),
			shape.Perimeter(),
		)
	}
}