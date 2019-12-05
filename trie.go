package sm

import (
	"sync"
	"sync/atomic"
)

// An implement of trie tree

type Node struct {
	IsKey    bool
	Children map[uint8]*Node
	Height   int
	Value    interface{}
	Lock     sync.Mutex
}

type Trie struct {
	Root       *Node
	NumberNode int32
	NumberKey  int32
}

func NewTrie() *Trie {
	root := &Node{IsKey: false, Children: make(map[uint8]*Node), Height: 0}
	trie := &Trie{Root: root, NumberNode: 0, NumberKey: 0}
	return trie
}

func CreateNode(isKey bool, height int) *Node {
	node := &Node{IsKey: isKey, Height: height, Children: make(map[uint8]*Node)}
	return node
}

func (node *Node) InsertChild(ord uint8, child *Node) {
	node.Children[ord] = child
}

func (node *Node) RemoveChild(ord uint8) {
	delete(node.Children, ord)
}

func (node *Node) GetChild(ord uint8) *Node {
	return node.Children[ord]
}

func (node *Node) Update(isKey bool, value interface{}) {
	node.Lock.Lock()
	defer node.Lock.Unlock()

	node.IsKey = isKey
	node.Value = value
}

func (trie *Trie) increaseNumberNode() {
	atomic.AddInt32(&trie.NumberNode, 1)
}

func (trie *Trie) decreaseNumberNode() {
	atomic.AddInt32(&trie.NumberNode, -1)
}

func (trie *Trie) increaseNumberKey() {
	atomic.AddInt32(&trie.NumberKey, 1)
}

func (trie *Trie) decreaseNumberKey() {
	atomic.AddInt32(&trie.NumberKey, -1)
}

func (trie *Trie) Walk(key []byte) (*Node, *Node, int) {
	var i int
	node := trie.Root
	parent := trie.Root

	for i = 0; i < len(key); i++ {
		order := key[i]
		parent = node
		node = node.GetChild(order)

		if node == nil {
			break
		}

	}

	return parent, node, i
}

func (trie *Trie) Insert(key []byte, value interface{}) (oldValue interface{}, ret int) {
	var i int
	var parent *Node
	var node *Node

	keyLen := len(key)
	ret = 0
	parent = trie.Root
	node = trie.Root

	// 此时walk是没有查找到
	// 但是在插入的时候，已经有其它协程插入了
	// 这样数据就丢失了，特别是节点是key的情况
	// 此时就丢掉了一个key
	for i = 0; i < len(key); i++ {
		order := key[i]
		// 当前节点添加读锁
		parent = node
		parent.Lock.Lock()

		node = node.GetChild(order)

		// 这里break之后，parent会持有读锁
		if node == nil {
			break
		}

		// 当前节点释放读锁
		parent.Lock.Unlock()
	}

	if i == keyLen && node != nil {
		ret = 1
		oldValue = node.Value
		if !node.IsKey {
			node.Update(true, value)
			trie.increaseNumberKey()

		}
		return oldValue, ret
	}

	if node == nil {
		node = parent
	}

	for ; i < keyLen; i++ {
		order := key[i]

		childNode := CreateNode(false, i)

		node.Children[order] = childNode
		node.Lock.Unlock()
		trie.increaseNumberNode()

		node = childNode
		node.Lock.Lock()
	}

	oldValue = node.Value
	node.IsKey = true
	node.Value = value
	node.Lock.Unlock()
	trie.increaseNumberKey()

	return oldValue, ret
}

func (trie *Trie) Find(key []byte) (ret bool, value interface{}) {
	ret = false
	value = nil
	keyLen := len(key)

	_, node, step := trie.Walk(key)

	if step == keyLen && node.IsKey {
		ret = true
		value = node.Value
	}

	return ret, value
}

func (trie *Trie) SeekAfter(key []byte) (it *Iterator) {
	_, node, _ := trie.Walk(key)

	if node == nil {
		return it
	}

	it = NewIterator(key, node)
	return it
}

func (trie *Trie) Remove(key []byte) bool {
	var i int
	var parent *Node
	var node *Node

	parent = trie.Root
	node = trie.Root
	keyLen := len(key)

	for i = 0; i < keyLen; i++ {
		order := key[i]
		// 当前节点添加读锁
		parent = node
		parent.Lock.Lock()

		node = node.GetChild(order)

		// 这里break之后，parent会持有读锁
		if node == nil {
			break
		}

		// 当前节点释放读锁
		parent.Lock.Unlock()
	}

	if i < keyLen || i == 0 {
		return false
	}

	if i == keyLen {
		if node != nil {
			// 不是key直接返回
			if !node.IsKey {
				return false
			}

			trie.decreaseNumberKey()
			// 是key但是有子节点
			if len(node.Children) > 0 {
				node.Update(false, nil)
				return true
			}

			// 是key没有子节点删除该节点
			parent.RemoveChild(key[i-1])
			trie.decreaseNumberNode()
			return true
		}
	}
	return false
}

func (trie *Trie) SeekBefore(key []byte) []int {
	var i int
	var flags []int
	node := trie.Root

	for i = 0; i < len(key); i++ {
		order := key[i]

		node = node.Children[order]

		if node == nil {
			break
		} else if node.IsKey {
			flags = append(flags, i)
		}

	}

	return flags
}

func (trie *Trie) BFS(fn func(key []byte, node *Node)) {
	queue := NewQueue()
	queue.Put(make([]byte, 0), trie.Root)

	for !queue.Empty() {
		suffix, node := queue.Get()
		fn(suffix, node)

		for ord, child := range node.Children {
			path := make([]byte, len(suffix))
			copy(path, suffix)
			path = append(path, ord)
			queue.Put(path, child)
		}
	}
}
