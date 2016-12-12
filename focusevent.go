package eui

import (
	"github.com/tryor/util/event"
)

type FocusEvent struct {
	event.Event
	Focus bool
}

func NewFocusEvent(source interface{}, focus bool) *FocusEvent {
	return &FocusEvent{Event: event.Event{Type: FOCUS_EVENT_TYPE, Source: source}, Focus: focus}
}
