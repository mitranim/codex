package codex

/**
 * The Traits type defines traits that characterise a word or group of words.
 * A valid Traits object can produce a generator that makes random
 * non-repeating words derived from the traits.
 *
 * This module also provides static functions to analyse words and extract
 * their traits.
 */

/*********************************** Type ************************************/

// Traits are rudimental characteristics of a word or group of words. A traits
// object unequivocally defines an unordered set of synthetic words that may be
// derived from it (see the definitions below).
type Traits struct {
	// Minimum and maximum number of sounds.
	MinNSounds int
	MaxNSounds int
	// Minimum and maximum number of vowels.
	MinNVowels int
	MaxNVowels int
	// Maximum number of consequtive vowels.
	MaxConseqVow int
	// Maximum number of consequtive consonants.
	MaxConseqCons int
	// Set of sounds that occur in the words.
	SoundSet Set
	// Set of pairs of sounds that occur in the words.
	PairSet PairSet

	// Replacement sound set to use instead of the default `knownSounds`.
	KnownSounds Set
	// Replacement sound set to use instead of the default `knownVowels`.
	KnownVowels Set
}

/**
 * Definitions of associated values.
 *
 * 1. Word sets.
 *
 * Given a traits object, we can derive an infinite set of sound sequences
 * using its sounds. Its limited subset qualifies as valid partial words,
 * defined by the Traits.validPart() criteria. This, in turn, has a limited
 * subset of sequences that qualify as valid complete words, defined by the
 * Traits.validComplete() criteria. This latter set is what we mean when
 * talking about the set of words defined by a traits group.
 *
 * 2. Virtual tree.
 *
 * Given a traits object, we can unequivocally derive an unordered tree that
 * efficiently represents the set of words defined by the traits. Such a tree
 * consists of nodes whose values represent sounds. Each node's path, starting
 * at the root, is unique, and represents a sequence of sounds, which may
 * qualify as a partial or complete word. The size of the tree is limited by
 * the criteria defined in Traits.validPart(), which puts a cap on the growth
 * of each branch.
 *
 * The tree doesn't have to exist in memory in order for us to traverse it.
 * For the sake of performance, we avoid building the entire tree, and instead
 * traverse its virtual equivalent through methods of a `state` object.
 */

/********************************** Methods **********************************/

/*--------------------------------- Public ----------------------------------*/

// Examines a slice of words and merges their traits into self.
func (this *Traits) Examine(words []string) error {
	if this == nil {
		return errType("can't examine with nil pointer")
	}

	// Examine each word and merge traits.
	for _, word := range words {
		if err := this.examineWord(word); err != nil {
			return err
		}
	}

	return nil
}

// Creates a generator function that returns a new word on each call. The words
// are guaranteed to never repeat and be randomly distributed in the traits'
// word set. When the set is exhausted, further calls return "".
func (this *Traits) Generator() func() string {
	st := &state{traits: this}
	return func() (result string) {
		// The trip iterator is not invoked after the state has been exhausted.
		st.trip(func(sounds ...string) { result = join(sounds, "") })
		return
	}
}

/*--------------------------------- Private ---------------------------------*/

// Takes a word, extracts its characteristics, and merges them into self. If the
// word doesn't satisfy our limitations, returns an error.
func (this *Traits) examineWord(word string) error {
	if this == nil {
		return errType("can't examine with nil pointer")
	}

	// Make sure the length is okay.
	if !validLength(word) {
		return errType("the word is too short or too long")
	}

	// Split into sounds.
	sounds, err := getSounds(word, this.knownSounds())
	if err != nil {
		return err
	}

	// Mandate that at least two sounds are found.
	if len(sounds) < 2 {
		return errType("less than two sounds found")
	}

	// Merge min and max number of consequtive sounds.
	n := len(sounds)
	if this.MinNSounds == 0 || n < this.MinNSounds {
		this.MinNSounds = n
	}
	if n > this.MaxNSounds {
		this.MaxNSounds = n
	}

	// Merge min and max total number of vowels.
	n = this.countVowels(sounds)
	if this.MinNVowels == 0 || n < this.MinNVowels {
		this.MinNVowels = n
	}
	if n > this.MaxNVowels {
		this.MaxNVowels = n
	}

	// Merge max number of consequtive vowels.
	n = this.maxConsequtiveVowels(sounds)
	if n > this.MaxConseqVow {
		this.MaxConseqVow = n
	}

	// Merge max number of consequtive consonants.
	n = this.maxConsequtiveConsonants(sounds)
	if n > this.MaxConseqCons {
		this.MaxConseqCons = n
	}

	// Merge set of used sounds.
	if this.SoundSet == nil {
		this.SoundSet = Set.New(nil, sounds...)
	} else {
		for sound := range Set.New(nil, sounds...) {
			this.SoundSet.Add(sound)
		}
	}

	// Find set of pairs of sounds.
	if this.PairSet == nil {
		this.PairSet = getPairs(sounds)
	} else {
		for pair := range getPairs(sounds) {
			this.PairSet.Add(pair)
		}
	}

	/*
		// Disabled for now; this causes a combinatorial explosion so bad that test
		// duration goes from seconds to minutes, if not hours. We should add an
		// option to enable this for the `WordsN()` static function.

		// Add reverse pairs.
		addReversePairs(this.PairSet)
	*/

	return nil
}

// Returns either the set of known sounds associated with the traits, or the
// default known sounds.
func (this *Traits) knownSounds() Set {
	if len(this.KnownSounds) > 0 {
		return this.KnownSounds
	}
	return knownSounds
}

// Returns either the set of known vowels associated with the traits, or the
// default known vowels.
func (this *Traits) knownVowels() Set {
	if len(this.KnownVowels) > 0 {
		return this.KnownVowels
	}
	return knownVowels
}

// Checks whether the given combination of sounds satisfies the conditions for
// a partial word. This is defined as follows:
//   1) the sounds don't exceed any of the numeric criteria in the given traits;
//   2) if there's only one sound, it must be the first sound in at least one
//      of the sound pairs in the given traits;
//   3) if there's at least one pair, the sequence of pairs must be valid as
//      defined in Traits.validPairs.
func (this *Traits) validPart(sounds ...string) bool {
	// Check numeric criteria.
	if this.countVowels(sounds) > this.MaxNVowels ||
		this.maxConsequtiveVowels(sounds) > this.MaxConseqVow ||
		this.maxConsequtiveConsonants(sounds) > this.MaxConseqCons {
		return false
	}

	// If there's only one sound, check if it's among the first sounds of pairs.
	if len(sounds) == 1 {
		for pair := range this.PairSet {
			if pair[0] == sounds[0] {
				return true
			}
		}
	}

	// Check if the pair sequence is valid per Traits.validPairs.
	if len(sounds) > 1 && !this.validPairs(sounds) {
		return false
	}

	return true
}

// Checks whether the given sequence of sounds satisfies the criteria for a
// complete word. This is defined as follows:
//   1) the sequence satisfies the partial criteria per Traits.validPart();
//   2) the sequence satisfies the complete criteria per Traits.checkPart().
func (this *Traits) validComplete(sounds ...string) bool {
	return this.validPart(sounds...) && this.checkPart(sounds...)
}

// Takes a valid partial word and checks if it's also a valid complete word,
// using the following criteria:
//   1) the number of vowels must fit within the bounds;
//   2) the number of sounds must fit within the bounds.
// The behaviour of this method for input values other than partial words is
// undefined.
func (this *Traits) checkPart(sounds ...string) bool {
	// Check vowel count.
	nVow := this.countVowels(sounds)
	if nVow < this.MinNVowels || nVow > this.MaxNVowels {
		return false
	}
	// Check sound count.
	if len(sounds) < this.MinNSounds || len(sounds) > this.MaxNSounds {
		return false
	}
	return true
}

// Verifies the validity of the sequence of sound pairs comprising the given
// word. Defined as follows:
//   1) the sequence must consist of sound pairs in the given traits; this is
//      implicitly guaranteed by the current tree traversal algorithms, so we
//      skip this check to save performance;
//   2) no sound pair immediately follows itself (e.g. "tata" in "ratatater");
//   3) no sound pair occurs more than twice.
// This has been somewhat optimised. Might stand for further improvement.
func (this *Traits) validPairs(sounds []string) bool {
	if len(sounds) < 2 {
		return true
	}

	// Variables to keep track of the last three pairs, up to current. This is
	// used for checking condition (2).
	var secondLastPair, lastPair, pair [2]string

	// Loop over the sequence, checking each condition.
	var prev string
	for index, current := range sounds {
		if index == 0 {
			prev = current
			continue
		}

		secondLastPair, lastPair, pair = lastPair, pair, [2]string{prev, current}

		// Check for condition (2). This can only be done starting at index 3.
		if index >= 3 {
			if secondLastPair == pair {
				return false
			}
		}

		// Check for condition (3). Originally we used a map of pairs to count pair
		// occurrences. This version is a performance optimisation, runs about
		// several dozen times faster for small datasets.
		if countPair(sounds[:index], prev, current) > 2 {
			return false
		}

		prev = current
	}

	return true
}

// Returns the biggest number of consequtive vowels that occurs in the given
// sound sequence.
func (this *Traits) maxConsequtiveVowels(sounds []string) int {
	return maxConsequtiveVowels(sounds, this.knownVowels())
}

// Returns the biggest number of consequtive consonants that occurs in the given
// sound sequence.
func (this *Traits) maxConsequtiveConsonants(sounds []string) int {
	return maxConsequtiveConsonants(sounds, this.knownVowels())
}

// Counts how many vowels occur in the given sound sequence.
func (this *Traits) countVowels(sounds []string) int {
	return countIntersections(sounds, this.knownVowels())
}

/********************************** Statics **********************************/

/*--------------------------------- Public ----------------------------------*/

// Shortcut to creating a traits object and calling its Traits.Examine().
func NewTraits(words []string) (*Traits, error) {
	traits := new(Traits)
	if err := traits.Examine(words); err != nil {
		return nil, err
	}
	return traits, nil
}
