package utils

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Listener struct {
	Id string
	Fn func(interface{})
}

type Event struct {
	listeners map[string]map[string]Listener
	counter   uint64
	mu        sync.RWMutex
}

func NewEvent() *Event {
	return &Event{
		listeners: make(map[string]map[string]Listener, 0),
	}
}

func (e *Event) OnEvent(eventName string, fn func(interface{})) func() {
	idNum := atomic.AddUint64(&e.counter, 1)
	id := fmt.Sprintf("%s:%d", eventName, idNum)
	listener := Listener{
		id,
		fn,
	}
	e.mu.Lock()
	if e.listeners[eventName] == nil {
		e.listeners[eventName] = make(map[string]Listener)
	}
	e.listeners[eventName][id] = listener
	e.mu.Unlock()
	return func() {
		e.OffEvent(eventName, id)
	}
}

func (e *Event) OnceEvent(eventName string, fn func(interface{})) func() {
	var once sync.Once
	var offFn func()
	fnNew := func(val interface{}) {
		once.Do(func() {
			fn(val)
			if offFn != nil {
				offFn()
			}
		})

	}
	offFn = e.OnEvent(eventName, fnNew)
	return offFn
}

func (e *Event) OffEvent(eventName string, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.listeners[eventName] == nil {
		return
	}
	delete(e.listeners[eventName], id)
}

func (e *Event) TriggerEvent(eventName string, data interface{}) {
	e.mu.RLock()
	if e.listeners[eventName] == nil {
		e.mu.RUnlock()
		return
	}

	copyListeners := make([]Listener, 0, len(e.listeners[eventName]))
	for _, item := range e.listeners[eventName] {
		copyListeners = append(copyListeners, item)
	}
	e.mu.RUnlock()

	for _, item := range copyListeners {
		item.Fn(data)
	}
}

func (e *Event) Clear(eventName string) {
	e.mu.Lock()
	delete(e.listeners, eventName)
	e.mu.Unlock()
}

func (e *Event) Count(eventName string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if listeners, ok := e.listeners[eventName]; ok {
		return len(listeners)
	}
	return 0
}
