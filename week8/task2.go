package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func atomicExample() {
    var counter atomic.Int64
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()
            counter.Add(1)
        }()
    }

    wg.Wait()
    fmt.Println("Atomic:", counter.Load())
}

func sliceExample() {
    results := make([]int, 1000)
    var wg sync.WaitGroup

    for i := range results {
        wg.Add(1)

        go func(index int) {
            defer wg.Done()
            results[index] = 1
        }(i)
    }

    wg.Wait()

    sum := 0
    for _, v := range results {
        sum += v
    }

    fmt.Println("Slice:", sum)
}

func recoverExample() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer wg.Done()

        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered:", r)
            }
        }()

        panic("Something went wrong")
    }()

    wg.Wait()
    fmt.Println("Program is still alive")
}

func main() {
    atomicExample()
    sliceExample()
    recoverExample()
}