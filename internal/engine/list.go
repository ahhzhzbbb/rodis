package engine

import (
	"strconv"
)

//===================================LIST=================================

func (k *KeyValue) SetList(key string, lPush bool, elements []string) (res int, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			list := NewQuickList(0, elements)
			res = len(elements)
			return NewObject(LIST, list), nil
		}

		obj := prev.(*Object)
		if obj.typ == LIST {
			oldList := obj.value.(*QuickList)
			oldList.mu.Lock()
			defer oldList.mu.Unlock()

			if lPush {
				oldList.PushFront(elements)
			} else {
				oldList.PushBack(elements)
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
		oldList := obj.value.(*QuickList)
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

		if intStop >= int(oldList.Length()) {
			intStop = int(oldList.Length()) - 1
		}

		if intStart > intStop {
			return values, true, err
		}

		values = oldList.GetElements()[intStart : intStop+1]

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
		oldList := obj.value.(*QuickList)
		oldList.mu.RLock()
		defer oldList.mu.RUnlock()

		values = make([]string, 0, min(countInt, int(oldList.Length())))

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
			if lpop {
				var temp []string
				for range countInt {
					temp = append(values, oldList.PopFront())
				}
				for i := len(temp) - 1; i >= 0; i-- {
					values = append(values, temp[i])
				}
			} else {
				for range countInt {
					values = append(values, oldList.PopBack())
				}
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

func (k *KeyValue) ListInsert(key string, position int, pivot, value string) (count int, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			return prev, ErrNotExists
		}

		obj := prev.(*Object)
		if obj.typ != LIST {
			return nil, ErrWrongType
		}
		oldList := obj.value.(*QuickList)
		oldList.mu.Lock()
		defer oldList.mu.Unlock()

		count = oldList.Insert(position, pivot, value)
		return obj, nil
	})
	if err != nil {
		if err == ErrNotExists {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}
