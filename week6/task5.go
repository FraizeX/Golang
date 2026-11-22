package main

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

func demonstrateLIFO() {
	defer fmt.Println("1: deferred first")
	defer fmt.Println("2: deferred second")
	defer fmt.Println("3: deferred third")

	fmt.Println("main body")
}

func demonstrateArguments() {
	x := 1

	defer fmt.Println("deferred x =", x)

	defer func() {
		fmt.Println("closure x =", x)
	}()

	x = 2

	fmt.Println("current x =", x)
}

func double() (result int) {
	defer func() {
		result *= 2
	}()

	return 21
}

func notDouble() int {
	result := 21

	defer func() {
		result *= 2
	}()

	return result
}

func trace(name string) func() {
	fmt.Println("enter", name)

	return func() {
		fmt.Println("exit", name)
	}
}

func work() {
	defer trace("work")()

	fmt.Println("working...")
}

type Resource struct {
	Name string
}

func Open(name string) (*Resource, error) {
	if name == "" {
		return nil, errors.New("empty resource name")
	}

	fmt.Println("open", name)

	return &Resource{
		Name: name,
	}, nil
}

func (resource *Resource) Close() {
	fmt.Println("close", resource.Name)
}

func processAllBad(names []string) error {
	for _, name := range names {
		resource, err := Open(name)

		if err != nil {
			return err
		}

		defer resource.Close()

		fmt.Println("process", name)
	}

	return nil
}

func processOne(name string) error {
	resource, err := Open(name)

	if err != nil {
		return err
	}

	defer resource.Close()

	fmt.Println("process", name)

	return nil
}

func processAllGood(names []string) error {
	for _, name := range names {
		if err := processOne(name); err != nil {
			return err
		}
	}

	return nil
}

func timer(name string) func() {
	start := time.Now()

	return func() {
		fmt.Printf("%s took %s\n", name, time.Since(start))
	}
}

func sortNumbers(numbers []int) {
	defer timer("sort")()

	sort.Ints(numbers)
}

func main() {
	fmt.Println("1. Порядок выполнения defer")
	demonstrateLIFO()

	fmt.Println("\n2. Вычисление аргументов")
	demonstrateArguments()

	fmt.Println("\n3. Именованные результаты")
	fmt.Println(double())    // 42
	fmt.Println(notDouble()) // 21

	fmt.Println("\n4. Двойной вызов")
	work()

	names := []string{"A", "B", "C"}

	fmt.Println("\n5. defer внутри внешнего цикла")

	if err := processAllBad(names); err != nil {
		fmt.Println("Ошибка:", err)
	}

	fmt.Println("\n6. Обработка через отдельную функцию")

	if err := processAllGood(names); err != nil {
		fmt.Println("Ошибка:", err)
	}

	fmt.Println("\n7. Измерение времени")

	numbers := make([]int, 100000)

	for i := range numbers {
		numbers[i] = (len(numbers) - i) * 37 % 100003
	}

	sortNumbers(numbers)

	fmt.Println("Первое число:", numbers[0])
	fmt.Println("Последнее число:", numbers[len(numbers)-1])
}