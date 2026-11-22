package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func busy(duration time.Duration) {
    deadline := time.Now().Add(duration)

    for time.Now().Before(deadline) {
        // Deliberately busy
    }
}

func main() {
    fmt.Println("CPUs:", runtime.NumCPU())
    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

    var wg sync.WaitGroup

    for i := 0; i < 8; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()
            busy(2 * time.Second)
        }()
    }

    wg.Wait()
}