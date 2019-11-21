package sm

import "container/list"

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
