package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Simulates reading a 1 MiB file.
func readFile(i int) []byte {
	data := make([]byte, 1<<20)

	copy(data, fmt.Sprintf("HEADER-%04d", i))

	return data
}

// Keeps small slices of large arrays.
func keepHeadersSlice(n int) [][]byte {
	headers := make([][]byte, 0, n)

	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, data[:11])
	}

	return headers
}

// Copies only the required 11 bytes.
func keepHeadersCopy(n int) [][]byte {
	headers := make([][]byte, 0, n)

	for i := 0; i < n; i++ {
		data := readFile(i)

		headers = append(
			headers,
			bytes.Clone(data[:11]),
		)
	}

	return headers
}

// Creates a large string.
func bigString(i int) string {
	return string(readFile(i))
}

// Keeps small substrings of large strings.
func keepSubstrings(n int) []string {
	result := make([]string, 0, n)

	for i := 0; i < n; i++ {
		big := bigString(i)
		result = append(result, big[:5])
	}

	return result
}

// Copies only the necessary characters.
func keepStringClones(n int) []string {
	result := make([]string, 0, n)

	for i := 0; i < n; i++ {
		big := bigString(i)

		result = append(
			result,
			strings.Clone(big[:5]),
		)
	}

	return result
}

// Measures memory retained by the result.
func measure(
	name string,
	build func(int) any,
	n int,
) {
	runtime.GC()

	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	result := build(n)

	runtime.GC()

	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	delta := int64(after.HeapAlloc) -
		int64(before.HeapAlloc)

	fmt.Printf(
		"%-20s %+7.2f MiB\n",
		name,
		float64(delta)/(1<<20),
	)

	runtime.KeepAlive(result)
}

func main() {
	const n = 100

	measure("headers slices", func(n int) any {
		return keepHeadersSlice(n)
	}, n)

	measure("headers clones", func(n int) any {
		return keepHeadersCopy(n)
	}, n)

	measure("substrings", func(n int) any {
		return keepSubstrings(n)
	}, n)

	measure("string clones", func(n int) any {
		return keepStringClones(n)
	}, n)
}
