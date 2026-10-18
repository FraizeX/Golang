package main

import (
	"fmt"
	"strings"
)

type List[T any] struct {
	next *List[T]
	val  T
}

// Возвращаем новый head.
// Поэтому работает даже если l == nil.
func (l *List[T]) Push(v T) *List[T] {
	return &List[T]{
		val:  v,
		next: l,
	}
}

func (l *List[T]) Len() int {
	count := 0

	for current := l; current != nil; current = current.next {
		count++
	}

	return count
}

func (l *List[T]) String() string {
	var parts []string

	for current := l; current != nil; current = current.next {
		parts = append(parts, fmt.Sprint(current.val))
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func main() {
	var list *List[int]

	// list == nil, но Push всё равно работает.
	list = list.Push(10)
	list = list.Push(20)
	list = list.Push(30)

	fmt.Println(list)
	fmt.Println("Len:", list.Len())

	var words *List[string]

	words = words.Push("Go")
	words = words.Push("Generics")

	fmt.Println(words)
	fmt.Println("Len:", words.Len())
}