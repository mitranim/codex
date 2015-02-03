package codex

// Data tree used inside State objects.

/*********************************** Type ************************************/

// A tree that defines a set of string sequences. Node values represent sounds.
// A sequence of sounds obtained by visiting a branch represents a part of a
// word or a complete word (the distinction is defined by Traits). We define a
// tree as unordered, regardless of the implementation.
type tree struct {
	// The node's own value.
	value string

	// The node's children, stored as a map where keys are children's values. This
	// contradicts the typical tree definition where nodes only have references to
	// their children, not their values. This also means there's no single source
	// of truth for child values. We're doing this to improve child lookup
	// performance in State tree traversal methods.
	nodes map[string]*tree

	// True if this node's path forms a valid partial word.
	part tern

	// True if this node's path forms a valid complete word.
	complete tern
}

/**
 * Performance and implementation notes.
 *
 * The original type was `type tree map[string]tree`, where map keys
 * represented node values. It was the most minimal definition I could think of.
 * Back then, a tree's lifecycle consisted of the following steps:
 *   1) build the entire tree;
 *   2) traverse the entire tree, calling an iterator function on each path.
 * For this lifecycle, the map structure turned out to be suboptimal. Maps are
 * better for random access, and slices are better for sequential looping.
 * Switching to `type tree struct{value string, nodes []*tree}` significantly
 * improved the traversal performance.
 *
 * In the current implementation, we never want to build and traverse a tree at
 * once. When generating an entire word set, we traverse a virtual tree
 * associated with a traits object. A real tree is used as a reference table
 * inside a State object, built on the fly and used for random lookups while
 * traversing a virtual tree in State methods. That's why we've switched back
 * to using maps for child nodes.
 */

/********************************** Methods **********************************/

// Walks the tree in depth-first pre-order, calling the given iterator function
// on each node, passing the ordered sequence of values from the root to the
// current node, inclusively. Ignores empty values.
func (this tree) walk(iterator func(...string), trail ...string) {
	if iterator == nil {
		return
	}

	// Because the trail (combined with the current node's value, if any) will be
	// passed recursively to multiple child nodes, where it will be subject to
	// `append` calls, we must ensure that its underlying array (its cap) equals
	// its length. If we don't, appending the current value will usually grow its
	// cap beyond its length, which will cause appended slices in child nodes to
	// share the same underlying array and mutate it for each other when appending
	// new values. If the iterator function stores slices of its arguments, they
	// will be mutated in unexpected ways by this. Therefore we must allocate a
	// new array to store each path, with the array's length matching the path's.
	var path []string
	if this.value == "" {
		path = make([]string, len(trail))
		copy(path, trail)
	} else {
		path = make([]string, len(trail), len(trail)+1)
		copy(path, trail)
		path = append(path, this.value)
		iterator(path...)
	}

	// Pass the newly allocated path to each child node. Because its length equals
	// its cap, when child nodes append their values to it, this is guaranteed to
	// have no effect on paths passed to other siblings.
	for _, node := range this.nodes {
		node.walk(iterator, path...)
	}
}

// Finds or creates a node under the given path. Each value in the path
// represents a value of a descendant node. When a node is created, it's given
// the current value.
func (this *tree) at(path ...string) (node *tree) {
	node = this
	for _, value := range path {
		if _, ok := node.nodes[value]; !ok {
			if node.nodes == nil {
				node.nodes = map[string]*tree{}
			}
			node.nodes[value] = &tree{value: value}
		}
		node = node.nodes[value]
	}
	return
}
