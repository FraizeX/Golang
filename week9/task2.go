package main

import (
    "fmt"
    "time"
)

func run(size int) {
    ch := make(chan int, size)
    done := make(chan struct{})

    start := time.Now()

    go func() {
        defer close(done)

        for i := 1; i <= 3; i++ {
            ch <- i

            fmt.Printf(
                "buffer=%d sent=%d time=%v\n",
                size, i,
                time.Since(start).Round(100*time.Millisecond),
            )
        }
    }()

    time.Sleep(300 * time.Millisecond)

    for i := 0; i < 3; i++ {
        <-ch
        time.Sleep(100 * time.Millisecond)
    }

    <-done
}

func main() {
    for _, size := range []int{0, 1, 2, 3, 10} {
        run(size)
        fmt.Println("----------")
    }
}