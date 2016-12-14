package eui

import (
	"github.com/tryor/commons/event"
)

type WheelEvent struct {
	MouseEvent
	Delta       int
	Orientation uint //Orientation
}

func NewWheelEvent(source interface{}, x, y int, buttons MButton, delta int, orientation uint) *WheelEvent {
	return &WheelEvent{MouseEvent: MouseEvent{Event: event.Event{Type: MOUSE_WHEEL_EVENT_TYPE, Source: source}, x: x, y: y, Buttons: buttons}, Delta: delta, Orientation: orientation}
}
