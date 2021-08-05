package main

import (
	"io"
	"testing"
)

func benchmarkCalculate(input int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		run(io.Discard, input)
	}
}

func BenchmarkBuffer1(b *testing.B) {
	benchmarkCalculate(1, b)
}
