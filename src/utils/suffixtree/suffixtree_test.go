package suffixtree

import (
	"fmt"
	"testing"
)

func TestSuffixTree(t *testing.T) {
	words := []string{"banana", "apple", "中文app"}
	tree := NewGeneralizedSuffixTree()
	for k, word := range words {
		tree.Put(word, k)
	}
	indexes := tree.Search("a", -1)

	if len(indexes) != 3 {
		t.Error("indexes len should be 3,but ", len(indexes))
	}
	fmt.Println(indexes)
	for _, index := range indexes {
		fmt.Println(words[index])
	}

	indexes = tree.Search("文", 0)

	if len(indexes) != 1 && indexes[0] != 2 {
		t.Error("indexes len should be 1 and indexes[0] must be 2,but ", len(indexes))
	}

	printnode("\t", tree.root)
}

func printnode(flag string, n *node) {
	for _, e := range n.edges {
		fmt.Printf("%s %s %v \n", flag, string(e.label), e.node.data)
		printnode(flag+"\t-", e.node)
	}
}
