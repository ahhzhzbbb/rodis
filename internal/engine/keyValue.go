package engine

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var ErrNotInteger = errors.New("value is not an integer or out of range")
var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
var ErrInternal = errors.New("something go wrong with server")

// var ErrOverflow = errors.New("")
// var ErrInvalidTTL = errors.New("")

type KeyValue struct {
	kv   Map
	et   Map
	keys []string
	mu   sync.RWMutex
}

type Entry struct {
	time  time.Time
	index int
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		kv: *New(1024),
		et: *New(1024),
	}
}

//=======================================STRING================================================

func (k *KeyValue) GetString(key string) (value string, found bool, err error) {
	if k.CheckExpireKey(key) {
		k.Del(key)
		return "", false, nil
	}

	prev, ok := k.kv.Get(key)
	if !ok {
		return "", false, err
	}

	temp := prev.(*Object)
	if temp.typ == STRING {
		value = temp.value.(string)
		found = true
		err = nil
	} else {
		err = ErrWrongType
	}
	return value, found, err
}

func (k *KeyValue) SetString(key, value string) error {
	k.Del(key)

	obj := NewObject(STRING, value)
	k.kv.Set(key, obj)
	return nil
}

func (k *KeyValue) IncrString(key string) (newValue int64, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (any, error) {
		if !exists {
			newValue = 1
			return NewObject(STRING, strconv.FormatInt(newValue, 10)), nil
		}

		current := prev.(*Object)
		if current.typ != STRING {
			return nil, ErrWrongType
		}

		i64, err := strconv.ParseInt(current.value.(string), 10, 64)
		if err != nil {
			return nil, ErrNotInteger
		}
		if i64 == math.MaxInt64 {
			return nil, ErrNotInteger
		}

		i64++
		newValue = i64
		obj := NewObject(STRING, strconv.FormatInt(newValue, 10))
		return obj, nil
	})
	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func (k *KeyValue) AppendString(key string, value string) (newLen int, err error) {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return 0, ErrInternal
		}
	}

	_, err = k.kv.Compute(key, func(prev any, exists bool) (newValue any, err error) {
		if !exists {
			newLen = len(value)
			return NewObject(STRING, value), nil
		}

		current := prev.(*Object)
		if current.typ != STRING {
			return nil, ErrWrongType
		}

		current.value = current.value.(string) + value
		newLen = len(current.value.(string))
		return current, nil
	})
	if err != nil {
		return 0, err
	}

	return newLen, nil
}

//======================================TTL\KEY===================================================

func (k *KeyValue) removeAtIndex(index int) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if index < 0 || index >= len(k.keys) {
		return
	}

	lastKey := k.keys[len(k.keys)-1]
	val, ok := k.et.Get(lastKey)
	if !ok {
		return
	}
	entryValue := val.(Entry)
	entryValue.index = index
	k.et.Set(lastKey, entryValue)

	k.keys[index] = lastKey
	k.keys = k.keys[:len(k.keys)-1]
}

func (k *KeyValue) Del(key string) bool {
	_, ok1 := k.kv.Delete(key)

	temp, ok2 := k.et.Get(key)
	if ok2 {
		value := temp.(Entry)
		k.removeAtIndex(value.index)
		k.et.Delete(key)
	}

	return ok1
}

func (k *KeyValue) SetExpireTime(key string, t time.Time) bool {
	if _, exists := k.kv.Get(key); !exists {
		return false
	}

	if k.CheckExpireKey(key) {
		k.kv.Delete(key)
		return false
	}

	if val, ok := k.et.Get(key); ok {
		temp := val.(Entry)
		oldIdx := temp.index
		k.et.Set(key, Entry{time: t, index: oldIdx})
	} else {
		k.mu.Lock()
		k.et.Set(key, Entry{time: t, index: len(k.keys)})
		k.keys = append(k.keys, key)
		k.mu.Unlock()
	}

	return true
}

func (k *KeyValue) CheckExistsKey(key string) bool {
	if k.CheckExpireKey(key) {
		ok := k.Del(key)
		if !ok {
			return false
		}
	}
	_, exists := k.kv.Get(key)
	return exists
}

func (k *KeyValue) CheckExpireKey(key string) bool {
	val, ok := k.et.Get(key)
	temp, _ := val.(Entry)
	t := temp.time

	if ok && t.Before(time.Now()) {
		return true
	}
	return false
}

func (k *KeyValue) ActiveExpiration(sampleSize int, expireThreshold float64, timeBudgetMs int) {
	startTime := time.Now()
	budget := time.Duration(timeBudgetMs) * time.Millisecond

	for {
		removedKeyNums := 0.0
		k.mu.RLock()
		totalKeys := len(k.keys)
		k.mu.RUnlock()
		if totalKeys == 0 {
			return
		}

		actualSampleSize := min(totalKeys, sampleSize)

		for range actualSampleSize {
			k.mu.RLock()
			if len(k.keys) == 0 {
				k.mu.RUnlock()
				break
			}
			idx := rand.Intn(len(k.keys))
			key := k.keys[idx]
			k.mu.RUnlock()

			if k.CheckExpireKey(key) {
				if ok := k.Del(key); ok {
					// fmt.Printf("Deleted Key: %s\n", key)
					removedKeyNums++
				}
			}
		}
		ratio := 0.0
		if actualSampleSize > 0 {
			ratio = removedKeyNums / float64(actualSampleSize)
		}

		if ratio < expireThreshold || time.Since(startTime) >= budget {
			return
		}
	}
}

/*
hàm này dùng để test thôi, test xong thì xóa đi cho đỡ tốn tài nguyên
// */
// func (k *KeyValue) LenET() int {
// 	return k.et.Len()
// }
