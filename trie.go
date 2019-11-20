package sm

// An implement of trie tree

import (
	"container/list"
)

type Node struct {
	IsKey    bool
	Children map[uint8]*Node
	Height   int
	Value    interface{}
}

type Trie struct {
	Root       *Node
	NumberNode uint
	NumberKey  uint
}

type Item struct {
	Key  []byte
	Node *Node
}

type Queue struct {
	List *list.List
}

func NewQueue() *Queue {
	queue := &Queue{}
	queue.List = list.New()
	return queue
}

func (q *Queue) Put(key []byte, node *Node) {
	item := &Item{Key: key, Node: node}
	q.List.PushBack(item)
}

func (q *Queue) Empty() bool {
	return q.List.Len() == 0
}

func (q *Queue) Get() (key []byte, node *Node) {
	if q.Empty() {
		return key, node
	}

	element := q.List.Front()
	q.List.Remove(element)
	value := element.Value
	item, _ := value.(*Item)
	return item.Key, item.Node
}

type Stack struct {
	List *list.List
}

func NewStack() (stack *Stack) {
	stack = &Stack{}
	stack.List = list.New()
	return stack
}

func (s *Stack) Empty() bool {
	return s.List.Len() == 0
}

func (s *Stack) Pop() (key []byte, node *Node) {
	if s.Empty() {
		return key, node
	}

	ele := s.List.Back()
	s.List.Remove(ele)
	value := ele.Value
	item, _ := value.(*Item)
	return item.Key, item.Node
}

func (s *Stack) Push(key []byte, node *Node) {
	s.List.PushBack(&Item{Key: key, Node: node})
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

func (trie *Trie) Walk(key []byte) (*Node, *Node, int) {
	var i int
	node := trie.Root
	parent := trie.Root

	for i = 0; i < len(key); i++ {
		order := key[i]
		parent = node
		node = node.Children[order]

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

	parent, node, step := trie.Walk(key)

	if step == keyLen {
		ret = 1
	}

	if node == nil {
		node = parent
	}

	for i = step; i < keyLen; i++ {
		order := key[i]
		node.Children[order] = CreateNode(false, i)
		trie.NumberNode += 1
		node = node.Children[order]
	}

	oldValue = node.Value
	node.IsKey = true
	node.Value = value
	trie.NumberKey += 1

	return oldValue, ret
}

func (trie *Trie) Find(key []byte) (ret bool, value interface{}) {
	ret = false
	value = nil
	keyLen := len(key)

	_, node, step := trie.Walk(key)

	if node == nil {
		return ret, value
	}

	if step == keyLen && node.IsKey {
		ret = true
		value = node.Value
	}
	return ret, value
}

func (trie *Trie) SeekAfter(key []byte) {

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
