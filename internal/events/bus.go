package events

import (
	"sync"
	"time"
)

type Handler func(Event)

type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *Bus) Publish(eventType string, data any) {
	event := Event{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	b.mu.RLock()
	handlers := append([]Handler(nil), b.subscribers[eventType]...)
	wildcard := append([]Handler(nil), b.subscribers["*"]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
	for _, h := range wildcard {
		h(event)
	}
}
