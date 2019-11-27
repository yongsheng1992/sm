package sm

import (
	"testing"
)

func TestTrie_Find(t *testing.T) {
	trie := NewTrie()
	keyList := []string{"ABCD", "ABC", "AB", "A", "B", "C", "BCD"}

	for _, key := range keyList {
		trie.Insert([]byte(key), key)
	}

	for _, key := range keyList {
		ret, value := trie.Find([]byte(key))
		if ret == false {
			t.Error("find or insert error")
		}
		if value != key {
			t.Error("key value is error")
		}
	}
}

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
		if trie.NumberKey != int32(len(keyList)) {
			t.Error("remove failed")
		}
	}

	for idx, key := range keyList {
		trie.Remove([]byte(key))
		if trie.NumberKey != int32(len(keyList)-idx-1) {
			t.Error("remove test failed")
		}
	}
}
