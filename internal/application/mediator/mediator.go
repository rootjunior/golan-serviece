package mediator

import (
	"context"
	"fmt"
	. "go-service/internal/core/interfaces"
	"reflect"
	"sync"
)

// ImplMediator реализует паттерн Mediator для маршрутизации команд, запросов и событий
// между компонентами приложения без прямых зависимостей между ними.
type ImplMediator struct {
	queries  map[reflect.Type]QueryHandler
	commands map[reflect.Type][]CommandHandler
	events   map[reflect.Type][]EventHandler
	bus      EventBus
	mu       sync.RWMutex
}

// NewMediator создаёт новый экземпляр Mediator.
// Принимает EventBus для асинхронной публикации событий.
func NewMediator(bus EventBus) Mediator {
	return &ImplMediator{
		queries:  make(map[reflect.Type]QueryHandler),
		commands: make(map[reflect.Type][]CommandHandler),
		events:   make(map[reflect.Type][]EventHandler),
		bus:      bus,
	}
}

// Queries

// RegisterQuery регистрирует обработчик для указанного типа запроса.
// queryPrototype — пустой экземпляр типа запроса, используется для получения reflect.Type.
// Повторная регистрация перезаписывает предыдущий обработчик.
func (m *ImplMediator) RegisterQuery(queryPrototype Query, handler QueryHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := reflect.TypeOf(queryPrototype)
	m.queries[t] = handler
}

// HandleQuery выполняет запрос и возвращает результат.
// Определяет обработчик по типу query и вызывает его Handle.
// Возвращает ошибку если обработчик для данного типа не зарегистрирован.
func (m *ImplMediator) HandleQuery(query Query) (Result, error) {
	t := reflect.TypeOf(query)
	m.mu.RLock()
	handler, ok := m.queries[t]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler registered for query type %v", t)
	}
	return handler.Execute(query) // было handler(query)
}

// Commands

// RegisterCommand регистрирует один или несколько обработчиков для указанного типа команды.
// commandPrototype — пустой экземпляр типа команды, используется для получения reflect.Type.
// Несколько вызовов RegisterCommand для одного типа добавляют обработчики к существующим.
func (m *ImplMediator) RegisterCommand(commandPrototype Command, handlers ...CommandHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := reflect.TypeOf(commandPrototype)
	m.commands[t] = append(m.commands[t], handlers...)
}

// HandleCommand выполняет команду через все зарегистрированные обработчики последовательно.
// Возвращает список результатов от каждого обработчика.
// Прерывает выполнение и возвращает ошибку если любой обработчик вернул ошибку.
// Возвращает ошибку если обработчики для данного типа не зарегистрированы.
func (m *ImplMediator) HandleCommand(ctx context.Context, command Command) ([]Result, error) {
	t := reflect.TypeOf(command)
	m.mu.RLock()
	handlers, ok := m.commands[t]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handlers registered for command type %v", t)
	}

	results := make([]Result, 0, len(handlers))
	for _, handler := range handlers {
		result, err := handler.Execute(ctx, command) // было handler(command)
		if err != nil {
			return results, fmt.Errorf("command handler error: %w", err)
		}
		results = append(results, result)
	}
	return results, nil
}

// Events

// RegisterEvent регистрирует один или несколько обработчиков для указанного типа события.
// eventPrototype — пустой экземпляр типа события, используется для получения reflect.Type.
// Несколько вызовов RegisterEvent для одного типа добавляют обработчики к существующим событиям.
func (m *ImplMediator) RegisterEvent(eventPrototype Event, handlers ...EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := reflect.TypeOf(eventPrototype)
	m.events[t] = append(m.events[t], handlers...)
}

// HandleEvent синхронно обрабатывает событие через все зарегистрированные обработчики.
// Прерывает выполнение и возвращает ошибку если любой обработчик вернул ошибку.
// Возвращает ошибку если обработчики для данного типа не зарегистрированы.
func (m *ImplMediator) HandleEvent(ctx context.Context, event Event) error {
	t := reflect.TypeOf(event)
	m.mu.RLock()
	handlers, ok := m.events[t]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no handlers registered for event type %v", t)
	}
	for _, handler := range handlers {
		if err := handler.Execute(ctx, event); err != nil { // было handler(event)
			return fmt.Errorf("event handler error: %w", err)
		}
	}
	return nil
}

// PublishEvents асинхронно отправляет события в стрим событий.
// Не блокирует основной поток — события обрабатываются воркерами в фоне.
// Использует неблокирующую отправку: если буфер канала заполнен, событие отбрасывается.
// Таким образом, метод предоставляет гарантию **at-most-once** (не более одного раза):
//   - Событие либо будет доставлено и обработано ровно один раз, либо потеряно совсем.
//   - Дублирование событий невозможно.
//   - Потеря возможна при переполнении буфера, а также при аварийном завершении процесса
//     (все непрочитанные события из канала исчезнут).
//
// Если ctx отменяется до отправки всех событий — возвращает ошибку,
// уже отправленные (переданные в канал) события продолжат обработку либо будут потеряны
// в случае сетевых или системных сбоев.
func (m *ImplMediator) PublishEvents(ctx context.Context, events ...Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			return fmt.Errorf("publish event cancelled: %w", ctx.Err())
		default:
			m.bus.Send(event)
		}
	}
	return nil
}
