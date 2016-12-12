package eui

import (
	"github.com/tryor/util/event"
)

type SelectEvent struct {
	event.Event
	Selected bool
}

func NewSelectEvent(source interface{}, selected bool) *SelectEvent {
	return &SelectEvent{Event: event.Event{Type: SELECT_EVENT_TYPE, Source: source}, Selected: selected}
}
