package main

import (
	"fmt"
	"reflect"
	"sync"
)

// QueryHandlerFunc — тип функции обработчика запроса
type QueryHandlerFunc func(query interface{}) (interface{}, error)

// Mediator хранит мапу query -> handler
type Mediator struct {
	queries map[reflect.Type]QueryHandlerFunc
	mu      sync.RWMutex
}

// NewMediator создаёт новый медиатор
func NewMediator() *Mediator {
	return &Mediator{
		queries: make(map[reflect.Type]QueryHandlerFunc),
	}
}

// RegisterQuery регистрирует обработчик для типа query
func (m *Mediator) RegisterQuery(queryPrototype interface{}, handler QueryHandlerFunc) {
	t := reflect.TypeOf(queryPrototype)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queries[t] = handler
}

// Execute вызывает обработчик для переданного запроса
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
