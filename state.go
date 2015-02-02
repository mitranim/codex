package codex

// Type that encapsulates word traits and maintains an internal state that
// persists between calls to its word generation methods.

import (
	"strings"
)

/*********************************** Type ************************************/

// A State object encapsulates word traits and maintains an internal state that
// affects its word generator methods. It defines its own methods for traversing
// the virtual tree defined by the traits. The internal state, represented with
// a tree type, reflects the visited parts of the virtual tree, keeping track of
// previously generated words. It allows us to speed up repeated traversals and
// guarantee no repeated words.
//
// A state must always be created with a NewState() call, or given an existing
// Traits object obtained with NewTraits(). Its behaviour without an associated
// Traits object is undefined.
type State struct {
	// Word traits.
	Traits *Traits

	// Tree that reflects the visited parts of the virtual tree defined by the
	// state's traits. It's built by State.walk() calls.
	tree *tree
}

/********************************** Statics **********************************/

// Takes a sample group of words, analyses their traits, and builds a State.
func NewState(words []string) (*State, error) {
	traits, err := examineMany(words)
	if err != nil {
		return nil, err
	}
	state := State{Traits: traits}
	return &state, nil
}

/********************************** Methods **********************************/

/*--------------------------------- Public ----------------------------------*/

// Generates the entire set of words defined by the state. The virtual pool of
// words is shared with State.WordsN(). This method exhausts the pool
// completely; subsequent calls to State.Words() and State.WordsN() return an
// empty result.
func (this *State) Words() (words Set) {
	iterator := func(sounds ...string) {
		words.Add(strings.Join(sounds, ""))
	}
	this.walk(iterator)
	return
}

// Generates a subset of the set of words defined by the state, limited to the
// given count. The words are guaranteed to never repeat. The virtual pool of
// words is shared with State.Words(). Subsequent calls to State.Words() are
// guaranteed to not include any of the words that have been returned by
// State.WordsN(). If this is called enough times to exhaust the entire pool,
// subsequent calls to State.Words() and State.WordsN() return an empty result.
func (this *State) WordsN(num int) (words Set) {
	iterator := func(sounds ...string) {
		words.Add(strings.Join(sounds, ""))
	}
	for i := 0; i < num; i++ {
		this.trip(iterator)
	}
	return
}

/*--------------------------------- Private ---------------------------------*/

// Walks the virtual tree of valid partial words associated with the state's
// traits, calling the given iterator function on paths that qualify as valid
// complete words. The principle is the same as Traits.walk(), with the
// difference that visited parts of the virtual tree are cached inside the
// state. This lets us avoid recalculating validity of partial words and calling
// the iterator on previously used paths.
//
// The state's internal tree plays three roles here:
//   1) it invalidates virtual paths that don't qualify for partial words,
//      marking them on their parent nodes; this lets us avoid calling
//      Traits.validPart() for the same paths later; it also lets us avoid
//      repeating Traits.validPart() checks for paths that have already been
//      validated;
//   2) it marks paths that have been passed to the iterator function; this lets
//      us guarantee that the state never returns the same word twice;
//   3) it invalidates subtrees that don't have any unused paths left; this lets
//      us avoid revisiting those subtrees on subsequent calls, speeding up
//      the State.WordsN() method approximately up to ten times over the course
//      of many repeated calls.
//
// This method is compatible with early exits via panic; see State.trip().
func (this *State) walk(iterator func(...string), sounds ...string) {
	if iterator == nil {
		return
	}

	if this.tree == nil {
		this.tree = new(tree)
	}

	// If no sounds are passed, start from the root.
	if len(sounds) == 0 {
		// The values of the first-level nodes are the first values of the pair set
		// associated with the traits.
		for _, first := range randFirsts(this.Traits.PairSet) {
			// Check for blocked paths.
			if this.tree.blocked.Has(first) {
				continue
			}
			// Continue recursively.
			this.walk(iterator, first)
			// If this code is reached, the child subtree is guaranteed to have been
			// used up, so we mark it as blocked.
			this.tree.blocked.Add(first)
		}

		// If sounds are passed, continue from that path onward.
	} else {
		// Find or create a tree node under this path.
		node := this.tree.at(sounds...)

		// [ ... sounds ... ( last sound ] <- pair -> next sound )
		//
		// We investigate pairs that begin with the last sound of the given
		// preceding sounds. Their second sounds form a set that, when individually
		// appended to the preceding sounds, form foundation paths for child
		// subtrees. For each of those paths, a child subtree may exist if the path
		// is a valid partial word.
		for _, second := range randSeconds(this.Traits.PairSet, sounds[len(sounds)-1]) {
			// Check for blocked paths.
			if node.blocked.Has(second) {
				continue
			}

			// Form the continued path.
			path := make([]string, len(sounds), len(sounds)+1)
			copy(path, sounds)
			path = append(path, second)

			// If the node doesn't have a child under this path, verify that the path
			// is a valid partial word. If it's not, block the sound and skip this
			// subtree. If it's valid, register a child under this sound. If the node
			// already has a child under this path, skip the checks.
			if node.nodes[second] == nil {
				if !this.Traits.validPart(path...) {
					node.blocked.Add(second)
					continue
				}
			}
			child := this.tree.at(path...)

			// Continue deeper.
			this.walk(iterator, path...)

			// If we have reached a leaf (a tip of a branch), call the iterator in
			// random order on each subpath of the leaf's path that qualifies as a
			// valid word. Before calling the iterator on a subpath, we find the node
			// corresponding to that subpath, and mark it as used.
			//
			// Performance note: deferring iterator calls until reaching a leaf, then
			// randomising the subpaths, slows down State.trip() by about two times.
			// We consider this an acceptable cost because State.WordsN() is still
			// fast enough for our purposes, and State.Words() is almost unaffected.
			if len(child.nodes) == 0 {
				for _, index := range permutate(len(path) + 1) {
					if index < 2 {
						continue
					}
					subpath := path[:index]
					subnode := this.tree.at(subpath...)
					if !subnode.used && this.Traits.checkPart(subpath...) {
						subnode.used = true
						iterator(subpath...)
					}
				}
			}

			// If this code is reached, the child subtree is guaranteed to have been
			// used up, so we mark it as blocked.
			node.blocked.Add(second)
		}
	}
}

// Uses State.walk() to traverse the tree and interrupts the walking after one
// successful call to the given iterator function.
func (this *State) trip(iterator func(...string)) {
	defer aid()
	this.walk(func(sounds ...string) {
		iterator(sounds...)
		interrupt()
	})
}
