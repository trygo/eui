package eui

import (
	"syscall"
)

type ICursor interface {
	SetStyle(style syscall.Handle)
	GetStyle() syscall.Handle
}

type Cursor struct {
	Style syscall.Handle
}

func NewCursor(style syscall.Handle) *Cursor {
	return &Cursor{Style: style}
}

func (this *Cursor) SetStyle(style syscall.Handle) {
	this.Style = style
}

func (this *Cursor) GetStyle() syscall.Handle {
	return this.Style
}
