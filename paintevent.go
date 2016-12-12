package eui

import (
	"github.com/tryor/util/event"
)

type PaintEvent struct {
	event.Event
}

func NewPaintEvent(source interface{}) *PaintEvent {
	return &PaintEvent{Event: event.Event{Type: PAINT_EVENT_TYPE, Source: source}}
}
