package main

import "fmt"

func main() {
	var functions []func()

	for i := 0; i < 3; i++ {
		functions = append(functions, func() {
			fmt.Print(i, " ")
		})
	}

	for _, function := range functions {
		function()
	}

	fmt.Println()
}