package engine

type Node struct {
	val  string
	prev *Node
	next *Node
}

type List struct {
	head *Node
	tail *Node
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

	for i := 1; i < len(args); i++ {
		newNode := NewNode(args[i])
		res.tail.next = newNode
		res.tail = newNode
	}

	return &res
}

func (l *List) LinkToList(newList *List) {
	l.tail.next = newList.head
	l.tail = newList.tail
}

func (l *List) Len() int {
	res := 0
	cur := l.head
	for {
		if cur == nil {
			break
		}
		res++
		cur = cur.next
	}
	return res
}

func (k *KeyValue) SetList(key string, elements []string) {
	list := NewList(elements)
	oldValue, exist := k.Get(key)
	if !exist {
		obj := NewObject(LIST, list)
		k.Set(key, obj)
	} else {
		temp := oldValue.(*Object)
		if temp.typ == LIST {
			oldList := temp.value.(*List)
			oldList.LinkToList(list)
		}
	}
}

func (k *KeyValue) GetList(key string) []string {
	var res []string
	val, ok := k.Get(key)
	if ok {
		obj := val.(*Object)
		oldList := obj.value.(*List)
		cur := oldList.head
		for {
			if cur == nil {
				break
			}

			res = append(res, cur.val)

			cur = cur.next
		}
	}

	return res
}
