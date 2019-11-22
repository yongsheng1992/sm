package sm

import (
	"testing"
)

func TestTrie_Remove(t *testing.T) {
	trie := NewTrie()

	keyList := []string{"ABCD", "ABC", "AB", "A", "B", "C", "BCD"}

	for _, key := range keyList {
		trie.Insert([]byte(key), nil)
	}

	for range keyList {
		ok := trie.Remove([]byte("BC"))
		if ok {
			t.Error("remove failed")
		}
		if trie.NumberKey != uint(len(keyList)) {
			t.Error("remove failed")
		}
	}

	for idx, key := range keyList {
		trie.Remove([]byte(key))
		if trie.NumberKey != uint(len(keyList)-idx-1) {
			t.Error("remove test failed")
		}
	}
}
