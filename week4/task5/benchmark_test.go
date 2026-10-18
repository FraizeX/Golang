package main

import "testing"

type Report struct {
	Title  string
	Values [1024]float64
}

//go:noinline
func (r Report) FirstValue() float64 {
	return r.Values[0]
}

//go:noinline
func (r *Report) FirstPointer() float64 {
	return r.Values[0]
}

type Point struct {
	X float64
	Y float64
}

//go:noinline
func (p Point) FirstValue() float64 {
	return p.X
}

//go:noinline
func (p *Point) FirstPointer() float64 {
	return p.X
}

var result float64

func BenchmarkReportValue(b *testing.B) {
	r := Report{}

	for i := 0; i < b.N; i++ {
		result = r.FirstValue()
	}
}

func BenchmarkReportPointer(b *testing.B) {
	r := Report{}

	for i := 0; i < b.N; i++ {
		result = r.FirstPointer()
	}
}

func BenchmarkPointValue(b *testing.B) {
	p := Point{}

	for i := 0; i < b.N; i++ {
		result = p.FirstValue()
	}
}

func BenchmarkPointPointer(b *testing.B) {
	p := Point{}

	for i := 0; i < b.N; i++ {
		result = p.FirstPointer()
	}
}