package main

import "fmt"

func main() {
    array := [5]int{10, 20, 30, 40, 50}

    slice := array[1:4]

    fmt.Println("Original array:", array)
    fmt.Println("Original slice:", slice)

    slice[0] = 999

    fmt.Println("\nAfter changing slice:")
    fmt.Println("Array:", array)
    fmt.Println("Slice:", slice)

    array[2] = 777

    fmt.Println("\nAfter changing array:")
    fmt.Println("Array:", array)
    fmt.Println("Slice:", slice)

	array[1] = 555

	fmt.Println("\nAfter changing slice again:")
	fmt.Println("Array:", array)
	fmt.Println("Slice:", slice)
}