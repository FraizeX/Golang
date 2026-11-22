package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func isPrime(n int) bool {
    if n < 2 {
        return false
    }

    for d := 2; d*d <= n; d++ {
        if n%d == 0 {
            return false
        }
    }
    return true
}

func countPrimes(from, to int) int {
    count := 0

    for n := from; n < to; n++ {
        if isPrime(n) {
            count++
        }
    }
    return count
}

func cpuWork(workers int) (int, time.Duration) {
    const N = 6_000_000

    start := time.Now()
    results := make([]int, workers)
    var wg sync.WaitGroup

    for w := 0; w < workers; w++ {
        from := w * N / workers
        to := (w + 1) * N / workers

        wg.Add(1)

        go func(index, from, to int) {
            defer wg.Done()
            results[index] = countPrimes(from, to)
        }(w, from, to)
    }

    wg.Wait()

    total := 0
    for _, count := range results {
        total += count
    }

    return total, time.Since(start)
}

func ioWork(workers int) time.Duration {
    start := time.Now()
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()
            time.Sleep(100 * time.Millisecond)
        }()
    }

    wg.Wait()
    return time.Since(start)
}

func main() {
    cpus := runtime.NumCPU()

    fmt.Println("CPUs:", cpus)
    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

    for _, workers := range []int{1, 2, 4, cpus, cpus * 2} {
        count, duration := cpuWork(workers)
        fmt.Printf("CPU workers=%d primes=%d time=%v\n",
            workers, count, duration)
    }

    for _, workers := range []int{1, 10, 1000} {
        fmt.Printf("I/O workers=%d time=%v\n",
            workers, ioWork(workers))
    }
}