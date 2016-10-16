package eui

import (
	//	"github.com/google/gxui"
	"github.com/trygo/util/event"
)

type IMouseEvent interface {
	event.IEvent
	Pressed(button MButton) bool
	X() int
	Y() int
	SetX(x int)
	SetY(y int)
	GetButtons() MButton
	GetKeySequence() *KeySequence
	GetModifier() KeyboardModifier
}

type MouseEvent struct {
	event.Event
	x, y    int //当前层或当前元素上的坐标
	Buttons MButton

	KeySequence *KeySequence
	PageEvent   IMouseEvent //在画布上的鼠标事件
	LayerEvent  IMouseEvent //在层上的鼠标事件

	Modifier KeyboardModifier
}

func NewMouseEvent(t event.Type, source interface{}, x, y int, buttons MButton, modifier KeyboardModifier) *MouseEvent {
	return &MouseEvent{Event: event.Event{Type: t, Source: source}, x: x, y: y, Buttons: buttons, Modifier: modifier}
}

/**
 * 检查按钮状态
 */
func (m *MouseEvent) Pressed(button MButton) bool {
	return (m.Buttons & button) > 0
}

func (m *MouseEvent) X() int {
	return m.x
}
func (m *MouseEvent) Y() int {
	return m.y
}

func (m *MouseEvent) SetX(x int) {
	m.x = x
}
func (m *MouseEvent) SetY(y int) {
	m.y = y
}

func (m *MouseEvent) GetButtons() MButton {
	return m.Buttons
}

func (m *MouseEvent) GetKeySequence() *KeySequence {
	return m.KeySequence
}

func (m *MouseEvent) GetModifier() KeyboardModifier {
	return m.Modifier
}

func IsMouseEvent(t event.Type) bool {
	return t > MOUSE_EVENT_MIN_TYPE && t < MOUSE_EVENT_MAX_TYPE
}
