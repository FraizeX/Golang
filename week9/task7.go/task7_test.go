package main

import (
    "fmt"
    "testing"
)

var resultSink int

func TestPassN(t *testing.T) {
    expected := N * (N - 1) / 2

    for _, size := range []int{0, 1, 100, 10000} {
        got := passN(size)

        if got != expected {
            t.Errorf(
                "buffer=%d got=%d want=%d",
                size, got, expected,
            )
        }
    }
}

func TestPassBatches(t *testing.T) {
    expected := N * (N - 1) / 2

    for _, size := range []int{0, 1, 100, 10000} {
        if got := passBatches(size); got != expected {
            t.Errorf("buffer=%d got=%d", size, got)
        }
    }
}

func BenchmarkPassN(b *testing.B) {
    for _, size := range []int{0, 1, 100, 10000} {
        b.Run(fmt.Sprintf("buf=%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                resultSink = passN(size)
            }
        })
    }
}

func BenchmarkPassBatches(b *testing.B) {
    for _, size := range []int{0, 1, 100, 10000} {
        b.Run(fmt.Sprintf("buf=%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                resultSink = passBatches(size)
            }
        })
    }
}