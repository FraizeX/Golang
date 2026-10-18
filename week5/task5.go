package main

import (
	"cmp"
	"fmt"
)

type Number interface {
	~int | ~int64 | ~float64
}

type Score int

func Sum[T Number](nums []T) T {
	var total T

	for _, n := range nums {
		total += n
	}

	return total
}

func Average[T Number](nums []T) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}

	total := Sum(nums)

	return float64(total) / float64(len(nums)), true
}

func MaxOf[T cmp.Ordered](items []T) (T, bool) {
	var zero T

	if len(items) == 0 {
		return zero, false
	}

	maxValue := items[0]

	for _, v := range items[1:] {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue, true
}

func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}

	return -1
}

func Filter[T any](items []T, keep func(T) bool) []T {
	result := make([]T, 0)

	for _, v := range items {
		if keep(v) {
			result = append(result, v)
		}
	}

	return result
}

func main() {
	fmt.Println("Sum:")
	fmt.Println(Sum([]int{1, 2, 3}))
	fmt.Println(Sum([]float64{1.5, 2.5}))

	avg, ok := Average([]int{10, 20, 30})
	fmt.Println("Average:", avg, ok)

	maxInt, ok := MaxOf([]int{4, 9, 2})
	fmt.Println("Max int:", maxInt, ok)

	maxString, ok := MaxOf([]string{"apple", "banana", "cat"})
	fmt.Println("Max string:", maxString, ok)

	_, ok = MaxOf([]int{})
	fmt.Println("Empty max:", ok)

	fmt.Println(
		"Index:",
		IndexOf([]string{"Go", "Java", "Rust"}, "Java"),
	)

	scores := []Score{50, 75, 90, 60, 100}

	passed := Filter(scores, func(s Score) bool {
		return s >= 75
	})

	fmt.Println("Passed:", passed)

	fmt.Println("Score sum:", Sum(scores))
}