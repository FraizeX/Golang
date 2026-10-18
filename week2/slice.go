package main
import "fmt"
func main() {
	array := [5]int{10, 20, 30, 40, 50}

	s := array[1:3]

	fmt.Println(s, len(s), cap(s))

	c := array[0:4]
	fmt.Println(c, len(c), cap(c))

	v := s[0:4]
	fmt.Println(v, len(v), cap(v))
}