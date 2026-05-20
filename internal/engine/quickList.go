package engine

import "sync"

const DefaultCapNode = 4

type QuickList struct {
	head    *QLNode
	tail    *QLNode
	len     int
	capNode int
	mu      sync.RWMutex
}

type QLNode struct {
	zip  *ZipList
	prev *QLNode
	next *QLNode
}

func NewQuickList(capNode int, args []string) *QuickList {
	if capNode <= 0 {
		capNode = DefaultCapNode
	}
	var res QuickList

	qlNode := &QLNode{
		zip: NewZipList(),
	}

	res.head = qlNode

	index := 0
	for index < len(args) {
		remain := capNode - int(qlNode.zip.Length())
		for i := 0; i < remain && index < len(args); i++ {
			qlNode.zip.PushBack(args[index])
			index++
		}
		if index < len(args) {
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
	res.capNode = capNode

	return &res
}

func (ql *QuickList) LinkToQuickList(newList *QuickList) (headList *QLNode) {
	ql.tail.next = newList.head
	ql.tail = newList.tail
	headList = ql.head
	ql.len += newList.len
	return headList
}

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
	if ql.tail.zip.Length()+uint16(len(elements)) <= uint16(ql.capNode) {
		for _, element := range elements {
			ql.tail.zip.PushBack(element)
		}
	} else {
		index := 0
		remain := ql.capNode - int(ql.tail.zip.Length())
		for index < remain && index < len(elements) {
			ql.tail.zip.PushBack(elements[index])
			index++
		}
		for index < len(elements) {
			newNode := &QLNode{
				zip: NewZipList(),
			}
			for i := 0; i < ql.capNode && index < len(elements); i++ {
				newNode.zip.PushBack(elements[index])
				index++
			}
			ql.tail.next = newNode
			newNode.prev = ql.tail
			ql.tail = newNode
		}
	}
	ql.len += len(elements)
}

func (ql *QuickList) PushFront(elements []string) {
	if ql.head.zip.Length()+uint16(len(elements)) <= uint16(ql.capNode) {
		for _, element := range elements {
			ql.head.zip.PushFront(element)
		}
	} else {
		index := 0
		remain := ql.capNode - int(ql.head.zip.Length())
		for index < remain && index < len(elements) {
			ql.head.zip.PushFront(elements[index])
			index++
		}
		for index < len(elements) {
			newNode := &QLNode{
				zip: NewZipList(),
			}
			for i := 0; i < ql.capNode && index < len(elements); i++ {
				newNode.zip.PushFront(elements[index])
				index++
			}
			ql.head.prev = newNode
			newNode.next = ql.head
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

func (ql *QuickList) Length() int {
	return ql.len
}
