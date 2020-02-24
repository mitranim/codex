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

var testDefWords = testWords

const testDefCount = 12

/*********************************** Tests ***********************************/

// NewTraits()
func Test_NewTraits(t *testing.T) {
	// t.SkipNow()

	traits, err := NewTraits(testLimitedWords)
	tmust(t, err)

	if traits == nil {
		t.Fatal("missing traits object")
	}

	// MinNSounds
	if traits.MinNSounds != 5 {
		t.Fatalf("MinNSounds mismatch: expected %v, got %v", 5, traits.MinNSounds)
	}

	// MaxNSounds
	if traits.MaxNSounds != 6 {
		t.Fatalf("MaxNSounds mismatch: expected %v, got %v", 6, traits.MaxNSounds)
	}

	// MinNVowels
	if traits.MinNVowels != 2 {
		t.Fatalf("MinNVowels mismatch: expected %v, got %v", 2, traits.MinNVowels)
	}

	// MaxNVowels
	if traits.MaxNVowels != 2 {
		t.Fatalf("MaxNVowels mismatch: expected %v, got %v", 2, traits.MaxNVowels)
	}

	// MaxConseqVow
	if traits.MaxConseqVow != 1 {
		t.Fatalf("MaxConseqVow mismatch: expected %v, got %v", 1, traits.MaxConseqVow)
	}

	// MaxConseqCons
	if traits.MaxConseqCons != 2 {
		t.Fatalf("MaxConseqCons mismatch: expected %v, got %v", 2, traits.MaxConseqCons)
	}

	// SoundSet
	sounds := Set{}
	for _, word := range testLimitedWords {
		sequence, err := getSounds(word, traits.knownSounds())
		tmust(t, err)
		for _, sound := range sequence {
			sounds.Add(sound)
		}
	}
	if !reflect.DeepEqual(traits.SoundSet, sounds) {
		t.Fatalf("SoundSet mismatch")
	}

	// PairSet
	pairs := PairSet{}
	for _, word := range testLimitedWords {
		sequence, err := getSounds(word, traits.knownSounds())
		tmust(t, err)
		for pair := range getPairs(sequence) {
			pairs.Add(pair)
		}
	}
	if !reflect.DeepEqual(traits.PairSet, pairs) {
		t.Fatalf("PairSet mismatch")
	}
}

// Traits.Generator()
func Test_Traits_Generator(t *testing.T) {
	// t.SkipNow()

	traits, err := NewTraits(testDefWords)
	tmust(t, err)
	gen := traits.Generator()

	if gen == nil {
		t.Fatal("missing generator function")
	}

	if gen() == "" {
		t.Fatal("no output received from generator")
	}
}

// Checks a generator's output. Also verifies that a generator eventually
// exhausts its word set (will enter an infinite loop otherwise).
func Test_Generator(t *testing.T) {
	// t.SkipNow()

	Test_Traits_Generator(t)

	traits, _ := NewTraits(testDefWords)
	gen := traits.Generator()
	words := Set{}

	// Collect the total output, check each word's validity and uniqueness.
	for word := gen(); word != ""; word = gen() {
		sounds, err := getSounds(word, traits.knownSounds())
		tmust(t, err)
		if !traits.validComplete(sounds...) {
			t.Fatal("invalid output from generator:", word)
		}
		if words.Has(word) {
			t.Fatal("repeated output from generator:", word)
		}
		words.Add(word)
	}

	// The output for a dozen source words can easily reach tens of thousands of
	// results. We're being very conservative here.
	if len(words) < 100 {
		t.Fatal("unexpectedly small number of words:", len(words))
	}

	// t.Log("total words in sample:", len(words))
	// t.Log("words in sample:", words)
}

// Verifies that the words returned from a generator match its source traits.
func Test_Generator_Words_Match_Traits(t *testing.T) {
	// t.SkipNow()

	Test_Generator(t)

	traits, _ := NewTraits(testLimitedWords)
	words := collectAll(traits)

	test_Words_Match_Traits(t, traits, words)
}

// Verifies that NewTraits() produces an error with invalid input.
func Test_Invalid_Input(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)

	invalids := []string{
		"", "a", "CAPITALS", "Capitalised", "with space",
		"numbers134125", "łàtîñôñè", "кириллица",
	}

	for _, invalid := range invalids {
		traits, err := NewTraits([]string{invalid})
		if traits != nil || err == nil {
			t.Fatalf("expected nil traits and non-nil error, got %v and %v", traits, err)
		}
	}
}

// Verifies that Traits.Examine() and NewTraits() are equivalent.
func Test_Traits_Examine(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)

	traits := new(Traits)
	tmust(t, traits.Examine(testDefWords))

	other, _ := NewTraits(testDefWords)

	if !reflect.DeepEqual(traits, other) {
		t.Fatal("expected new(Traits) + Traits.Examine() to be equivalent to NewTraits()")
	}

	if !reflect.DeepEqual(collectAll(traits), collectAll(other)) {
		t.Fatal("expected resulting word sets to be equivalent")
	}
}

// Verifies that a Traits object uses internal known sounds, if available.
func Test_Traits_KnownSounds(t *testing.T) {
	// t.SkipNow()

	Test_Traits_Examine(t)

	traits := new(Traits)

	traits.KnownSounds = Set.New(nil,
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	)

	if reflect.DeepEqual(traits.knownSounds(), knownSounds) {
		t.Fatal("expected Traits.knownSounds() to return the internal sound set")
	}

	if traits.Examine(testDefWords) == nil {
		t.Fatal("expected Traits.Examine() to fail when used with a limited sound set")
	}
}

// Verifies that a Traits object uses internal known vowels, if available.
func Test_Traits_KnownVowels(t *testing.T) {
	// t.SkipNow()

	Test_NewTraits(t)
	Test_Traits_Examine(t)

	traits := new(Traits)

	traits.KnownVowels = Set.New(nil,
		"a", "e", "i",
	)

	if reflect.DeepEqual(traits.knownVowels(), knownVowels) {
		t.Fatal("expected Traits.knownVowels() to return the internal vowel set")
	}

	if traits.Examine(testDefWords) != nil {
		t.Fatal("expected Traits.Examine() to complete successfully when used with a limited vowel set")
	}

	other, _ := NewTraits(testDefWords)
	if reflect.DeepEqual(traits, other) {
		t.Fatalf("expected Traits.Examine() with custom vowels to produce different traits")
	}

	if reflect.DeepEqual(collectAll(traits), collectAll(other)) {
		t.Fatal("expected resulting word sets to be different")
	}
}

// Verifies that words from a generator are randomly distributed. Rudimental and
// naive, todo remember some math and use a real probability function.
func Test_Generator_Random_Distribution(t *testing.T) {
	// t.SkipNow()

	Test_Generator(t)

	traits, _ := NewTraits(testDefWords)

	// Make a sorted list of words.
	gen := traits.Generator()
	unordered := Set{}
	for word := gen(); word != ""; word = gen() {
		unordered.Add(word)
	}
	words := make([]string, 0, len(unordered))
	for word := range unordered {
		words = append(words, word)
	}
	sort.Strings(words)

	// Limit of how many tight groups to permit.
	maxTightGroups := len(words) / testDefCount / 10
	if maxTightGroups == 0 {
		maxTightGroups = 1
	}

	// Counter of tight group occurrences.
	count := 0

	// Prepare a generator that makes words in chunks.
	wordsN := generatorN(traits)

	// Loop over generator results and count how many times all indices from
	// a sample fall within a tight range (let's say 1/5th the length).
	for sample := wordsN(testDefCount); len(sample) > 0; sample = wordsN(testDefCount) {
		indices := make([]int, 0, len(sample))
		for word := range sample {
			indices = append(indices, findIndex(words, word))
		}
		if maximum(indices)-minimum(indices) < len(words)/5 {
			count++
		}
	}

	if count > maxTightGroups {
		t.Fatalf("for %v sorted words, %v out of %v samples were tightly grouped", len(words), count, len(words)/testDefCount+1)
	}
}

/********************************** Helpers **********************************/

// Words_Match_Traits helper.
func test_Words_Match_Traits(t *testing.T, traits *Traits, words Set) {
	for word := range words {
		// MinNSounds
		sounds, err := getSounds(word, traits.knownSounds())
		tmust(t, err)
		if len(sounds) < traits.MinNSounds {
			t.Fatalf("\"%v\" MinNSounds mismatch: expected >=%v, got %v", word, traits.MinNSounds, len(sounds))
		}

		// MaxNSounds
		if len(sounds) > traits.MaxNSounds {
			t.Fatalf("\"%v\" MaxNSounds mismatch: expected <=%v, got %v", word, traits.MaxNSounds, len(sounds))
		}

		// MinNVowels
		if n := traits.countVowels(sounds); n < traits.MinNVowels {
			t.Fatalf("\"%v\" MinNVowels mismatch: expected >=%v, got %v", word, traits.MinNVowels, n)
		}

		// MaxNVowels
		if n := traits.countVowels(sounds); n > traits.MaxNVowels {
			t.Fatalf("\"%v\" MaxNVowels mismatch: expected <=%v, got %v", word, traits.MaxNVowels, n)
		}

		// MaxConseqVow
		if n := traits.maxConsequtiveVowels(sounds); n > traits.MaxConseqVow {
			t.Fatalf("\"%v\" MaxConseqVow mismatch: expected <=%v, got %v", word, traits.MaxConseqVow, n)
		}

		// MaxConseqCons
		if n := traits.maxConsequtiveConsonants(sounds); n > traits.MaxConseqCons {
			t.Fatalf("\"%v\" MaxConseqCons mismatch: expected <=%v, got %v", word, traits.MaxConseqCons, n)
		}

		// SoundSet
		for sound := range Set.New(nil, sounds...) {
			if !traits.SoundSet.Has(sound) {
				t.Fatalf("\"%v\" SoundSet mismatch, unexpected sound: %v", word, sound)
			}
		}

		// PairSet
		for pair := range getPairs(sounds) {
			if !traits.PairSet.Has(pair) {
				t.Fatalf("\"%v\" PairSet mismatch, unexpected pair: %v", word, pair)
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

// Creates a function that generates words in chunks until its inner generator
// is exhausted.
func generatorN(traits *Traits) func(int) Set {
	gen := traits.Generator()

	return func(num int) Set {
		words := Set{}
		for word := gen(); word != "" && len(words) < num; word = gen() {
			words.Add(word)
		}
		return words
	}
}

// Collects all words from the given traits.
func collectAll(traits *Traits) Set {
	words := Set{}
	gen := traits.Generator()
	for word := gen(); word != ""; word = gen() {
		words.Add(word)
	}
	return words
}
