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
		NewTraits(defWords)
	}
}

// NewState()
func Benchmark_NewState(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		NewState(defWords)
	}
}

// Words()
func Benchmark_Words(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		Words(defWords)
	}
}

// Words() with a larger source dataset
func Benchmark_Words_LargeDataset(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		Words(testManyWords)
	}
}

// WordsN()
func Benchmark_WordsN(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		WordsN(defWords, defCount)
	}
}

// WordsN() with a larger source dataset
func Benchmark_WordsN_LargeDataset(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		WordsN(testManyWords, defCount)
	}
}

// Traits.Words()
func Benchmark_Traits_Words(b *testing.B) {
	// b.SkipNow()

	traits, _ := NewTraits(defWords)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		traits.Words()
	}
}

// NewState() && State.Words()
func Benchmark_NewState_State_Words(b *testing.B) {
	// b.SkipNow()

	for i := 0; i < b.N; i++ {
		state, _ := NewState(defWords)
		state.Words()
	}
}

// State.WordsN()
func Benchmark_State_WordsN(b *testing.B) {
	// b.SkipNow()

	state, _ := NewState(defWords)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if len(state.WordsN(defCount)) == 0 {
			state, _ = NewState(defWords)
			state.WordsN(defCount)
		}
	}
}

// State.WordsN() with a larger source dataset
func Benchmark_State_WordsN_LargeDataset(b *testing.B) {
	// b.SkipNow()

	state, _ := NewState(testManyWords)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if len(state.WordsN(defCount)) == 0 {
			state, _ = NewState(testManyWords)
			state.WordsN(defCount)
		}
	}
}
