package codex

// Benchmarks.

import (
	"testing"
)

/********************************** Globals **********************************/

var defSounds = []string{"n", "e", "b", "u"}

/******************************** Benchmarks *********************************/

// NewTraits()
func Benchmark_NewTraits(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		NewTraits(testDefWords)
	}
}

// NewTraits()
func Benchmark_NewTraits_LargeDataset(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		NewTraits(testManyWords)
	}
}

// Traits.Generator() -> complete words set
func Benchmark_Generator_All(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		traits, _ := NewTraits(testDefWords)
		gen := traits.Generator()
		for gen() != "" {
		}
	}
}

// Large source dataset -> Traits.Generator() -> complete words set
func Benchmark_Generator_All_LargeDataset(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		traits, _ := NewTraits(testManyWords)
		gen := traits.Generator()
		for gen() != "" {
		}
	}
}

// Traits.Generator() -> generate default count
func Benchmark_Generator_N(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		traits, _ := NewTraits(testDefWords)
		gen := traits.Generator()
		for i := 0; i < testDefCount; i++ {
			gen()
		}
	}
}

// Large source dataset -> Traits.Generator() -> generate default count
func Benchmark_Generator_N_LargeDataset(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		traits, _ := NewTraits(testManyWords)
		gen := traits.Generator()
		for i := 0; i < testDefCount; i++ {
			gen()
		}
	}
}
