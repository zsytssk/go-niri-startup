package utils

import "fmt"

type Event struct {
	Listeners map[string][]Listener
	Counter   int
}

func (state *Event) OnEvent(eventName string, fn func(interface{})) {
	state.Counter++
	id := fmt.Sprintf("%s_%d", eventName, state.Counter)
	state.Listeners[eventName] = append(state.Listeners[eventName], Listener{
		id,
		fn,
	})
}

func (state *Event) OnceEvent(eventName string, fn func(interface{})) {
	var listener Listener
	listener = Listener{
		fmt.Sprintf("%s_%d", eventName, state.Counter),
		func(val interface{}) {
			fn(val)
			state.OffEvent(eventName, listener)
		},
	}
	state.Listeners[eventName] = append(state.Listeners[eventName], listener)
}

func (state *Event) OffEvent(eventName string, fn Listener) {
	index := -1
	for i, item := range state.Listeners[eventName] {
		if item.Id == fn.Id {
			index = i
		}
	}
	state.Listeners[eventName] = append(state.Listeners[eventName][:index], state.Listeners[eventName][index+1:]...)
}

func (state *Event) TriggerEvent(eventName string, data interface{}) {
	for _, item := range state.Listeners[eventName] {
		item.Fn(data)
	}
}

type Listener struct {
	Id string
	Fn func(interface{})
}
