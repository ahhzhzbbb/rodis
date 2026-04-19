package engine

import (
	"errors"
	"strconv"
	"time"
)

var ErrValueNotInteger = errors.New("value is not an integer or out of range")

type KeyValue struct {
	kv Map
	et Map
	// Mu sync.RWMutex
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		kv: *New(1024),
		et: *New(1024),
	}
}

func (k *KeyValue) Get(key string) (string, bool) {

	if k.CheckExpireKey(key) {
		return "", false
	}

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

	_ = k.CheckExpireKey(key)

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
	_ = k.CheckExpireKey(key)

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
	_, ok := k.kv.Delete(key)
	k.et.Delete(key)
	return ok
}

func (k *KeyValue) DelValue(key string) {
	k.kv.Delete(key)
}

func (k *KeyValue) SetExpireTime(key string, t time.Time) bool {
	if _, exists := k.kv.Get(key); !exists {
		return false
	}

	k.et.Set(key, t)

	return true
}

func (k *KeyValue) CheckExpireKey(key string) bool {
	val, ok := k.et.Get(key)
	t, _ := val.(time.Time)

	if ok && t.Before(time.Now()) {
		k.kv.Delete(key)
		k.et.Delete(key)
		return true
	}
	return false
}
