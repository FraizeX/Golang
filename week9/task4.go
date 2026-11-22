package main

import (
    "fmt"
    "runtime"
    "time"

    "golang.org/x/tour/tree"
)

func walkNode(t *tree.Tree, ch chan<- int) {
    if t == nil {
        return
    }

    walkNode(t.Left, ch)
    ch <- t.Value
    walkNode(t.Right, ch)
}

func Walk(t *tree.Tree, ch chan int) {
    defer close(ch)
    walkNode(t, ch)
}

func Same(t1, t2 *tree.Tree) bool {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go Walk(t1, ch1)
    go Walk(t2, ch2)

    same := true

    for {
        v1, ok1 := <-ch1
        v2, ok2 := <-ch2

        if !ok1 && !ok2 {
            return same
        }

        if ok1 != ok2 || v1 != v2 {
            same = false
        }
    }
}

func main() {
    fmt.Println(
        "Same trees:",
        Same(tree.New(1), tree.New(1)),
    )

    fmt.Println(
        "Different trees:",
        Same(tree.New(1), tree.New(2)),
    )

    time.Sleep(10 * time.Millisecond)

    fmt.Println(
        "Goroutines:",
        runtime.NumGoroutine(),
    )
}