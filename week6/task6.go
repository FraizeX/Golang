package main

import (
	"errors"
	"fmt"
)

var ErrTaskFailed = errors.New("service unavailable")

type Task struct {
	Name string
	Run  func() error
}

func runOne(task Task) (err error) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			err = fmt.Errorf(
				"task %s: panic: %v",
				task.Name,
				panicValue,
			)
		}
	}()

	taskError := task.Run()

	if taskError != nil {
		return fmt.Errorf(
			"task %s: %w",
			task.Name,
			taskError,
		)
	}

	return nil
}

func SafeRun(tasks []Task) []error {
	var taskErrors []error

	for _, task := range tasks {
		err := runOne(task)

		if err != nil {
			taskErrors = append(taskErrors, err)
		}
	}

	return taskErrors
}

func main() {
	tasks := []Task{
		{
			Name: "success",
			Run: func() error {
				fmt.Println("Первая успешная задача выполнена")
				return nil
			},
		},
		{
			Name: "index out of range",
			Run: func() error {
				numbers := []int{1, 2, 3}

				fmt.Println(numbers[10])

				return nil
			},
		},
		{
			Name: "nil map",
			Run: func() error {
				var values map[string]int

				values["answer"] = 42

				return nil
			},
		},
		{
			Name: "missing config",
			Run: func() error {
				panic("config is missing")
			},
		},
		{
			Name: "returned error",
			Run: func() error {
				return ErrTaskFailed
			},
		},
		{
			Name: "last success",
			Run: func() error {
				fmt.Println("Последняя успешная задача тоже выполнена")
				return nil
			},
		},
	}

	taskErrors := SafeRun(tasks)

	fmt.Println("\nПолученные ошибки:")

	for _, err := range taskErrors {
		fmt.Println(err)
	}

	joined := errors.Join(taskErrors...)

	fmt.Println("\nКоличество ошибок:", len(taskErrors))

	fmt.Println(
		"Sentinel сохранился:",
		errors.Is(joined, ErrTaskFailed),
	)
}