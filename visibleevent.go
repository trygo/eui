package eui

import (
	"github.com/tryor/commons/event"
)

type VisibleEvent struct {
	event.Event
	Visible bool
}

func NewVisibleEvent(source interface{}, visible bool) *VisibleEvent {
	return &VisibleEvent{Event: event.Event{Type: VISIBLE_EVENT_TYPE, Source: source}, Visible: visible}
}
