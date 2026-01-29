package main

import (
	"fmt"
	"reflect"
	"sync"
)

type QueryHandlerFunc func(query interface{}) (interface{}, error)

type Mediator struct {
	queries map[reflect.Type]QueryHandlerFunc
	mu      sync.RWMutex
}

func NewMediator() *Mediator {
	return &Mediator{
		queries: make(map[reflect.Type]QueryHandlerFunc),
	}
}

func (m *Mediator) RegisterQuery(queryPrototype interface{}, handler QueryHandlerFunc) {
	t := reflect.TypeOf(queryPrototype)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queries[t] = handler
}

func (m *Mediator) Execute(query interface{}) (interface{}, error) {
	t := reflect.TypeOf(query)
	m.mu.RLock()
	handler, ok := m.queries[t]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler registered for query type %v", t)
	}
	return handler(query)
}
