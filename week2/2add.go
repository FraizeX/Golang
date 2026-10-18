package main

import "fmt"

func modifySlice(numbers []int) {
    numbers[0] = 100
    numbers[1] = 200
}


func addNumber1(numbers []int) []int {
    numbers = append(numbers, 100)

    return numbers
}

func addNumber(numbers []int) {
    numbers = append(numbers, 100)
}

func main() {
    numbers := []int{1, 2, 3, 4}

    fmt.Println("Before:", numbers)

    modifySlice(numbers)

    fmt.Println("After:", numbers)

	addNumber(numbers)

    fmt.Println(numbers)

	numbers = addNumber1(numbers)
	fmt.Println("After addNumber1:", numbers)
}