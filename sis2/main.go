package main

import (
	"fmt"
	"os"

	"sis2/internal/texttask"
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

	defer service.Shutdown()
	fmt.Printf("Text Analyzer - %d workers, %d queue slots\n", texttask.WorkerCount, texttask.QueueCapacity)
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
		if err := service.Submit(task.ID); err != nil {
			return err
		}
		fmt.Printf("Accepted task #%d\n", task.ID)
	}

	if err := service.Submit(tasks[0].ID); err != nil {
		fmt.Println("Expected duplicate submission error:", err)
	}
	service.Shutdown()
	for _, task := range tasks {
		stored, err := service.Get(task.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Task #%d: state=%s, words=%d, characters=%d, lines=%d\n",
			stored.ID, stored.State, stored.Result.WordCount,
			stored.Result.CharacterCount, stored.Result.LineCount)
	}

	if _, err := service.Create(" \t\n "); err != nil {
		fmt.Println("Expected validation error:", err)
	}
	if _, err := service.Get(999); err != nil {
		fmt.Println("Expected lookup error:", err)
	}
	if err := service.Submit(tasks[0].ID); err != nil {
		fmt.Println("Expected submission error after shutdown:", err)
	}

	return nil
}
