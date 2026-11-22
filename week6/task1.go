package main

import "fmt"

// apply принимает функцию как аргумент.
func apply(numbers []int, operation func(int) int) []int {
	result := make([]int, len(numbers))

	for i, number := range numbers {
		result[i] = operation(number)
	}

	return result
}

// counter возвращает функцию с собственным состоянием.
func counter() func() int {
	count := 0

	return func() int {
		count++
		return count
	}
}

// makeAccumulator возвращает функцию,
// которая накапливает переданные ей числа.
func makeAccumulator() func(int) int {
	total := 0

	return func(number int) int {
		total += number
		return total
	}
}

// memoize возвращает функцию,
// запоминающую результаты предыдущих вычислений.
func memoize(f func(int) int) func(int) int {
	cache := make(map[int]int)

	return func(number int) int {
		// Проверяем, вычисляли ли мы результат для этого аргумента.
		if value, ok := cache[number]; ok {
			return value
		}

		// Если результата нет, вызываем исходную функцию.
		value := f(number)

		// Сохраняем результат.
		cache[number] = value

		return value
	}
}

func main() {
	fmt.Println("1. Функция как значение")

	square := func(number int) int {
		return number * number
	}

	fmt.Println(square(5)) // 25

	result := apply([]int{1, 2, 3}, square)
	fmt.Println(result) // [1 4 9]

	fmt.Println("\n2. Независимые счетчики")

	firstCounter := counter()
	secondCounter := counter()

	fmt.Println(firstCounter())  // 1
	fmt.Println(firstCounter())  // 2
	fmt.Println(secondCounter()) // 1
	fmt.Println(firstCounter())  // 3

	fmt.Println("\n3. Захват переменной")

	x := 10

	show := func() {
		fmt.Println(x)
	}

	x = 20

	show() // 20

	fmt.Println("\n4. Аккумулятор")

	acc := makeAccumulator()

	fmt.Println(acc(10)) // 10
	fmt.Println(acc(5))  // 15
	fmt.Println(acc(-3)) // 12

	fmt.Println("\n5. Мемоизация")

	calls := 0

	slowFunction := func(number int) int {
		calls++

		fmt.Println("Исходная функция вычисляет результат для:", number)

		return number * number
	}

	cachedFunction := memoize(slowFunction)

	fmt.Println(cachedFunction(4)) // Вычисление: 16
	fmt.Println(cachedFunction(4)) // Из кеша: 16
	fmt.Println(cachedFunction(5)) // Вычисление: 25
	fmt.Println(cachedFunction(4)) // Из кеша: 16
	fmt.Println(cachedFunction(5)) // Из кеша: 25

	fmt.Println("Количество реальных вычислений:", calls) // 2
}