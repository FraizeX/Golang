package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Specify case")
        return
    }

    switch os.Args[1] {

    // NIL CHANNEL

    case "nil-send":
        var ch chan int
        ch <- 1

    case "nil-receive":
        var ch chan int
        fmt.Println(<-ch)

    case "nil-close":
        var ch chan int
        close(ch)

    case "nil-lencap":
        var ch chan int
        fmt.Println(len(ch), cap(ch))

    // OPEN CHANNEL

    case "open-send":
        ch := make(chan int)

        go func() {
            time.Sleep(100 * time.Millisecond)
            fmt.Println("received:", <-ch)
        }()

        ch <- 42
        fmt.Println("send completed")
        time.Sleep(10 * time.Millisecond)

    case "open-receive":
        ch := make(chan int)

        go func() {
            time.Sleep(100 * time.Millisecond)
            ch <- 42
        }()

        fmt.Println("received:", <-ch)

    case "open-close":
        ch := make(chan int)
        close(ch)
        fmt.Println("closed successfully")

    case "open-lencap":
        ch := make(chan int, 3)
        ch <- 10
        ch <- 20

        fmt.Println(len(ch), cap(ch))

    // CLOSED CHANNEL

    case "closed-send":
        ch := make(chan int)
        close(ch)
        ch <- 42

    case "closed-receive":
        ch := make(chan int, 1)
        ch <- 7
        close(ch)

        v, ok := <-ch
        fmt.Println(v, ok)

        v, ok = <-ch
        fmt.Println(v, ok)

    case "closed-close":
        ch := make(chan int)
        close(ch)
        close(ch)

    case "closed-lencap":
        ch := make(chan int, 3)
        ch <- 10
        close(ch)

        fmt.Println(len(ch), cap(ch))

    default:
        fmt.Println("Unknown case")
    }
}