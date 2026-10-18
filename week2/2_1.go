package main

import "fmt"

func main() {
    // 1. Slice literal
    slice1 := []int{10, 20, 30}
    fmt.Println("slice1:", slice1)

    // 2. Nil slice using var
    var slice2 []int
    fmt.Println("slice2:", slice2)
	slice2 = append(slice2, 1, 2, 3)
	fmt.Println("slice2 after append:", slice2)
	slice2 = append(slice2, slice1...)
	fmt.Println("slice2 after appending slice1:", slice2)

    // 3. Slice using make
    slice3 := make([]int, 3, 5)
    fmt.Println("slice3:", slice3)
    fmt.Println("len:", len(slice3))
    fmt.Println("cap:", cap(slice3))

    // 4. Slice from an array
    array := [5]int{10, 20, 30, 40, 50}
    slice4 := array[0:4]

    fmt.Println("slice4:", slice4, len(slice4), cap(slice4))

	var slice5 []int = []int{1, 2, 3, 4, 5}
	fmt.Println("slice5:", slice5)

	array2 := [5]int{10, 20, 30, 40, 50}
    slice9 := array2[0:4]

    // Меняем первый элемент через слайс
    slice9[0] = 999

    fmt.Println("Массив:", array2)  // Выведет: [999 20 30 40 50] — исходный массив изменился!
    fmt.Println("Слайс: ", slice9) // Выведет: [999 20 30 40]
	slice9 = append(slice9, 1000)
	fmt.Println("Массив после append:", array2)
	fmt.Println("Слайс после append: ", slice9)
}