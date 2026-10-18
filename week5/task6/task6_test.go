package main

import "testing"

type Number interface {
	~int | ~int64 | ~float64
}

func SumInts(nums []int) int {
	total := 0

	for _, n := range nums {
		total += n
	}

	return total
}

func Sum[T Number](nums []T) T {
	var total T

	for _, n := range nums {
		total += n
	}

	return total
}

func SumAny(nums []any) int {
	total := 0

	for _, n := range nums {
		total += n.(int)
	}

	return total
}

func ToAny(nums []int) []any {
	result := make([]any, len(nums))

	for i, n := range nums {
		result[i] = n
	}

	return result
}

var result int

func makeNumbers() []int {
	nums := make([]int, 1000)

	for i := range nums {
		nums[i] = i
	}

	return nums
}

func BenchmarkSumInts(b *testing.B) {
	nums := makeNumbers()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result = SumInts(nums)
	}
}

func BenchmarkSumGeneric(b *testing.B) {
	nums := makeNumbers()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result = Sum(nums)
	}
}

func BenchmarkSumAny(b *testing.B) {
	nums := makeNumbers()

	// Конвертируем только ОДИН раз,
	// до начала измерения.
	anyNums := ToAny(nums)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result = SumAny(anyNums)
	}
}

func BenchmarkSumAnyWithBoxing(b *testing.B) {
	nums := makeNumbers()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result = SumAny(ToAny(nums))
	}
}