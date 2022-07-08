package event

import (
	"context"
	"sort"
)

type FrameEvent interface {
	GetSource() interface{}
	Name() string
}

type FrameListener interface {
	OnEvent(local context.Context, event FrameEvent) error
	Order() int
	WatchEvent() []FrameEvent
}

type Dispatcher struct {
	listeners map[string][]FrameListener
}

func (f *Dispatcher) DispatchEvent(local context.Context, event FrameEvent) {
	eventName := event.Name()
	if listeners, ok := f.listeners[eventName]; ok {
		for _, listener := range listeners {
			listener.OnEvent(local, event)
		}
	}
}

func (f *Dispatcher) AddEventListener(local context.Context, listener FrameListener) {

	events := listener.WatchEvent()

	for _, event := range events {
		eventName := event.Name()
		if listeners, ok := f.listeners[eventName]; ok {
			l := append(listeners, listener)
			sort.Slice(l, func(i, j int) bool {
				return l[i].Order() < l[j].Order()
			})
			f.listeners[eventName] = l
		} else {
			f.listeners[eventName] = []FrameListener{listener}
		}
	}

}
