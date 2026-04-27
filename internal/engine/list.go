package engine

import (
	"strconv"
	"sync"
)

type Node struct {
	val  string
	prev *Node
	next *Node
}

type List struct {
	head *Node
	tail *Node
	len  int
	mu   sync.RWMutex
}

func NewNode(val string) *Node {
	return &Node{
		val: val,
	}
}

func NewList(args []string) *List {
	var res List

	head := NewNode(args[0])

	res.head = head
	res.tail = head
	res.len = len(args)

	for i := 1; i < len(args); i++ {
		newNode := NewNode(args[i])
		res.tail.next = newNode
		res.tail = newNode
	}

	return &res
}

func (l *List) LinkToList(newList *List) (headList *Node) {
	l.tail.next = newList.head
	l.tail = newList.tail
	headList = l.head
	l.len += newList.len
	return headList
}

//===================================LIST=================================

func (k *KeyValue) SetList(key string, lPush bool, elements []string) (res int, err error) {
	list := NewList(elements)
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			res = list.len
			return NewObject(LIST, list), nil
		}

		obj := prev.(*Object)
		if obj.typ == LIST {
			oldList := obj.value.(*List)
			oldList.mu.Lock()
			defer oldList.mu.Unlock()

			if lPush {
				list.tail.next = oldList.head
				oldList.head = list.head
				oldList.len += list.len
			} else {
				oldList.LinkToList(list)
			}
			res = oldList.len
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
		oldList := obj.value.(*List)
		oldList.mu.RLock()
		defer oldList.mu.RUnlock()
		if intStart < 0 {
			intStart += oldList.len
		}

		if intStop < 0 {
			intStop += oldList.len
		}

		if intStart >= oldList.len {
			return values, false, err
		}

		cur := oldList.head
		count := 0
		for {
			if cur == nil {
				break
			}

			if count > int(intStop) {
				return values, true, err
			}

			if count >= int(intStart) {
				values = append(values, cur.val)
			}
			cur = cur.next
			count++
		}
		found = true
	}

	return values, found, err
}
