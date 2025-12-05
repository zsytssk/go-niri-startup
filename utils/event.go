package utils

import "fmt"

type Event struct {
	Listeners map[string][]Listener
	Counter   int
}

func (e *Event) OnEvent(eventName string, fn func(interface{})) func() {
	e.Counter++
	id := fmt.Sprintf("%s_%d", eventName, e.Counter)
	listener := Listener{
		id,
		fn,
	}
	e.Listeners[eventName] = append(e.Listeners[eventName], listener)

	return func() {
		e.OffEvent(eventName, listener)
	}
}

func (e *Event) OnceEvent(eventName string, fn func(interface{})) {
	var listener Listener
	listener = Listener{
		fmt.Sprintf("%s_%d", eventName, e.Counter),
		func(val interface{}) {
			fn(val)
			e.OffEvent(eventName, listener)
		},
	}
	e.Listeners[eventName] = append(e.Listeners[eventName], listener)
}

func (e *Event) OffEvent(eventName string, fn Listener) {
	index := -1
	for i, item := range e.Listeners[eventName] {
		if item.Id == fn.Id {
			index = i
		}
	}
	e.Listeners[eventName] = append(e.Listeners[eventName][:index], e.Listeners[eventName][index+1:]...)
}

func (e *Event) TriggerEvent(eventName string, data interface{}) {
	for _, item := range e.Listeners[eventName] {
		item.Fn(data)
	}
}

type Listener struct {
	Id string
	Fn func(interface{})
}
