package main

import (
	"flag"
	"fmt"
)

//go:noinline
func depth64(n int) int {
	var pad [64]byte

	pad[n%len(pad)] = byte(n)

	if n == 0 {
		return int(pad[0])
	}

	return depth64(n-1) + 1 +
		int(pad[(n+1)%len(pad)])
}

//go:noinline
func depth1024(n int) int {
	var pad [1024]byte

	pad[n%len(pad)] = byte(n)

	if n == 0 {
		return int(pad[0])
	}

	return depth1024(n-1) + 1 +
		int(pad[(n+1)%len(pad)])
}

// Version without reading the array.
//go:noinline
func depthNoRead(n int) int {
	var pad [64]byte

	pad[n%len(pad)] = byte(n)

	if n == 0 {
		return 0
	}

	return depthNoRead(n-1) + 1
}

func forever(n int) int {
	return forever(n+1) + 1
}

func main() {
	mode := flag.String("mode", "64", "Test mode")
	n := flag.Int("n", 10000, "Recursion depth")
	flag.Parse()

	switch *mode {
	case "64":
		fmt.Println(depth64(*n))

	case "1024":
		fmt.Println(depth1024(*n))

	case "noread":
		fmt.Println(depthNoRead(*n))

	case "overflow":
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered:", r)
			}
		}()

		fmt.Println(forever(0))
	}
}
