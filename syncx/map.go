package syncx

import (
	"errors"
	"fmt"
	"sync"
)

type Map[T any] interface {
	Load(key string) (*T, error)
	Store(key string, value T)
}

type impl[T any] struct {
	m *sync.Map
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
