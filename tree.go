package codex

// Data tree used inside State objects.

/*********************************** Type ************************************/

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

/********************************** Methods **********************************/

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

/********************************* Utilities *********************************/

// Creates shallow child nodes for a tree from the given pairs on the given
// path. Same algorithm as in Traits.walk(). See that method's comments.
func sprout(pairs PairSet, path ...string) (nodes map[string]*tree) {
	nodes = map[string]*tree{}
	if len(path) == 0 {
		// If there's no preceding path, use the first sounds from the pairs.
		for pair := range pairs {
			if _, ok := nodes[pair[0]]; ok {
				continue
			}
			nodes[pair[0]] = nil
		}
	} else {
		// Take pairs that begin with the last sound of the given path, collect
		// their second sounds, and register on the child node map.
		for pair := range pairs {
			if pair[0] == path[len(path)-1] {
				nodes[pair[1]] = nil
			}
		}
	}
	return
}
