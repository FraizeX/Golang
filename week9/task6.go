package main

import (
    "errors"
    "fmt"
    "runtime"
    "strings"
    "sync"
    "time"
)

type Result struct {
    Word  string
    Upper string
}

var ErrQueueFull = errors.New("queue is full")

func producer(
    words []string,
    done <-chan struct{},
) <-chan string {

    jobs := make(chan string)

    go func() {
        defer close(jobs)

        for _, word := range words {
            select {
            case jobs <- word:
            case <-done:
                return
            }
        }
    }()

    return jobs
}

func startWorkers(
    n int,
    jobs <-chan string,
    done <-chan struct{},
) <-chan Result {

    results := make(chan Result)
    var wg sync.WaitGroup

    for i := 0; i < n; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()

            for {
                select {
                case <-done:
                    return

                case word, ok := <-jobs:
                    if !ok {
                        return
                    }

                    result := Result{
                        Word:  word,
                        Upper: strings.ToUpper(word),
                    }

                    select {
                    case results <- result:
                    case <-done:
                        return
                    }
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}

func Submit(
    queue chan<- string,
    job string,
    timeout time.Duration,
) error {

    timer := time.NewTimer(timeout)
    defer timer.Stop()

    select {
    case queue <- job:
        return nil

    case <-timer.C:
        return ErrQueueFull
    }
}

func main() {
    words := []string{
        "golang", "goroutine", "channel",
        "select", "worker", "pipeline",
        "scheduler", "context", "race",
    }

    fmt.Println("=== Normal execution ===")

    done1 := make(chan struct{})
    jobs1 := producer(words, done1)
    results1 := startWorkers(3, jobs1, done1)

    for result := range results1 {
        fmt.Println(result.Word, "->", result.Upper)
    }

    fmt.Println("=== Cancellation ===")

    done2 := make(chan struct{})
    jobs2 := producer(words, done2)
    results2 := startWorkers(3, jobs2, done2)

    for i := 0; i < 3; i++ {
        result, ok := <-results2
        if !ok {
            break
        }

        fmt.Println("Got:", result.Upper)
    }

    close(done2)

    // Drain remaining values until the closer finishes.
    for range results2 {
    }

    fmt.Println("Goroutines:", runtime.NumGoroutine())

    fmt.Println("=== Submit ===")

    queue := make(chan string, 1)

    err := Submit(queue, "first", time.Second)
    fmt.Println("First:", err)

    err = Submit(queue, "second", 50*time.Millisecond)
    fmt.Println("Second:", err)
}