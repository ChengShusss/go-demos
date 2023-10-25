package main

import "testing"

func BenchmarkGetAlive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Call the function to be tested here
		getAlive()
	}
}
