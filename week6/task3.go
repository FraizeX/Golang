package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("student not found")
var ErrTimeout = errors.New("timeout")

type GradeError struct {
	Student string
	Grade   int
}

func (e *GradeError) Error() string {
	return fmt.Sprintf(
		"invalid grade %d for %s",
		e.Grade,
		e.Student,
	)
}

type GradeBook struct {
	grades map[string]int
}

func (gb *GradeBook) Get(name string) (int, error) {
	grade, ok := gb.grades[name]

	if !ok {
		return 0, ErrNotFound
	}

	return grade, nil
}

func (gb *GradeBook) Set(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{
			Student: name,
			Grade:   grade,
		}
	}

	if gb.grades == nil {
		gb.grades = make(map[string]int)
	}

	gb.grades[name] = grade

	return nil
}

func AddBonus(gb *GradeBook, name string, bonus int) error {
	grade, err := gb.Get(name)

	if err != nil {
		return fmt.Errorf(
			"add bonus for %q: get: %w",
			name,
			err,
		)
	}

	err = gb.Set(name, grade+bonus)

	if err != nil {
		return fmt.Errorf(
			"add bonus for %q: set: %w",
			name,
			err,
		)
	}

	return nil
}

// Здесь специально используется %v вместо %w
// в ветке ошибки Get.
func AddBonusBroken(gb *GradeBook, name string, bonus int) error {
	grade, err := gb.Get(name)

	if err != nil {
		return fmt.Errorf(
			"add bonus for %q: get: %v",
			name,
			err,
		)
	}

	err = gb.Set(name, grade+bonus)

	if err != nil {
		return fmt.Errorf(
			"add bonus for %q: set: %w",
			name,
			err,
		)
	}

	return nil
}

func report(err error) {
	var gradeError *GradeError

	switch {
	case err == nil:
		fmt.Println("OK")

	case errors.Is(err, ErrNotFound):
		fmt.Println("404:", err)

	case errors.As(err, &gradeError):
		fmt.Printf(
			"400: %v; student=%s; grade=%d\n",
			err,
			gradeError.Student,
			gradeError.Grade,
		)

	default:
		fmt.Println("500:", err)
	}
}

// Пример собственного типа, который оборачивает другую ошибку.
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query %q: %v", e.Query, e.Err)
}

func (e *QueryError) Unwrap() error {
	return e.Err
}

func main() {
	var book GradeBook

	fmt.Println("1. Добавление корректной оценки")
	report(book.Set("Aruzhan", 90))

	fmt.Println("\n2. Корректный бонус")
	report(AddBonus(&book, "Aruzhan", 5))

	grade, err := book.Get("Aruzhan")

	if err != nil {
		report(err)
	} else {
		fmt.Println("Текущая оценка:", grade) // 95
	}

	fmt.Println("\n3. Студент отсутствует")
	report(AddBonus(&book, "Timur", 5)) // 404

	fmt.Println("\n4. Оценка выходит за диапазон")
	report(AddBonus(&book, "Aruzhan", 10)) // 400: 105

	fmt.Println("\n5. Потеря причины через текстовое форматирование")
	report(AddBonusBroken(&book, "Timur", 5)) // 500

	fmt.Println("\n6. Неизвестная ошибка")
	report(errors.New("database is unavailable")) // 500

	fmt.Println("\n7. Сравнение ошибки")

	wrapped := fmt.Errorf("load student: %w", ErrNotFound)

	fmt.Println(wrapped == ErrNotFound)          // false
	fmt.Println(errors.Is(wrapped, ErrNotFound)) // true

	fmt.Println("\n8. Собственный Unwrap")

	queryError := &QueryError{
		Query: "SELECT * FROM students",
		Err:   ErrTimeout,
	}

	outerError := fmt.Errorf("load students: %w", queryError)

	fmt.Println(outerError)
	fmt.Println(errors.Is(outerError, ErrTimeout)) // true

	fmt.Println("\n9. Typed nil")

	var pointer *GradeError
	var interfaceError error = pointer

	fmt.Println(pointer == nil)        // true
	fmt.Println(interfaceError == nil) // false
}