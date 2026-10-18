package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "invalid field: " + e.Field
}

func validate(age int) error {
	if age < 0 {
		return &ValidationError{
			Field: "age",
		}
	}

	// ВАЖНО:
	return nil
}

func main() {
	err := validate(20)

	fmt.Println("valid:", err == nil)

	err = validate(-1)

	if err != nil {
		fmt.Println("error:", err)

		var validationErr *ValidationError

		if errors.As(err, &validationErr) {
			fmt.Println("Field:", validationErr.Field)
		}
	}
}