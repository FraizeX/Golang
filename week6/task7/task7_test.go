package task7

import (
	"errors"
	"sync"
	"testing"
)

var ErrNegative = errors.New("negative value")

//go:noinline
func checkErr(number int) (int, error) {
	if number < 0 {
		return 0, ErrNegative
	}

	return number, nil
}

//go:noinline
func checkPanic(number int) int {
	if number < 0 {
		panic(ErrNegative)
	}

	return number
}

func safeCheck(number int) (value int, err error) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			err = panicValue.(error)
		}
	}()

	return checkPanic(number), nil
}

type Counter struct {
	mu sync.Mutex
	n  int
}

func (counter *Counter) IncNoDefer() {
	counter.mu.Lock()

	counter.n++

	counter.mu.Unlock()
}

func (counter *Counter) IncDefer() {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	counter.n++
}

// Результаты остаются наблюдаемыми для компилятора.
var sinkInt int
var sinkError error

func BenchmarkErrorReturn(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkInt, sinkError = checkErr(-1)
	}
}

func BenchmarkPanicRecover(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkInt, sinkError = safeCheck(-1)
	}
}

func BenchmarkIncNoDefer(b *testing.B) {
	counter := &Counter{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.IncNoDefer()
	}

	b.StopTimer()

	sinkInt = counter.n
}

func BenchmarkIncDefer(b *testing.B) {
	counter := &Counter{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.IncDefer()
	}

	b.StopTimer()

	sinkInt = counter.n
}