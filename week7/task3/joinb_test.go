package task3

import "testing"

var sinkS string

var ids = []int{
	1, 2, 3, 4, 5,
	6, 7, 8, 9, 10,
}

func TestJoinIDsAllocs(t *testing.T) {
	slow := testing.AllocsPerRun(1000, func() {
		sinkS = JoinIDsSlow(ids)
	})

	fast := testing.AllocsPerRun(1000, func() {
		sinkS = JoinIDsFast(ids)
	})

	t.Logf("Slow: %.0f allocations", slow)
	t.Logf("Fast: %.0f allocations", fast)

	if fast > 1 {
		t.Fatalf(
			"Fast made %.0f allocations; want <= 1",
			fast,
		)
	}

	if JoinIDsSlow(ids) != JoinIDsFast(ids) {
		t.Fatal("Results are different")
	}
}

func BenchmarkJoinIDsSlow(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkS = JoinIDsSlow(ids)
	}
}

func BenchmarkJoinIDsFast(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sinkS = JoinIDsFast(ids)
	}
}
