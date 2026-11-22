package main

import (
	"bytes"
	"strconv"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		b := new(bytes.Buffer)
		b.Grow(4096)
		return b
	},
}

// Common rendering function.
func writeIDs(buf *bytes.Buffer, ids []int) {
	var tmp [20]byte

	for _, id := range ids {
		buf.WriteString("id=")

		digits := strconv.AppendInt(
			tmp[:0],
			int64(id),
			10,
		)

		buf.Write(digits)
		buf.WriteByte('\n')
	}
}

// Creates a new buffer every time.
func RenderNew(ids []int) int {
	var buf bytes.Buffer
	buf.Grow(4096)

	writeIDs(&buf, ids)

	return buf.Len()
}

// Reuses a buffer from sync.Pool.
func RenderPool(ids []int) int {
	buf := bufPool.Get().(*bytes.Buffer)

	buf.Reset()

	writeIDs(buf, ids)

	n := buf.Len()

	bufPool.Put(buf)

	return n
}
