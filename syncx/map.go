package syncx

import (
	"errors"
	"fmt"
	"sync"
)

type Map[T any] interface {
	Load(key string) (*T, error)
	Store(key string, value T)
	Delete(key string)
	Export() map[string]T
	Size() int
}

type impl[T any] struct {
	m *sync.Map
}

func (i *impl[T]) Delete(key string) {
	i.m.Delete(key)
}

func (i *impl[T]) Export() map[string]T {
	m := map[string]T{}
	i.m.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(T)
		return true
	})
	return m
}

func (i *impl[T]) Size() int {
	count := 0
	i.m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func Define[T any]() Map[T] {
	return &impl[T]{
		m: &sync.Map{},
	}
}

func (i *impl[T]) Load(key string) (*T, error) {
	val, ok := i.m.Load(key)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s not found", key))
	}
	t, ok := val.(T)
	if !ok {
		return nil, errors.New("type not match")
	}
	return &t, nil
}

func (i *impl[T]) Store(key string, value T) {
	i.m.Store(key, value)
}
