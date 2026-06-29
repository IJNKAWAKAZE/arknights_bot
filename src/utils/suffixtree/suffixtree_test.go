package suffixtree

import (
	"testing"
)

func TestSuffixTree(t *testing.T) {
	words := []string{"banana", "apple", "中文app"}
	tree := NewGeneralizedSuffixTree()
	for k, word := range words {
		tree.Put(word, k)
	}
	indexes := tree.Search("a")

	assertIndexCounts(t, indexes, map[int]int{0: 3, 1: 1, 2: 1})

	indexes = tree.Search("文")

	if len(indexes) != 1 || indexes[0] != 2 {
		t.Error("indexes len should be 1 and indexes[0] must be 2,but ", len(indexes))
	}
}

func assertIndexCounts(t *testing.T, indexes []int, expected map[int]int) {
	t.Helper()
	counts := make(map[int]int, len(indexes))
	for _, index := range indexes {
		counts[index]++
	}
	if len(counts) != len(expected) {
		t.Fatalf("unexpected distinct indexes len: got %d want %d, values=%v", len(counts), len(expected), indexes)
	}
	for index, want := range expected {
		if got := counts[index]; got != want {
			t.Fatalf("count for index %d = %d, want %d (values=%v)", index, got, want, indexes)
		}
	}
}
