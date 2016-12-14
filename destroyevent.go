package eui

import (
	"github.com/tryor/commons/event"
)

type DestroyEvent struct {
	event.Event
}

func NewDestroyEvent(source interface{}) *DestroyEvent {
	return &DestroyEvent{Event: event.Event{Type: DESTROY_EVENT_TYPE, Source: source}}
}
