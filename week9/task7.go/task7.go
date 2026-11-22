package main

const N = 10_000

func passN(size int) int {
    ch := make(chan int, size)

    go func() {
        defer close(ch)

        for i := 0; i < N; i++ {
            ch <- i
        }
    }()

    total := 0

    for value := range ch {
        total += value
    }

    return total
}

func passBatches(size int) int {
    ch := make(chan []int, size)

    go func() {
        defer close(ch)

        for i := 0; i < N; i += 100 {
            batch := make([]int, 0, 100)

            for j := i; j < i+100 && j < N; j++ {
                batch = append(batch, j)
            }

            ch <- batch
        }
    }()

    total := 0

    for batch := range ch {
        for _, value := range batch {
            total += value
        }
    }

    return total
}