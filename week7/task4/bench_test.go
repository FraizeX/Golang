package task4

import "testing"

var (
	pointResult    Point
	pointPtrResult *Point
	sliceResult    []int
)

func BenchmarkNewPointValue(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		pointResult = NewPointValue(i, i+1)
	}
}

func BenchmarkNewPointPointer(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		pointPtrResult = NewPointPointer(i, i+1)
	}
}

func BenchmarkCollectNoCap(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sliceResult = CollectNoCap(1000)
	}
}

func BenchmarkCollectWithCap(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sliceResult = CollectWithCap(1000)
	}
}
