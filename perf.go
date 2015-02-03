package codex

// Utilities optimised with benchmarks. Keeping this in a separate file to keep
// track of what has and hasn't been optimised.

import (
	"strings"
)

// Returns the biggest number of consequtive vowels that occurs in the given
// sound sequence.
func maxConsequtiveVowels(sounds []string, vowels Set) (max int) {
	var count int
	for _, sound := range sounds {
		if !vowels.Has(sound) {
			count = 0
		} else {
			count++
			if count > max {
				max = count
			}
		}
	}
	return
}

// Returns the biggest number of consequtive consonants that occurs in the given
// sound sequence.
func maxConsequtiveConsonants(sounds []string, vowels Set) (max int) {
	var count int
	for _, sound := range sounds {
		if vowels.Has(sound) {
			count = 0
		} else {
			count++
			if count > max {
				max = count
			}
		}
	}
	return
}

// Counts how many strings in the given slice occur in the given set.
func countIntersections(strings []string, set Set) (count int) {
	for _, value := range strings {
		if set.Has(value) {
			count++
		}
	}
	return
}

// Counts the occurrences of the given pair of strings in the given slice. This
// is used in Traits.validPairs as a performance optimisation.
func countPair(strings []string, prev, current string) (count int) {
	var ownPrev string
	for index, ownCurrent := range strings {
		if index == 0 {
			ownPrev = ownCurrent
			continue
		}
		if ownPrev == prev && ownCurrent == current {
			count++
		}
		ownPrev = ownCurrent
	}
	return
}

/************************************ Set ************************************/

// Set behaves like a set of strings. Tried a map version and a slice version.
// The slice version was marginally faster for very small datasets and with
// little lookup. The map version is significantly faster for anything more than
// a handful of values, or with many lookups. The difference is huge for big
// datasets, which this package has aplenty.
type Set map[string]struct{}

// Creates a new set from the given keys. Usage:
//   Set.New(nil, "one", "other")
func (Set) New(keys ...string) Set {
	set := make(Set, len(keys))
	for _, key := range keys {
		set.Add(key)
	}
	return set
}

// Adds the given element.
func (this *Set) Add(key string) {
	if *this == nil {
		*this = Set{}
	}
	(*this)[key] = struct{}{}
}

// Deletes the given element.
func (this *Set) Del(key string) {
	delete((*this), key)
}

// Checks for the presence of the given element.
func (this *Set) Has(key string) bool {
	_, ok := (*this)[key]
	return ok
}

// Prints itself nicely in fmt(%#v).
func (this Set) GoString() string {
	keys := make([]string, 0, len(this))
	for key := range this {
		keys = append(keys, `"`+key+`"`)
	}
	return "{" + strings.Join(keys, ", ") + "}"
}

// Prints itself nicely in println().
func (this Set) String() string {
	return this.GoString()
}
