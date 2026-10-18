package main

import (
	"fmt"
	"os"

	"sis1/internal/texttask"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	service, err := texttask.NewService(texttask.TextAnalyzer{})
	if err != nil {
		return err
	}

	fmt.Println("Text Analyzer - SIS Checkpoint 1")
	inputs := []string{
		"Hello Go world",
		"Привет, Go!\nИзучаем Go.",
	}

	tasks := make([]texttask.Task, 0, len(inputs))
	for _, text := range inputs {
		task, err := service.Create(text)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
		fmt.Printf("Created task #%d: state=%s, text=%q\n", task.ID, task.State, task.Text)
	}

	for _, task := range tasks {
		completed, err := service.Process(task.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Processed task #%d: state=%s, words=%d, characters=%d, lines=%d\n",
			completed.ID, completed.State, completed.Result.WordCount,
			completed.Result.CharacterCount, completed.Result.LineCount)
	}

	stored, err := service.Get(tasks[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("Retrieved task #%d: state=%s, result=%+v\n", stored.ID, stored.State, stored.Result)

	if _, err := service.Create(" \t\n "); err != nil {
		fmt.Println("Expected validation error:", err)
	}
	if _, err := service.Get(999); err != nil {
		fmt.Println("Expected lookup error:", err)
	}
	if _, err := service.Process(tasks[0].ID); err != nil {
		fmt.Println("Expected state error:", err)
	}

	return nil
}
