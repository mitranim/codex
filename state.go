package codex

// Type that encapsulates word traits and maintains an internal state that is
// mutated by, and affects, its tree traversal methods.

/*********************************** Type ************************************/

// A state object encapsulates word traits and maintains an internal state that
// affects its tree traversal methods. The internal state, represented with a
// tree type, reflects the visited parts of the traits' virtual tree, keeping
// track of previously generated words. It allows us to speed up repeated
// traversals and guarantee no repeated words.
type state struct {
	// Word traits.
	traits *Traits

	// Tree that reflects the visited parts of the virtual tree defined by the
	// state's traits. It's built by state.walk() calls.
	tree *tree
}

/********************************** Methods **********************************/

// Walks the virtual tree of the state's traits, caching the visited parts in
// the state's inner tree. This caching lets us skip repeated Traits.validPart()
// checks, individual visited nodes, and fully visited subtrees. This
// significantly speeds up state.trip() traversals that restart from the root on
// each call, and lets us avoid revisiting nodes. This method also randomises
// the order of visiting subtrees from each node.
func (this *state) walk(iterator func(...string), sounds ...string) {
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
		node.nodes = sprout(this.traits.PairSet, sounds...)
	}

	// Loop over remaining child nodes and investigate their subtrees.
	for _, sound := range randNodeValues(node.nodes) {
		// Appending to sounds mutates their underlying array unless their cap was
		// <= 2 or so. If the iterator was expected to store sound slices, we would
		// allocate a new array for each path to avoid unexpected mutations. Right
		// now, we can get away with passing the slices as-is, because this method
		// is not exposed publicly and our own iterators don't store slices.
		path := append(sounds, sound)
		// Invalidate the path if it doesn't qualify as a partial word.
		if !this.traits.validPart(path...) {
			delete(node.nodes, sound)
			continue
		}
		// (1)(2) -> pre-order, (2)(1) -> post-order. Post-order is required by
		// state.walkRandom().
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

// Walks the state's virtual tree; for each path given to the wrapper function,
// we visit its subpaths in random order, marking the corresponding nodes as
// visited. For the distribution to be random, the tree needs to be traversed in
// post-order. We only visit paths that qualify as valid complete words and
// haven't been visited before.
func (this *state) walkRandom(iterator func(...string)) {
	iter := func(sounds ...string) {
		for _, index := range permutate(len(sounds)) {
			if index < 1 {
				continue
			}
			path := sounds[:index+1]
			node := this.tree.at(path...)
			if !node.visited {
				node.visited = true
				if this.traits.checkPart(path...) {
					iterator(path...)
				}
			}
		}
	}
	this.walk(iter)
}

// Uses state.walkRandom() to traverse the tree and interrupts the walking after
// one successful call to the given iterator function.
func (this *state) trip(iterator func(...string)) {
	defer aid()
	this.walkRandom(func(sounds ...string) {
		iterator(sounds...)
		interrupt()
	})
}
