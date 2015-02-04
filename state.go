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
	this.walkStraight(func(sounds ...string) {
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

// Walks the virtual tree of the state's traits, caching the visited parts in
// the state's inner tree. This caching lets us skip repeated Traits.validPart()
// checks, individual visited nodes, and fully visited subtrees. This has no
// benefit for a one-shot traversal that visits the entire tree at once (see
// State.Words()), but significantly speeds up traversals that restart from the
// root after exiting early via panic (see State.trip()), and lets us avoid
// revisiting nodes.
//
// This also randomises the order of visiting subtrees from each node.
func (this *State) walk(iterator func(...string), sounds ...string) {
	if iterator == nil {
		return
	}
	if this.tree == nil {
		this.tree = new(tree)
	}

	// Find or create a matching node for this path. If it doesn't have child
	// nodes yet, make a shallow map to track valid paths.
	node := this.tree.at(sounds...)
	if node.nodes == nil {
		node.nodes = sprout(this.Traits.PairSet, sounds...)
	}

	// Loop over remaining child nodes and investigate their subtrees.
	for _, sound := range randNodeValues(node.nodes) {
		path := append(sounds, sound)
		// Invalidate the path if it doesn't qualify as a partial word.
		if !this.Traits.validPart(path...) {
			delete(node.nodes, sound)
			continue
		}
		// (1)(2) -> pre-order, (2)(1) -> post-order. Post-order is required by
		// State.walkRandom(); it slows down State.Words() by about 10-15%, which
		// doesn't warrant its own separate algorithm.
		// (2) Continue recursively.
		this.walk(iterator, path...)
		// (1) If this path hasn't yet been visited, feed it to the iterator.
		if !node.at(sound).visited {
			iterator(path...)
		}
		// If this code is reached, the subtree is used up, so we forget about it.
		delete(node.nodes, sound)
	}
}

// Walks the state's virtual tree, visiting paths that qualify as valid complete
// words.
func (this *State) walkStraight(iterator func(...string)) {
	iter := func(sounds ...string) {
		this.tree.at(sounds...).visited = true
		if this.Traits.checkPart(sounds...) {
			iterator(sounds...)
		}
	}
	this.walk(iter)
}

// Walks the state's virtual tree; for each paths given to the wrapper function,
// we visit its subpaths in random order, marking the corresponding nodes as
// visited. For the distribution to be random, the tree needs to be traversed in
// post-order. We only visit paths that qualify as valid complete words and
// haven't been visited before.
func (this *State) walkRandom(iterator func(...string)) {
	iter := func(sounds ...string) {
		for _, index := range permutate(len(sounds)) {
			if index < 1 {
				continue
			}
			path := sounds[:index+1]
			node := this.tree.at(path...)
			if !node.visited {
				node.visited = true
				if this.Traits.checkPart(path...) {
					iterator(path...)
				}
			}
		}
	}
	this.walk(iter)
}

// Uses State.walkRandom() to traverse the tree and interrupts the walking after
// one successful call to the given iterator function.
func (this *State) trip(iterator func(...string)) {
	defer aid()
	this.walkRandom(func(sounds ...string) {
		iterator(sounds...)
		interrupt()
	})
}
