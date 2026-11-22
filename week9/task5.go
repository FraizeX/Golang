package main

import "fmt"

func main() {
    a := make(chan int, 1000)
    b := make(chan int, 1000)

    for i := 0; i < 1000; i++ {
        a <- 1
        b <- 2
    }

    countA, countB := 0, 0

    for i := 0; i < 1000; i++ {
        select {
        case <-a:
            countA++
        case <-b:
            countB++
        }
    }

    fmt.Println("A:", countA, "B:", countB)
}