package main

import "sync"

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

func CountSeq(n int) int {
    count := 0

    for x := 0; x < n; x++ {
        if isPrime(x) {
            count++
        }
    }

    return count
}

func CountPar(n, parts int) int {
    if parts <= 0 {
        panic("parts must be positive")
    }

    results := make([]int, parts)
    var wg sync.WaitGroup

    for i := 0; i < parts; i++ {
        from := i * n / parts
        to := (i + 1) * n / parts

        wg.Add(1)

        go func(index, from, to int) {
            defer wg.Done()

            count := 0
            for x := from; x < to; x++ {
                if isPrime(x) {
                    count++
                }
            }

            results[index] = count
        }(i, from, to)
    }

    wg.Wait()

    total := 0
    for _, value := range results {
        total += value
    }

    return total
}