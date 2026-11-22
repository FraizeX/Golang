package main

import "testing"

var sinkInt int

var ids200 = func() []int {
	ids := make([]int, 200)

	for i := range ids {
		ids[i] = i
	}

	return ids
}()

func TestRender(t *testing.T) {
	newLength := RenderNew(ids200)
	poolLength := RenderPool(ids200)

	if newLength != poolLength {
		t.Fatalf(
			"Different lengths: %d vs %d",
			newLength,
			poolLength,
		)
	}

	if newLength != 1290 {
		t.Fatalf(
			"Got %d bytes, want 1290",
			newLength,
		)
	}
}

func BenchmarkRenderNew(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkInt = RenderNew(ids200)
	}
}

func BenchmarkRenderPool(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkInt = RenderPool(ids200)
	}
}
