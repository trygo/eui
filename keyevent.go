package eui

import (
	"github.com/google/gxui"
	"github.com/tryor/commons/event"
)

type KeyEvent struct {
	event.Event
	KeySequence *KeySequence //键序列
	Key         gxui.KeyboardKey
	Char        rune //字符

	Modifier KeyboardModifier
}

//func NewKeyEvent(t event.Type, source interface{}, keys ...byte) *KeyEvent {
func NewKeyEvent(t event.Type, source interface{}, modifier KeyboardModifier, keys ...gxui.KeyboardKey) *KeyEvent {
	ke := &KeyEvent{Event: event.Event{Type: t, Source: source}, Modifier: modifier}
	ke.KeySequence = NewKeySequence(keys...)
	if len(keys) > 0 {
		ke.Key = keys[0]
	}
	return ke
}

func NewKeyCharEvent(t event.Type, source interface{}, char rune, modifier KeyboardModifier, keys ...gxui.KeyboardKey) *KeyEvent {
	ke := &KeyEvent{Event: event.Event{Type: t, Source: source}, Char: char, Modifier: modifier}
	ke.KeySequence = NewKeySequence(keys...)
	if len(keys) > 0 {
		ke.Key = keys[0]
	}
	return ke
}

func (ke *KeyEvent) Test(k gxui.KeyboardKey) bool {
	return ke.KeySequence.Test(k)
}

func (ke *KeyEvent) Tests(matchMode int, ks ...gxui.KeyboardKey) bool {
	return ke.KeySequence.Tests(matchMode, ks...)
}

func IsKeyEvent(t event.Type) bool {
	return t == KEY_PRESS_EVENT_TYPE || t == KEY_RELEASE_EVENT_TYPE
}
