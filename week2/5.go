package main

import "fmt"

func main() {
    array := [5]int{10, 20, 30, 40, 50}
    slice4 := array[0:4] // len: 4, cap: 5

    // 1. Сначала заполняем слайс до предела емкости базового массива
    slice4 = append(slice4, 99) // Теперь len: 5, cap: 5. Свободных мест в array больше нет.
    fmt.Println("После 1-го append:")
    fmt.Println("Массив:", array) // [10 20 30 40 99] (массив изменился!)
    fmt.Println("Слайс: ", slice4) // [10 20 30 40 99]

    // 2. Добавляем элемент ЗА пределы емкости массива
    slice4 = append(slice4, 777) 
    
    fmt.Println("\nПосле 2-го append (выход за пределы):")
    fmt.Println("Массив:", array) // [10 20 30 40 99] — ОСТАЛСЯ ПРЕЖНИМ!
    fmt.Println("Слайс: ", slice4) // [10 20 30 40 99 777] — УШЕЛ В СВОЙ МАССИВ!
    fmt.Println("Новая емкость слайса:", cap(slice4)) // Выведет 10 (емкость удвоилась)
}