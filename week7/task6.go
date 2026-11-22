package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

type Record struct {
	ID   int
	Data [256]byte
}

var live []*Record
var latest *Record

//go:noinline
func consume(r *Record) int {
	// Storing the pointer globally forces
	// the Record to escape to the heap.
	latest = r

	return r.ID + int(r.Data[0])
}

func main() {
	nLive := flag.Int(
		"live", 80000, "Number of retained records",
	)

	nGarbage := flag.Int(
		"garbage", 10000000, "Temporary records",
	)

	flag.Parse()

	start := time.Now()

	// Objects that stay reachable.
	live = make([]*Record, 0, *nLive)

	for i := 0; i < *nLive; i++ {
		live = append(live, &Record{ID: i})
	}

	var sum int64

	// Short-lived heap objects.
	for i := 0; i < *nGarbage; i++ {
		r := &Record{ID: i}
		sum += int64(consume(r))
	}

	latest = nil

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("Sum:", sum)
	fmt.Println("Time:", time.Since(start))

	fmt.Println("GC cycles:", m.NumGC)

	fmt.Printf(
		"GC pauses: %.3f ms\n",
		float64(m.PauseTotalNs)/1e6,
	)

	fmt.Printf(
		"HeapAlloc: %.1f MiB\n",
		float64(m.HeapAlloc)/(1<<20),
	)

	fmt.Printf(
		"HeapSys: %.1f MiB\n",
		float64(m.HeapSys)/(1<<20),
	)

	fmt.Printf(
		"Sys: %.1f MiB\n",
		float64(m.Sys)/(1<<20),
	)

	runtime.KeepAlive(live)
}
