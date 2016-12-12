package eui

import (
	"github.com/tryor/util/event"
)

type MovedEvent struct {
	event.Event
	Dx, Dy int
	Angle  float32
}

func NewMovedEvent(source interface{}) *MovedEvent {
	return &MovedEvent{Event: event.Event{Type: MOVED_EVENT_TYPE, Source: source}}
}
