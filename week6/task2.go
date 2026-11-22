package main

import "fmt"

func fibonacci() func() int {
	a, b := 0, 1

	return func() int {
		current := a

		a, b = b, a+b

		return current
	}
}

func main() {
	next := fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Print(next(), " ")
	}

	fmt.Println()

	// Новый вызов fibonacci создаёт независимую последовательность.
	another := fibonacci()

	fmt.Println(another()) // 0
	fmt.Println(another()) // 1

	// Первый генератор продолжает свою последовательность.
	fmt.Println(next()) // 55
}