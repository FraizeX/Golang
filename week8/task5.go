package main

import (
    "context"
    "errors"
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

func slowQuery() int {
    time.Sleep(50 * time.Millisecond)
    return 42
}

func fetchLeaky(timeout time.Duration) (int, error) {
    ch := make(chan int)

    go func() {
        ch <- slowQuery()
    }()

    select {
    case v := <-ch:
        return v, nil
    case <-time.After(timeout):
        return 0, errors.New("timeout")
    }
}

func fetchFixed(timeout time.Duration) (int, error) {
    ch := make(chan int, 1)

    go func() {
        ch <- slowQuery()
    }()

    select {
    case v := <-ch:
        return v, nil
    case <-time.After(timeout):
        return 0, errors.New("timeout")
    }
}

func startTickerLeaky() {
    go func() {
        for {
            time.Sleep(10 * time.Millisecond)
        }
    }()
}

func startTickerFixed(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(10 * time.Millisecond)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                // Periodic work
            }
        }
    }()
}

func main() {
    fmt.Println("Start:", runtime.NumGoroutine())

    for i := 0; i < 100; i++ {
        fetchLeaky(time.Millisecond)
    }

    time.Sleep(100 * time.Millisecond)
    fmt.Println("After leaks:", runtime.NumGoroutine())

    pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
    pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)

    for i := 0; i < 100; i++ {
        fetchFixed(time.Millisecond)
    }

    time.Sleep(100 * time.Millisecond)
    fmt.Println("After fixed fetch:", runtime.NumGoroutine())

    for i := 0; i < 10; i++ {
        startTickerLeaky()
    }

    fmt.Println("After leaky tickers:", runtime.NumGoroutine())

    ctx, cancel := context.WithCancel(context.Background())

    for i := 0; i < 10; i++ {
        startTickerFixed(ctx)
    }

    fmt.Println("After fixed tickers:", runtime.NumGoroutine())

    cancel()
    time.Sleep(50 * time.Millisecond)

    fmt.Println("After cancel:", runtime.NumGoroutine())
}