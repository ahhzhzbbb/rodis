package engine

import "sync"

const DefaultCapNode = 10000

type QuickList struct {
	head        *QLNode
	tail        *QLNode
	len         int
	bytesOfNode int
	mu          sync.RWMutex
}

type QLNode struct {
	zip  *ZipList
	prev *QLNode
	next *QLNode
}

func (n *QLNode) isFull(bytesOfNode int) bool {
	return n.zip.GetBytes() >= uint32(bytesOfNode)
}

func (n *QLNode) isEmpty() bool {
	return n.zip.Length() == 0
}

func NewQuickList(bytesOfNode int, args []string) *QuickList {
	if len(args) == 0 {
		return nil
	}

	if bytesOfNode <= 0 {
		bytesOfNode = DefaultCapNode
	}
	var res QuickList

	qlNode := &QLNode{
		zip: NewZipList(),
	}

	res.head = qlNode

	index := 0
	for index < len(args) {
		if !qlNode.isFull(bytesOfNode) {
			qlNode.zip.PushBack(args[index])
			index++
		} else {
			newNode := &QLNode{
				zip: NewZipList(),
			}
			qlNode.next = newNode
			newNode.prev = qlNode
			qlNode = newNode
		}
	}

	res.tail = qlNode
	res.len = len(args)
	res.bytesOfNode = bytesOfNode

	return &res
}

// func (ql *QuickList) LinkToQuickList(newList *QuickList) (headList *QLNode) {
// 	ql.tail.next = newList.head
// 	ql.tail = newList.tail
// 	headList = ql.head
// 	ql.len += newList.len
// 	return headList
// }

func (ql *QuickList) GetElements() (elements []string) {
	curr := ql.head
	for {
		if curr == nil {
			break
		}
		elements = append(elements, curr.zip.GetElements()...)
		curr = curr.next
	}
	return elements
}

func (ql *QuickList) PushBack(elements []string) {
	for _, element := range elements {
		if !ql.tail.isFull(ql.bytesOfNode) {
			ql.tail.zip.PushBack(element)
		} else {
			newNode := &QLNode{
				zip: NewZipList(),
			}
			newNode.zip.PushBack(element)
			ql.tail.next = newNode
			newNode.prev = ql.tail
			ql.tail = newNode
		}
	}
	ql.len += len(elements)
}

func (ql *QuickList) PushFront(elements []string) {
	for i := len(elements) - 1; i >= 0; i-- {
		element := elements[i]
		if !ql.head.isFull(ql.bytesOfNode) {
			ql.head.zip.PushFront(element)
		} else {
			newNode := &QLNode{
				zip: NewZipList(),
			}
			newNode.zip.PushFront(element)
			newNode.next = ql.head
			ql.head.prev = newNode
			ql.head = newNode
		}
	}
	ql.len += len(elements)
}

func (ql *QuickList) PopBack() string {
	if ql.len == 0 {
		return ""
	}
	value := ql.tail.zip.PopBack()
	if ql.tail.zip.Length() == 0 && ql.len > 1 {
		ql.tail = ql.tail.prev
		ql.tail.next = nil
	}
	ql.len--
	return value
}

func (ql *QuickList) PopFront() string {
	if ql.len == 0 {
		return ""
	}
	value := ql.head.zip.PopFront()
	if ql.head.zip.Length() == 0 && ql.len > 1 {
		ql.head = ql.head.next
		ql.head.prev = nil
	}
	ql.len--
	return value
}

func (ql *QuickList) GetIndexOFElement(element string) (*QLNode, int, bool) {
	curr := ql.head
	for curr != nil {
		if pos, found := curr.zip.GetIndexOfElement(element); found {
			return curr, pos, true
		} else {
			curr = curr.next
		}
	}
	return nil, -1, false
}

func (ql *QuickList) Insert(node *QLNode, indexInNode int, value string) bool {
	if indexInNode < 0 || node == nil {
		return false
	}

	if node == ql.head && indexInNode == -1 {
		ql.PushFront([]string{value})
		return true
	}

	if node == ql.tail && indexInNode == int(node.zip.Length()) {
		ql.PushBack([]string{value})
		return true
	}

	elementSize := 2 + uint32(len(value))
	if node.zip.GetBytes()+elementSize > uint32(ql.bytesOfNode) {
		newNode := &QLNode{
			zip: node.zip.SplitList(indexInNode),
		}
		if newNode == nil {
			return false
		}
		node.zip.PushBack(value)
	}
}

func (ql *QuickList) Length() int {
	return ql.len
}
