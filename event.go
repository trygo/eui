package eui

import (
	"reflect"
)

type EventSubscription interface {
	Unlisten()
}

type Event interface {
	Fire(args ...interface{})
	Listen(interface{}) EventSubscription
	ParameterTypes() []reflect.Type
}

type SimpleEvent struct {
	EventBase
}

func CreateEvent(signature interface{}) Event {
	e := &SimpleEvent{}
	e.init(signature)
	return e
}
