package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    const N = 100_000

    stop := make(chan struct{})
    var wg sync.WaitGroup

    var before, after runtime.MemStats
    runtime.ReadMemStats(&before)

    start := time.Now()

    for i := 0; i < N; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()
            <-stop
        }()
    }

    duration := time.Since(start)

    runtime.ReadMemStats(&after)

    delta := int64(after.StackSys) - int64(before.StackSys)

    fmt.Println("Goroutines:", runtime.NumGoroutine())
    fmt.Println("Creation time:", duration)
    fmt.Printf("StackSys delta: %.2f MB\n",
        float64(delta)/(1024*1024))
    fmt.Printf("Approx KB per goroutine: %.2f\n",
        float64(delta)/1024/float64(N))

    close(stop)
    wg.Wait()

    fmt.Println("After close:", runtime.NumGoroutine())
}