package engine

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// type Node struct {
// 	val  string
// 	prev *Node
// 	next *Node
// }

// type List struct {
// 	head *Node
// 	tail *Node
// 	len  int
// 	mu   sync.RWMutex
// }

// func NewNode(val string) *Node {
// 	return &Node{
// 		val: val,
// 	}
// }

// func NewList(args []string) *List {
// 	var res List

// 	head := NewNode(args[0])

// 	res.head = head
// 	res.tail = head
// 	res.len = len(args)

// 	for i := 1; i < len(args); i++ {
// 		newNode := NewNode(args[i])
// 		res.tail.next = newNode
// 		res.tail = newNode
// 	}

// 	return &res
// }

// func (l *List) LinkToList(newList *List) (headList *Node) {
// 	l.tail.next = newList.head
// 	l.tail = newList.tail
// 	headList = l.head
// 	l.len += newList.len
// 	return headList
// }

// func (l *List) GetElements() (elements []string) {
// 	curr := l.head
// 	for {
// 		if curr == nil {
// 			break
// 		}
// 		elements = append(elements, curr.val)
// 		curr = curr.next
// 	}
// 	return elements
// }

//===================================LIST=================================

func (k *KeyValue) SetList(key string, lPush bool, element string) (res int, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			list := NewZipList()
			res = 1
			list.PushBack(element)
			return NewObject(LIST, list), nil
		}

		obj := prev.(*Object)
		if obj.typ == LIST {
			oldList := obj.value.(*ZipList)
			oldList.mu.Lock()
			defer oldList.mu.Unlock()

			if lPush {
				oldList.PushFront(element)

			} else {
				oldList.PushBack(element)
			}
			res = int(oldList.Length())
			return obj, nil
		}
		return nil, ErrWrongType
	})
	return res, err
}

func (k *KeyValue) GetListBetween(key, start, stop string) (values []string, found bool, err error) {
	i64Start, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return values, false, ErrNotInteger
	}
	intStart := int(i64Start)

	i64Stop, err := strconv.ParseInt(stop, 10, 64)
	if err != nil {
		return values, false, ErrNotInteger
	}
	intStop := int(i64Stop)

	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return values, false, ErrInternal
		}
	}
	prev, ok := k.kv.Get(key)
	if ok {
		obj := prev.(*Object)
		if obj.typ != LIST {
			return values, found, ErrWrongType
		}
		oldList := obj.value.(*ZipList)
		oldList.mu.RLock()
		defer oldList.mu.RUnlock()
		if intStart < 0 {
			intStart += int(oldList.Length())
		}

		if intStop < 0 {
			intStop += int(oldList.Length())
		}

		if intStart >= int(oldList.Length()) {
			return values, false, err
		}

		count := 0
		offset := HEADER_SIZE
		listLen := binary.LittleEndian.Uint32(oldList.buf[0:4])
		for offset < int(listLen) {
			if count > int(intStop) {
				return values, true, err
			}

			if count >= int(intStart) {
				encoding := uint8(oldList.buf[offset+1])
				s := string(oldList.buf[offset+2 : offset+2+int(encoding)])
				values = append(values, s)
				offset += 2 + int(encoding)
			}
			count++
		}

		found = true
	}

	return values, found, err
}

func (k *KeyValue) PopList(key, count string, lpop bool) (values []string, poped bool, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return nil, false, ErrInternal
		}
	}

	temp, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return values, false, ErrNotInteger
	}
	countInt := int(temp)

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			return prev, ErrNotExists
		}

		obj := prev.(*Object)
		if obj.typ != LIST {
			return values, ErrWrongType
		}
		oldList := obj.value.(*ZipList)
		oldList.mu.RLock()
		defer oldList.mu.RUnlock()

		values = make([]string, 0, min(countInt, int(oldList.Length())))

		fmt.Println(oldList.Length())
		if countInt >= int(oldList.Length()) {
			if lpop {
				values = append(values, oldList.GetElements()...)
			} else {
				elements := oldList.GetElements()
				for i := len(elements) - 1; i >= 0; i-- {
					values = append(values, elements[i])
				}
			}
			return nil, nil
		} else {
			for range countInt {
				var value string
				if lpop {
					value = oldList.PopFront()
				} else {
					value = oldList.PopBack()
				}
				values = append(values, value)
			}
		}

		return obj, err
	})
	if err != nil {
		if err == ErrNotExists {
			return values, false, nil
		}
		return values, false, err
	}

	return values, true, err
}
