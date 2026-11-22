package task3

import (
	"fmt"
	"strconv"
	"strings"
)

// Slow version: repeated allocations.
func JoinIDsSlow(ids []int) string {
	s := ""

	for i, id := range ids {
		if i > 0 {
			s += ","
		}

		s += fmt.Sprint(id)
	}

	return s
}

// Fast version: reuse one buffer.
func JoinIDsFast(ids []int) string {
	var sb strings.Builder

	sb.Grow(len(ids) * 8)

	var tmp [20]byte

	for i, id := range ids {
		if i > 0 {
			sb.WriteByte(',')
		}

		digits := strconv.AppendInt(
			tmp[:0],
			int64(id),
			10,
		)

		sb.Write(digits)
	}

	return sb.String()
}
