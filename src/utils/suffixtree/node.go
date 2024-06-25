package suffixtree

import (
	"sort"
)

type node struct {
	/*
	 * The payload array used to store the data (indexes) associated with this node.
	 * In this case, it is used to store all property indexes.
	 */
	data []int
	/**
	 * The set of edges starting from this node
	 */
	edges []*edge
	/**
	 * The suffix link as described in Ukkonen's paper.
	 * if str is the string denoted by the path from the root to this, this.suffix
	 * is the node denoted by the path that corresponds to str without the first rune.
	 */
	suffix *node
}

/*
 * getData returns the first numElements elements from the ones associated to this node.
 *
 * Gets data from the payload of both this node and its children, the string representation
 * of the path to this node is a substring of the one of the children nodes.
 *
 * @param numElements the number of results to return. Use <=0 to get all
 * @return the first numElements associated to this node and children
 */
func (n *node) getData(numElements int) (ret []int) {

	if numElements > 0 {
		if numElements > len(n.data) {
			numElements -= len(n.data)
			ret = n.data
		} else {
			ret = n.data[:numElements]
			return
		}
	} else {
		ret = n.data
	}

	// need to get more matches from child nodes. This is what may waste time
	for _, edge := range n.edges {
		data := edge.node.getData(numElements)
	NEXTIDX:
		for _, idx := range data {
			for _, v := range ret {
				if v == idx {
					continue NEXTIDX
				}
			}

			if numElements > 0 {
				numElements--
			}
			ret = append(ret, idx)
		}

		if numElements == 0 {
			break
		}
	}

	return
}

// addRef adds the given index to the set of indexes associated with this
func (n *node) addRef(index int) {
	if n.contains(index) {
		return
	}
	n.addIndex(index)
	// add this reference to all the suffixes as well
	iter := n.suffix
	for iter != nil {
		if iter.contains(index) {
			break
		}
		iter.addRef(index)
		iter = iter.suffix
	}
}

func (n *node) contains(index int) bool {
	i := sort.SearchInts(n.data, index)
	return i < len(n.data) && n.data[i] == index
}

func (n *node) addEdge(r rune, e *edge) {
	if idx := n.search(r); idx == -1 {
		n.edges = append(n.edges, e)
		sort.Slice(n.edges, func(i, j int) bool { return n.edges[i].label[0] < n.edges[j].label[0] })
	} else {
		n.edges[idx] = e
	}

}

func (n *node) getEdge(r rune) *edge {
	idx := n.search(r)
	if idx < 0 {
		return nil
	}
	return n.edges[idx]
}

func (n *node) search(r rune) int {
	idx := sort.Search(len(n.edges), func(i int) bool { return n.edges[i].label[0] >= r })
	if idx < len(n.edges) && n.edges[idx].label[0] == r {
		return idx
	}

	return -1
}

func (n *node) addIndex(idx int) {
	n.data = append(n.data, idx)
}

func newNode() *node {
	return &node{}
}
