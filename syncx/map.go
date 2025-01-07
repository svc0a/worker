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

type instance[T any] struct {
	m *sync.Map
}

func (i *instance[T]) Delete(key string) {
	i.m.Delete(key)
}

func (i *instance[T]) Export() map[string]T {
	m := map[string]T{}
	i.m.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(T)
		return true
	})
	return m
}

func (i *instance[T]) Size() int {
	count := 0
	i.m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func Define[T any]() Map[T] {
	return &instance[T]{
		m: &sync.Map{},
	}
}

func (i *instance[T]) Load(key string) (*T, error) {
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

func (i *instance[T]) Store(key string, value T) {
	i.m.Store(key, value)
}
