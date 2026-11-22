package main

import (
    "fmt"
    "testing"
)

const N = 200_000

var sink int

func TestCountPar(t *testing.T) {
    expected := CountSeq(N)

    for _, parts := range []int{1, 2, 8, 1000, 200_000} {
        actual := CountPar(N, parts)

        if actual != expected {
            t.Errorf(
                "parts=%d: got %d, want %d",
                parts, actual, expected,
            )
        }
    }
}

func BenchmarkCountSeq(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sink = CountSeq(N)
    }
}

func BenchmarkCountPar(b *testing.B) {
    for _, parts := range []int{2, 8, 1000, 200_000} {
        b.Run(fmt.Sprintf("parts=%d", parts), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                sink = CountPar(N, parts)
            }
        })
    }
}