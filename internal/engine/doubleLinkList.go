package engine

type Node struct {
	prev  *Node
	next  *Node
	value string
}

type DoubleLinkList struct {
	head *Node
	tail *Node
	len  int
}

func NewNode(value string) *Node {
	return &Node{
		value: value,
	}
}

func NewDoubleLinkList() *DoubleLinkList {
	return &DoubleLinkList{}
}

func (l *DoubleLinkList) PushBack(node *Node) {
	if l.tail == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}
	l.len++
}

func (l *DoubleLinkList) PushFront(node *Node) {
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		l.head.prev = node
		node.next = l.head
		l.head = node
	}
	l.len++
}

func (l *DoubleLinkList) PopFront() *Node {
	if l.head == nil {
		return nil
	}
	node := l.head
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}
	l.len--
	return node
}

func (l *DoubleLinkList) PopBack() *Node {
	if l.tail == nil {
		return nil
	}
	node := l.tail
	l.tail = l.tail.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	l.len--
	return node
}

func (l *DoubleLinkList) GetNodeByValue(value string) *Node {
	current := l.head
	for current != nil {
		if current.value == value {
			return current
		}
		current = current.next
	}
	return nil
}

func (l *DoubleLinkList) InsertAfter(node *Node, value string) {
	newNode := NewNode(value)
	newNode.prev = node
	newNode.next = node.next
	if node.next != nil {
		node.next.prev = newNode
	} else {
		l.tail = newNode
	}
	node.next = newNode
	l.len++
}

func (l *DoubleLinkList) Len() int {
	return l.len
}
