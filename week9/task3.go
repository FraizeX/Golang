package main

import (
    "fmt"
    "sync"
)

func collect(workers int) <-chan int {
    results := make(chan int)
    var wg sync.WaitGroup

    for w := 1; w <= workers; w++ {
        wg.Add(1)

        go func(worker int) {
            defer wg.Done()

            for i := 0; i < 3; i++ {
                results <- worker*10 + i
            }
        }(w)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}

func main() {
    total := 0

    for value := range collect(3) {
        total += value
    }

    fmt.Println("sum:", total)
}