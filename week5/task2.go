package main

import (
	"errors"
	"fmt"
)

type Celsius float64

func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}

func describe(v any) string {
	switch x := v.(type) {

	case nil:
		return "nil"

	case int, int64:
		return fmt.Sprintf("integer %v", x)

	case string:
		return fmt.Sprintf("string of %d bytes", len(x))

	case []int:
		return fmt.Sprintf("slice with %d ints", len(x))

	case bool:
		return fmt.Sprintf("bool %v", x)

	case map[string]int:
		return fmt.Sprintf("map with %d entries", len(x))

	case error:
		return "error: " + x.Error()

	case fmt.Stringer:
		return "Stringer: " + x.String()

	default:
		return fmt.Sprintf("unknown type %T", x)
	}
}

func main() {
	values := []any{
		nil,
		7,
		int64(7),
		"привет",
		[]int{1, 2, 3},
		errors.New("disk full"),
		Celsius(21.5),
		3.14,
		struct{}{},
		true,
		map[string]int{"Go": 100},
	}

	for _, v := range values {
		fmt.Println(describe(v))
	}
}