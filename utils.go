package codex

// Utility functions and types.

import (
	"math/rand"
	"time"
)

/********************************* Utilities *********************************/

// Seed the random generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Takes a word and splits it into a series of known glyphs representing sounds.
func getSounds(word string, known Set) ([]string, error) {
	sounds := make([]string, 0, len(word))
	// Loop over the word, matching known glyphs. Break if no match is found.
	for i := 0; i < len(word); i++ {
		// Check for a known digraph.
		if i+2 <= len(word) && known.Has(word[i:i+2]) {
			sounds = append(sounds, word[i:i+2])
			i++
			// Check for a known monograph.
		} else if known.Has(word[i : i+1]) {
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

// Checks if the given word is too short or too long.
func validLength(word string) bool {
	return len(word) > 1 && len(word) < 33
}

// Copy of Join from the standard package `strings`.
func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}

// Republished rand.Perm.
func permutate(length int) []int {
	return rand.Perm(length)
}

// Shuffles a slice of strings in-place, using the Fisher–Yates method.
func shuffle(values []string) {
	for i := range values {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
}

// Gets the node values from the given map of child nodes.
func nodeValues(nodes map[string]*tree) (result []string) {
	if nodes == nil {
		return
	}
	if len(nodes) == 0 {
		return []string{}
	}
	result = make([]string, 0, len(nodes))
	for key := range nodes {
		result = append(result, key)
	}
	return
}

// Gets the node values from the given map of child nodes and shuffles it.
func randNodeValues(nodes map[string]*tree) (result []string) {
	result = nodeValues(nodes)
	if len(result) == 0 {
		return
	}
	shuffle(result)
	return
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

// Creates shallow child nodes for a tree from the given pairs on the given
// path.
func sprout(pairs PairSet, path ...string) (nodes map[string]*tree) {
	nodes = map[string]*tree{}
	if len(path) == 0 {
		// If no sound were passed, start from the root.
		for pair := range pairs {
			if _, ok := nodes[pair[0]]; ok {
				continue
			}
			nodes[pair[0]] = nil
		}
		// Otherwise continue from the given path.
	} else {
		// [ ... sounds ... ( last sound ] <- pair -> next sound )
		//
		// We investigate pairs that begin with the last sound of the given
		// preceding sounds. Their second sounds form a set that, when individually
		// appended to the preceding sounds, form foundation paths for child
		// subtrees. We register these second sounds on the child node map.
		for pair := range pairs {
			if pair[0] == path[len(path)-1] {
				nodes[pair[1]] = nil
			}
		}
	}
	return
}

/********************************** PairSet **********************************/

/**
 * Performance note: tried a slice version, and it significantly decreased the
 * package's benchmark performance. Sticking with a map for now.
 */

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

/********************************** errType **********************************/

type errType string

func (this errType) Error() string {
	return string(this)
}

/*********************************** tree ************************************/

// A tree that defines a set of string sequences. Node values represent sounds.
// A sequence of sounds obtained by visiting a branch represents a part of a
// word or a complete word (the distinction is defined by Traits). We define a
// tree as unordered, regardless of the implementation.
type tree struct {
	// The node's children, stored as a map where keys are children's values.
	nodes map[string]*tree
	// True if this node has been visited by an iterator.
	visited bool
}

// Finds or creates a node under the given path. Each value in the path
// represents a value of a descendant node.
func (this *tree) at(path ...string) (node *tree) {
	node = this
	for _, value := range path {
		if node.nodes[value] == nil {
			node.nodes[value] = new(tree)
		}
		node = node.nodes[value]
	}
	return
}
