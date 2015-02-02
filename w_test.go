package codex

// Tests.

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

/********************************** Globals **********************************/

var testWords = []string{
	"go", "nebula", "aurora", "theron", "thorax", "deity", "quasar",
}

var testManyWords = []string{
	"go", "nebula", "aurora", "theron", "thorax", "deity",
	"quasar", "graphene", "nanite", "orchestra", "eridium",
}

// Words with the following traits:
//   * 5-6 sounds;
//   * 2 vowels;
//   * 1 max consequtive vowel;
//   * 1-2 max consequtive consonants.
var testLimitedWords = []string{
	"theron", "thorax", "rocket", "proton", "filler", "absurd", "paper",
}

var defWords = testWords

const defCount = 12

/*********************************** Tests ***********************************/

// NewTraits()
func Test_NewTraits(t *testing.T) {
	// t.SkipNow()

	traits, err := NewTraits(testLimitedWords)
	tmust(t, err)

	if traits == nil {
		t.Fatal("!! missing traits object")
	}

	// MinNSounds
	if traits.MinNSounds != 5 {
		t.Fatalf("!! MinNSounds mismatch: expected %v, got %v", 5, traits.MinNSounds)
	}

	// MaxNSounds
	if traits.MaxNSounds != 6 {
		t.Fatalf("!! MaxNSounds mismatch: expected %v, got %v", 6, traits.MaxNSounds)
	}

	// MinNVowels
	if traits.MinNVowels != 2 {
		t.Fatalf("!! MinNVowels mismatch: expected %v, got %v", 2, traits.MinNVowels)
	}

	// MaxNVowels
	if traits.MaxNVowels != 2 {
		t.Fatalf("!! MaxNVowels mismatch: expected %v, got %v", 2, traits.MaxNVowels)
	}

	// MaxConseqVow
	if traits.MaxConseqVow != 1 {
		t.Fatalf("!! MaxConseqVow mismatch: expected %v, got %v", 1, traits.MaxConseqVow)
	}

	// MaxConseqCons
	if traits.MaxConseqCons != 2 {
		t.Fatalf("!! MaxConseqCons mismatch: expected %v, got %v", 2, traits.MaxConseqCons)
	}

	// SoundSet
	sounds := Set{}
	for _, word := range testLimitedWords {
		sequence, err := getSounds(word)
		tmust(t, err)
		for _, sound := range sequence {
			sounds.Add(sound)
		}
	}
	if !reflect.DeepEqual(traits.SoundSet, sounds) {
		t.Fatalf("!! SoundSet mismatch")
	}

	// PairSet
	pairs := PairSet{}
	for _, word := range testLimitedWords {
		sequence, err := getSounds(word)
		tmust(t, err)
		for pair := range getPairs(sequence) {
			pairs.Add(pair)
		}
	}
	if !reflect.DeepEqual(traits.PairSet, pairs) {
		t.Fatalf("!! PairSet mismatch")
	}
}

// NewState()
func Test_NewState(t *testing.T) {
	// t.SkipNow()

	state, err := NewState(defWords)
	tmust(t, err)

	if state == nil {
		t.Fatal("!! missing state object")
	}
}

// Words()
func Test_Words(t *testing.T) {
	// t.SkipNow()
	test_Words(t, defWords)
}

// Words() with a larger source dataset
func Test_Words_LargeDataset(t *testing.T) {
	// t.SkipNow()
	test_Words(t, testManyWords)
}

// WordsN()
func Test_WordsN(t *testing.T) {
	// t.SkipNow()
	test_WordsN(t, defWords)
}

// WordsN() with a larger source dataset
func Test_WordsN_LargeDataset(t *testing.T) {
	// t.SkipNow()
	test_WordsN(t, testManyWords)
}

// Traits.Words()
func Test_Traits_Words(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)
	Test_Words(t)

	traits, _ := NewTraits(defWords)
	words, _ := Words(defWords)

	if !reflect.DeepEqual(words, traits.Words()) {
		t.Fatal("!! word set mismatch between Words() and Traits.Words()")
	}
}

// State.Words()
func Test_State_Words(t *testing.T) {
	// t.SkipNow()

	Test_NewState(t)
	Test_Words(t)

	state, _ := NewState(defWords)
	words, _ := Words(defWords)

	if !reflect.DeepEqual(words, state.Words()) {
		t.Fatal("!! word set mismatch between Words() and State.Words()")
	}
}

// State.WordsN()
func Test_State_WordsN(t *testing.T) {
	// t.SkipNow()

	Test_NewState(t)
	Test_Words(t)
	Test_State_Words(t)

	state, _ := NewState(defWords)
	words, _ := Words(defWords)

	// Ascertain that the method returns the expected number of words.
	sample := state.WordsN(defCount)
	if sample == nil {
		t.Fatal("!! missing sample set")
	}
	if len(sample) != defCount {
		t.Fatal("!! word count mismatch")
	}

	// Ascertain that no results are repeated and that subsequent calls eventually
	// exhaust the word pool.
	total := Set{}
	for len(sample) > 0 {
		for word := range sample {
			if total.Has(word) {
				t.Fatal("!! repeated word", word, "at total length:", len(total))
			}
			total.Add(word)
		}
		sample = state.WordsN(defCount)
	}

	// Ascertain that the total set returned from all calls is equivalent to the
	// set from Words(), which is shown to be equivalent to State.Words() in
	// another test.
	if !reflect.DeepEqual(words, total) {
		t.Fatal("!! word set mismatch between Words() and total from State.WordsN()")
	}
}

// Verifies that words from Words() match their traits.
func Test_Words_Match_Traits(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)
	Test_Words(t)

	traits, _ := NewTraits(testLimitedWords)
	words, _ := Words(testLimitedWords)

	test_Words_Match_Traits(t, traits, words)
}

// Verifies that words from State.Words() match their traits.
func Test_State_Words_Match_Traits(t *testing.T) {
	// t.SkipNow()

	Test_NewState(t)
	Test_State_Words(t)

	state, _ := NewState(testLimitedWords)

	test_Words_Match_Traits(t, state.Traits, state.Words())
}

// Verifies that NewTraits() and NewState() produce an error with invalid input.
func Test_Invalid_Input(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)
	Test_NewState(t)

	invalids := []string{
		"", "a", "CAPITALS", "Capitalised", "with space",
		"numbers134125", "łàtîñôñè", "кириллица",
	}

	for _, invalid := range invalids {
		traits, err := NewTraits([]string{invalid})
		if traits != nil || err == nil {
			t.Fatalf("!! expected nil traits and non-nil error")
		}
		state, err := NewState([]string{invalid})
		if state != nil || err == nil {
			t.Fatalf("!! expected nil state and non-nil error")
		}
	}
}

// Verifies that words from State.WordsN() are randomly distributed. Rudimental
// and naive, todo remember some math and use a real probability function.
func Test_State_WordsN_Random_Distribution(t *testing.T) {
	// t.SkipNow()

	Test_NewState(t)
	Test_Words(t)
	Test_State_WordsN(t)

	state, _ := NewState(defWords)

	// Make a sorted list of words.
	unordered, _ := Words(defWords)
	words := make([]string, 0, len(unordered))
	for word := range unordered {
		words = append(words, word)
	}
	sort.Strings(words)

	// Limit of how many tight groups to permit.
	maxTightGroups := len(words) / defCount / 10
	if maxTightGroups == 0 {
		maxTightGroups = 1
	}

	// Counter of tight group occurrences.
	count := 0

	// Loop over State.WordsN() results and count how many times all indices from
	// a sample fall within a tight range (let's say 1/5th the length).
	for sample := state.WordsN(defCount); len(sample) > 0; sample = state.WordsN(defCount) {
		indices := make([]int, 0, len(sample))
		for word := range sample {
			indices = append(indices, findIndex(words, word))
		}
		if maximum(indices)-minimum(indices) < len(words)/5 {
			count++
		}
	}

	if count > maxTightGroups {
		t.Fatalf("!! for %v sorted words, %v samples were tightly grouped", len(words), count)
	}
}

/********************************** Helpers **********************************/

// Words() helper.
func test_Words(t *testing.T, source []string) {
	words, err := Words(source)
	tmust(t, err)
	if words == nil {
		t.Fatal("!! missing words set")
	}
	if len(words) == 0 {
		t.Fatal("!! zero words received")
	}
	// The output for a dozen source words can easily reach tens of thousands of
	// results. We're being very conservative here.
	if len(words) < 100 {
		t.Fatal("!! unexpectedly small number of words:", len(words))
	}
	t.Log("-- total words in sample:", len(words))
	// t.Log("-- words in sample:", words)
}

// WordsN() helper.
func test_WordsN(t *testing.T, source []string) {
	words, err := WordsN(source, defCount)
	tmust(t, err)
	if words == nil {
		t.Fatal("!! missing words set")
	}
	if len(words) != defCount {
		t.Fatalf("!! word count mismatch: expected %v, got %v", defCount, len(words))
	}
}

// Words_Match_Traits helper.
func test_Words_Match_Traits(t *testing.T, traits *Traits, words Set) {
	for word := range words {
		// MinNSounds
		sounds, err := getSounds(word)
		tmust(t, err)
		if len(sounds) < traits.MinNSounds {
			t.Fatalf("!! \"%v\" MinNSounds mismatch: expected >=%v, got %v", word, traits.MinNSounds, len(sounds))
		}

		// MaxNSounds
		if len(sounds) > traits.MaxNSounds {
			t.Fatalf("!! \"%v\" MaxNSounds mismatch: expected <=%v, got %v", word, traits.MaxNSounds, len(sounds))
		}

		// MinNVowels
		if n := countIntersections(sounds, knownVowels); n < traits.MinNVowels {
			t.Fatalf("!! \"%v\" MinNVowels mismatch: expected >=%v, got %v", word, traits.MinNVowels, n)
		}

		// MaxNVowels
		if n := countIntersections(sounds, knownVowels); n > traits.MaxNVowels {
			t.Fatalf("!! \"%v\" MaxNVowels mismatch: expected <=%v, got %v", word, traits.MaxNVowels, n)
		}

		// MaxConseqVow
		if n := maxConsequtiveVowels(sounds); n > traits.MaxConseqVow {
			t.Fatalf("!! \"%v\" MaxConseqVow mismatch: expected <=%v, got %v", word, traits.MaxConseqVow, n)
		}

		// MaxConseqCons
		if n := maxConsequtiveConsonants(sounds); n > traits.MaxConseqCons {
			t.Fatalf("!! \"%v\" MaxConseqCons mismatch: expected <=%v, got %v", word, traits.MaxConseqCons, n)
		}

		// SoundSet
		for sound := range Set.New(nil, sounds...) {
			if !traits.SoundSet.Has(sound) {
				t.Fatalf("!! \"%v\" SoundSet mismatch, unexpected sound: %v", word, sound)
			}
		}

		// PairSet
		for pair := range getPairs(sounds) {
			if !traits.PairSet.Has(pair) {
				t.Fatalf("!! \"%v\" PairSet mismatch, unexpected pair: %v", word, pair)
			}
		}
	}
}

/*********************************** Utils ***********************************/

// Prints expanded values.
func prn(values ...interface{}) {
	var result string
	for i := 0; i < len(values); i++ {
		if reflect.ValueOf(values[i]).Kind() == reflect.String {
			result += fmt.Sprintf("%v", values[i])
		} else {
			result += fmt.Sprintf("%#v", values[i])
		}
		if i < len(values)-1 {
			result += " "
		}
	}
	fmt.Println(result)
}

// Prints simple values.
func log(values ...interface{}) {
	fmt.Println(values...)
}

func tmust(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// Finds the position of the given value in the given slice.
func findIndex(group []string, value string) int {
	for index, val := range group {
		if val == value {
			return index
		}
	}
	return -1
}

func minimum(values []int) int {
	if len(values) == 0 {
		return 0
	}
	result := values[0]
	for _, value := range values {
		if value < result {
			result = value
		}
	}
	return result
}

func maximum(values []int) int {
	if len(values) == 0 {
		return 0
	}
	result := values[0]
	for _, value := range values {
		if value > result {
			result = value
		}
	}
	return result
}
