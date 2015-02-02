package codex

// The Traits type defines traits that characterise a word or group of words.
// This module also provides functions to analyse words and extract their
// traits.

import (
	"strings"
)

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
}

/**
 * Definitions of associated values.
 *
 * 1. Word sets.
 *
 * Given a traits object, we can derive an infinite set of sound sequences
 * using its sounds. A subset of it qualifies as valid partial words, defined
 * by the Traits.validPart() criteria. This, in turn, has a limited subset of
 * sequences that qualify as valid complete words, defined by the
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
 * the criteria defined in Traits.validComplete(), which put a cap on the
 * growth of each branch.
 *
 * The tree doesn't have to exist in memory in order for us to traverse it.
 * For the sake of performance, we avoid building the entire tree, and instead
 * traverse its virtual equivalent.
 */

/********************************** Methods **********************************/

/*--------------------------------- Public ----------------------------------*/

// Generates and returns the entire set of words defined by the traits.
func (this *Traits) Words() (words Set) {
	iterator := func(sounds ...string) {
		words.Add(strings.Join(sounds, ""))
	}
	this.walk(iterator)
	return
}

/*--------------------------------- Private ---------------------------------*/

// Checks whether the given combination of sounds satisfies the conditions for
// a partial word. This is defined as follows:
//   1) the sounds don't exceed any of the numeric criteria in the given traits;
//   2) if there's only one sound, it must be the first sound in at least one
//      of the sound pairs in the given traits;
//   3) if there's at least one pair, the sequence of pairs must be valid as
//      defined in Traits.validPairs.
func (this *Traits) validPart(sounds ...string) bool {
	// Check numeric criteria.
	if countIntersections(sounds, knownVowels) > this.MaxNVowels ||
		maxConsequtiveVowels(sounds) > this.MaxConseqVow ||
		maxConsequtiveConsonants(sounds) > this.MaxConseqCons {
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

	// Checks if the pair sequence is valid per Traits.validPairs.
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
	nVow := countIntersections(sounds, knownVowels)
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
//   1) the sequence must consist of sound pairs in the given traits;
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

		// Check for condition (1).
		if !this.PairSet.Has(pair) {
			return false
		}

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

// Traverses the virtual tree of valid partial words associated with the traits
// in depth-first pre-order, calling the given iterator function on each path
// (from the root) that forms a complete word.
func (this *Traits) walk(iterator func(...string), sounds ...string) {
	if iterator == nil {
		return
	}

	// If no sounds were passed, start from the root.
	if len(sounds) == 0 {
		firsts := Set{}
		// The values of the first-level nodes are the first values of the pair set
		// associated with the traits.
		for pair := range this.PairSet {
			first := pair[0]
			// Check for repeats. The same first sound may occur in several pairs.
			if firsts.Has(first) {
				continue
			}
			firsts.Add(first)
			// Continue recursively.
			this.walk(iterator, first)
		}
		// If sounds were passed, continue from that path onward.
	} else {
		// [ ... sounds ... ( last sound ] <- pair -> next sound )
		//
		// We investigate pairs that begin with the last sound of the given
		// preceding sounds. Their second sounds form a set that, when individually
		// appended to the preceding sounds, form foundation paths for child
		// subtrees. For each of those paths, a child subtree may exist if the path
		// is a valid partial word.
		last := sounds[len(sounds)-1]
		for pair := range this.PairSet {
			if pair[0] != last {
				continue
			}

			// Form the continued path and ensure that it's a valid partial word.
			path := make([]string, len(sounds), len(sounds)+1)
			copy(path, sounds)
			path = append(path, pair[1])
			if !this.validPart(path...) {
				continue
			}

			// If the path is actually a valid complete word, call the iterator.
			if this.checkPart(path...) {
				iterator(path...)
			}

			// Continue deeper.
			this.walk(iterator, path...)
		}
	}
}

/********************************** Statics **********************************/

/*--------------------------------- Public ----------------------------------*/

// Public version of examineMany().
func NewTraits(words []string) (traits *Traits, err error) {
	return examineMany(words)
}

/*--------------------------------- Private ---------------------------------*/

// Takes a word and returns a set of its characteristics, or an error if the
// word isn't valid, as per the validInput() function.
func examine(word string) (traits *Traits, err error) {
	// Make sure the word is valid.
	if !validInput(word) {
		return nil, errInvalid(word)
	}

	// Split into sounds.
	sounds, err := getSounds(word)
	if err != nil {
		return
	}

	// Mandate that at least two sounds are found.
	if len(sounds) < 2 {
		return nil, errType("less than two sounds found")
	}

	traits = new(Traits)

	// Add min and max number of consequtive sounds.
	traits.MinNSounds = len(sounds)
	traits.MaxNSounds = len(sounds)

	// Add max number of consequtive vowels.
	traits.MaxConseqVow = maxConsequtiveVowels(sounds)

	// Add max number of consequtive consonants.
	traits.MaxConseqCons = maxConsequtiveConsonants(sounds)

	// Add set of used sounds.
	traits.SoundSet = Set.New(nil, sounds...)

	// Find number of vowels vowels.
	nVow := countIntersections(sounds, knownVowels)

	// Add min and max total number of vowels.
	traits.MinNVowels = nVow
	traits.MaxNVowels = nVow

	// Find set of pairs of sounds.
	traits.PairSet = getPairs(sounds)

	/*
		// Disabled for now; this causes a combinatorial explosion so bad that test
		// duration goes from seconds to minutes, if not hours. We should add an
		// option to enable this for the `WordsN()` static function.

		// Add reverse pairs.
		addReversePairs(traits.PairSet)
	*/

	return
}

// Examines a slice of words and merges their traits, returning a Traits object
// that encompasses them all.
func examineMany(words []string) (traits *Traits, err error) {
	traits = &Traits{SoundSet: Set{}, PairSet: PairSet{}}

	// Examine each word and merge traits.
	for _, word := range words {
		tr, err := examine(word)
		if err != nil {
			return nil, err
		}

		// Merge MinNSounds and MaxNSounds.
		if n := tr.MinNSounds; traits.MinNSounds == 0 || n > 0 && n < traits.MinNSounds {
			traits.MinNSounds = n
		}
		if n := tr.MaxNSounds; n > traits.MaxNSounds {
			traits.MaxNSounds = n
		}

		// Merge MinNVowels and MaxNVowels.
		if n := tr.MinNVowels; traits.MinNVowels == 0 || n > 0 && n < traits.MinNVowels {
			traits.MinNVowels = n
		}
		if n := tr.MaxNVowels; n > traits.MaxNVowels {
			traits.MaxNVowels = n
		}

		// Merge MaxConseqVow.
		if tr.MaxConseqVow > traits.MaxConseqVow {
			traits.MaxConseqVow = tr.MaxConseqVow
		}

		// Merge MaxConseqCons.
		if tr.MaxConseqCons > traits.MaxConseqCons {
			traits.MaxConseqCons = tr.MaxConseqCons
		}

		// Merge SoundSet.
		for sound := range tr.SoundSet {
			traits.SoundSet.Add(sound)
		}

		// Merge PairSet.
		for pair := range tr.PairSet {
			traits.PairSet.Add(pair)
		}
	}

	return
}
