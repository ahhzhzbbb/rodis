package engine

import "strconv"

type Node struct {
	val  string
	prev *Node
	next *Node
}

type List struct {
	head *Node
	tail *Node
	len  int
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

// func (l *List) Len() int {
// 	res := 0
// 	cur := l.head
// 	for {
// 		if cur == nil {
// 			break
// 		}
// 		res++
// 		cur = cur.next
// 	}
// 	return res
// }

//===================================LIST=================================

func (k *KeyValue) SetList(key string, lPush bool, elements []string) (res int, err error) {
	list := NewList(elements)
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}
	oldValue, exist := k.kv.Get(key)
	if !exist {
		obj := NewObject(LIST, list)
		k.kv.Set(key, obj)
		res = list.len
	} else {
		temp := oldValue.(*Object)
		if temp.typ == LIST {
			oldList := temp.value.(*List)
			if lPush {
				oldList.head = list.LinkToList(oldList)
			} else {
				oldList.LinkToList(list)
			}
			return oldList.len, err
		} else {
			err = ErrWrongType
		}
	}
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
