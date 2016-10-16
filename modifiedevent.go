package eui

import (
	"github.com/trygo/util/event"
)

type ModifiedEvent struct {
	event.Event
	ModifiedSupport *ModifiedSupport
}

func NewModifiedEvent(source interface{}, support *ModifiedSupport) *ModifiedEvent {
	//func NewModifiedEvent(source interface{}) *ModifiedEvent {
	me := &ModifiedEvent{Event: event.Event{Type: MODIFIED_EVENT_TYPE, Source: source}}
	me.ModifiedSupport = support
	return me
}
