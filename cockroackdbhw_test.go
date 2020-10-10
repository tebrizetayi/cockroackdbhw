package main

import (
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("insert", func(bp *testing.B) { insert() })
	}
}

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("read", func(bp *testing.B) { read("ca5bd769-7b2c-40af-9d0b-5e370091ba8a") })
	}
}
