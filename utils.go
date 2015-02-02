package codex

// Public and private utility functions.

import (
	"math/rand"
	"regexp"
	"time"
)

/***************************** Public Functions ******************************/

// Static generator functions exposed by the package.

// Takes a sample group of words, analyses their traits, and builds a set of all
// synthetic words that may derived from those traits. This should only be used
// for very small samples. More than just a handful of sample words causes a
// combinatorial explosion, takes a lot of time to calculate, and produces too
// many results to be useful. The number of results can easily reach hundreds of
// thousands for just a dozen of sample words.
func Words(words []string) (Set, error) {
	traits, err := NewTraits(words)
	if err != nil {
		return nil, err
	}
	return traits.Words(), nil
}

// Takes a sample group of words and a count limiter. Analyses the words and
// builds a random sample of synthetic words that may be derived from those
// traits, limited to the given count.
func WordsN(words []string, num int) (Set, error) {
	state, err := NewState(words)
	if err != nil {
		return nil, err
	}
	return state.WordsN(num), nil
}

/********************************** Globals **********************************/

// Regexp to test validity of source words.
var matcher = regexp.MustCompile(`^\w+$`)

/********************************* Utilities *********************************/

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Takes a word and splits it into a series of known glyphs representing sounds.
func getSounds(word string) ([]string, error) {
	sounds := make([]string, 0, len(word))
	// Loop over the word, matching known glyphs. Break if no match is found.
	for i := 0; i < len(word); i++ {
		// Check for a known digraph.
		if i+2 <= len(word) && knownSounds.Has(word[i:i+2]) {
			sounds = append(sounds, word[i:i+2])
			i++
			// Check for a known monograph.
		} else if knownSounds.Has(word[i : i+1]) {
			sounds = append(sounds, word[i:i+1])
			// Otherwise return an error.
		} else {
			return nil, errType("encountered unknown symbol")
		}
	}
	// Return the found glyphs.
	return sounds, nil
}

// Takes a sequence of sounds and returns the set of consequtive pairs that
// occur in this sequence.
func getPairs(sounds []string) (pairs PairSet) {
	for i := 0; i < len(sounds)-1; i++ {
		pairs.Add([2]string{sounds[i], sounds[i+1]})
	}
	return
}

// Takes a set of pairs of sounds and adds their reverses.
func addReversePairs(pairs PairSet) {
	for key := range pairs {
		pairs.Add([2]string{key[1], key[0]})
	}
}

// Checks if the given word satisfies the following conditions:
//   1) only lowercase letters of the English alphabet;
//   2) no shorter than 2 letters and no longer than 64 letters.
func validInput(word string) bool {
	if len(word) < 2 || len(word) > 64 || !matcher.Match([]byte(word)) {
		return false
	}
	return true
}

// Republished rand.Perm.
func permutate(length int) []int {
	return rand.Perm(length)
}

func randFirsts(pairs PairSet) (results []string) {
	buffer := make([]string, 0, len(pairs))
outer:
	for pair := range pairs {
		// If existing value, skip.
		for _, value := range buffer {
			if pair[0] == value {
				continue outer
			}
		}
		// Otherwise add new value.
		buffer = append(buffer, pair[0])
	}
	results = make([]string, 0, len(buffer))
	for _, index := range permutate(len(buffer)) {
		results = append(results, buffer[index])
	}
	return
}

func randSeconds(pairs PairSet, first string) (results []string) {
	buffer := make([]string, 0, len(pairs))
	for pair := range pairs {
		// If doesn't match the given first value, skip.
		if pair[0] != first {
			continue
		}
		// Otherwise add new value.
		buffer = append(buffer, pair[1])
	}
	results = make([]string, 0, len(buffer))
	for _, index := range permutate(len(buffer)) {
		results = append(results, buffer[index])
	}
	return
}

// Returns the "invalid" error for the given word.
func errInvalid(word string) error {
	return errType("the word `" + word + "` is either too long, too short, or contains symbols other than lowercase Latin letters")
}

// Panic message used when breaking out from recursive iterations early.
const panicMsg = "early exit through panic"

// Wrapper for panic used when breaking out from recursive iterations early.
func interrupt() {
	panic(panicMsg)
}

// Wrapper for recovery from early iteration breakout through panic.
func aid() {
	msg := recover()
	if msg != nil && msg != panicMsg {
		panic(msg)
	}
}

/********************************** PairSet **********************************/

// PairSet behaves like a set of pairs of strings.
type PairSet map[[2]string]struct{}

// Creates a new set from the given keys. Usage:
//   PairSet.New(nil, [2]string{"one", "other"})
func (PairSet) New(keys ...[2]string) PairSet {
	set := make(PairSet, len(keys))
	for _, key := range keys {
		set.Add(key)
	}
	return set
}

// Adds the given element.
func (this *PairSet) Add(key [2]string) {
	if *this == nil {
		*this = PairSet{}
	}
	(*this)[key] = struct{}{}
}

// Deletes the given element.
func (this *PairSet) Del(key [2]string) {
	delete((*this), key)
}

// Checks for the presence of the given element.
func (this *PairSet) Has(key [2]string) bool {
	_, ok := (*this)[key]
	return ok
}

/*
// Commented out to avoid depending on fmt. If we include fmt at some point,
// this should be uncommented.

// Prints itself nicely in fmt(%#v).
func (this PairSet) GoString() string {
	keys := make([]string, 0, len(this))
	for key := range this {
		keys = append(keys, fmt.Sprintf("{%#v, %#v}", key[0], key[1]))
	}
	return "{" + strings.Join(keys, ", ") + "}"
}

// Prints itself nicely in println().
func (this PairSet) String() string {
	return this.GoString()
}
*/

/********************************** errType **********************************/

type errType string

func (this errType) Error() string {
	return string(this)
}
