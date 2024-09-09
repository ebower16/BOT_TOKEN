package main

import (
	"testing"
)

func BenchmarkParallelSumOfSquares(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParallelSumOfSquares(10000000)
	}
}

func BenchmarkSumOfSquares(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SumOfSquares(10000000)
	}
}
