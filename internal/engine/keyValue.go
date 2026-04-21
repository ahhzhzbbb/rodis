package engine

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var ErrValueNotInteger = errors.New("value is not an integer or out of range")

type KeyValue struct {
	kv   Map
	et   Map
	keys []string
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

func (k *KeyValue) Get(key string) (string, bool) {

	// if k.CheckExpireKey(key) {
	// 	k.Del(key)
	// 	return "", false
	// }

	temp, exists := k.kv.Get(key)
	res, _ := temp.(string)
	return res, exists
}

func (k *KeyValue) Set(key, value string) {
	k.et.Delete(key)
	k.kv.Set(key, value)
}

func (k *KeyValue) Incr(key string) (int64, error) {
	var result int64

	if k.CheckExpireKey(key) {
		k.Del(key)
	}

	_, err := k.kv.Compute(key, func(prev interface{}, exists bool) (interface{}, error) {
		if !exists {
			result = 1
			return "1", nil
		}

		current, ok := prev.(string)
		if !ok {
			return nil, ErrValueNotInteger
		}

		i64, err := strconv.ParseInt(current, 10, 64)
		if err != nil {
			return nil, ErrValueNotInteger
		}

		i64++
		result = i64
		return strconv.FormatInt(i64, 10), nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (k *KeyValue) Append(key string, value string) (int, error) {
	var result int
	if k.CheckExpireKey(key) {
		k.Del(key)
	}

	_, err := k.kv.Compute(key, func(prev interface{}, exists bool) (newValue interface{}, err error) {
		if !exists {
			result = len(value)
			return value, nil
		}

		current, ok := prev.(string)
		if !ok {
			return nil, ErrValueNotInteger
		}

		current = current + value
		result = len(current)
		return current, nil
	})
	if err != nil {
		return 0, err
	}

	return result, nil
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

func (k *KeyValue) DelValue(key string) {
	k.kv.Delete(key)
}

func (k *KeyValue) SetExpireTime(key string, t time.Time) bool {
	if _, exists := k.kv.Get(key); !exists {
		return false
	}

	if k.CheckExpireKey(key) {
		k.kv.Delete(key)
		return false
	}

	k.et.Set(key, Entry{time: t, index: len(k.keys)})
	k.keys = append(k.keys, key)

	return true
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

func (k *KeyValue) DeleteKeys(sampleSize int, expireThreshold float64, timeBudgetMs int) {
	startTime := time.Now()
	budget := time.Duration(timeBudgetMs) * time.Millisecond

	for {
		removedKeyNums := 0.0

		totalKeys := len(k.keys)
		if totalKeys == 0 {
			return
		}

		actualSampleSize := min(totalKeys, sampleSize)

		for range actualSampleSize {
			idx := rand.Intn(len(k.keys))
			key := k.keys[idx]

			if k.CheckExpireKey(key) {
				if ok := k.Del(key); ok {
					fmt.Printf("Deleted Key: %s\n", key)
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

func (k *KeyValue) removeAtIndex(index int) {
	if index < 0 || index >= len(k.keys) {
		return
	}

	lastKey := k.keys[len(k.keys)-1]
	val, ok := k.et.Get(lastKey)
	if !ok {
	}
	entryValue := val.(Entry)
	entryValue.index = index
	k.et.Set(lastKey, entryValue)

	k.keys[index] = lastKey
	k.keys = k.keys[:len(k.keys)-1]
}

/*
hàm này dùng để test thôi, test xong thì xóa đi cho đỡ tốn tài nguyên
// */
// func (k *KeyValue) LenET() int {
// 	return k.et.Len()
// }
