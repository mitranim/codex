package codex

// Type that encapsulates word traits and maintains an internal state that
// persists between calls to its word generation methods.

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
	traits, err := NewTraits(words)
	if err != nil {
		return nil, err
	}
	state := &State{Traits: traits}
	return state, nil
}

/********************************** Methods **********************************/

/*--------------------------------- Public ----------------------------------*/

// Generates the entire set of words defined by the state. The virtual pool of
// words is shared with State.WordsN(). This method exhausts the pool
// completely; subsequent calls to State.Words() and State.WordsN() return an
// empty result.
func (this *State) Words() (words Set) {
	this.walk(func(sounds ...string) {
		words.Add(join(sounds, ""))
	})
	return
}

// Generates a randomly distributed subset of the set of words defined by the
// state, limited to the given count. The words are guaranteed to never repeat.
// The virtual pool of words is shared with State.Words(). Subsequent calls to
// State.Words() are guaranteed to not include any of the words that have been
// returned by State.WordsN(). If this is called enough times to exhaust the
// entire pool, subsequent calls to State.Words() and State.WordsN() return an
// empty result.
func (this *State) WordsN(num int) (words Set) {
	iterator := func(sounds ...string) {
		words.Add(join(sounds, ""))
	}
	for i := 0; i < num; i++ {
		this.trip(iterator)
	}
	return
}

/*--------------------------------- Private ---------------------------------*/

// Walks the virtual tree of the state's traits, caching the visited parts as a
// real tree structure. The inner tree caches the results of Traits.validPart()
// and Traits.validComplete() checks, speeding up repeated traversals and
// letting us skip paths that have already been fed to an iterator function in a
// preceding traversal.
//
// Due to caching, this method is compatible with early exits via panic; see
// State.trip().
func (this *State) walk(iterator func(...string)) {
	if this.tree == nil {
		this.tree = new(tree)
	}
	this.Traits.walk(func(sounds ...string) bool {
		node := this.tree.at(sounds...)

		// Check partial validity. If this path doesn't qualify as a partial word,
		// abort the branch.
		if node.part == nil {
			node.part = ternary(this.Traits.validPart(sounds...))
		}
		if !*node.part {
			return false
		}

		// Check complete validity. If this path qualifies as a complete word, call
		// the iterator with it. If the value is already set and true, this implies
		// that an iterator has already been called with this path. In this case, we
		// don't call it again.
		if node.complete == nil {
			node.complete = ternary(this.Traits.checkPart(sounds...))
			if *node.complete {
				iterator(sounds...)
			}
		}
		return true
	})
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
