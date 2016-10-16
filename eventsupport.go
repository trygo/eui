package eui

import (
	"github.com/google/gxui"
	"github.com/trygo/util/event"
)

type IEventSupport interface {
	FireEvent(e event.IEvent) bool

	OnMouseMove(func(*MouseEvent)) gxui.EventSubscription
	OnMouseDown(func(*MouseEvent)) gxui.EventSubscription
	OnMouseUp(func(*MouseEvent)) gxui.EventSubscription
	OnMouseDoubleClick(func(*MouseEvent)) gxui.EventSubscription
	OnMouseEnter(func(*MouseEvent)) gxui.EventSubscription
	OnMouseExit(func(*MouseEvent)) gxui.EventSubscription
	OnMouseWheel(func(*WheelEvent)) gxui.EventSubscription

	OnKeyPress(func(*KeyEvent)) gxui.EventSubscription
	OnKeyRelease(func(*KeyEvent)) gxui.EventSubscription
	OnKeyRepeat(func(*KeyEvent)) gxui.EventSubscription
	OnKeyChar(func(*KeyEvent)) gxui.EventSubscription

	OnSelect(func(*SelectEvent)) gxui.EventSubscription
	OnModified(func(*ModifiedEvent)) gxui.EventSubscription
	OnPaint(func(*PaintEvent)) gxui.EventSubscription
	OnMoved(func(*MovedEvent)) gxui.EventSubscription
	OnFocus(func(*FocusEvent)) gxui.EventSubscription

	OnDestroy(func(*DestroyEvent)) gxui.EventSubscription

	//OnMouseScroll(func(*MouseEvent)) gxui.EventSubscription
}

type EventSupport struct {
	//	Self interface{}

	onMouseMoveEvent   gxui.Event
	onMouseDownEvent   gxui.Event
	onMouseUpEvent     gxui.Event
	onMouseDoubleClick gxui.Event
	onMouseEnterEvent  gxui.Event
	onMouseExitEvent   gxui.Event
	onMouseWheelEvent  gxui.Event
	//	onMouseScrollEvent gxui.Event
	onKeyPressEvent   gxui.Event
	onKeyReleaseEvent gxui.Event
	onKeyRepeatEvent  gxui.Event
	onKeyCharEvent    gxui.Event
	onSelectEvent     gxui.Event
	onModifiedEvent   gxui.Event
	onPaintEvent      gxui.Event
	onMovedEvent      gxui.Event
	onFocusEvent      gxui.Event
	onDestroyEvent    gxui.Event
}

//func initEvent(esupport *EventSupport) {
//	esupport.onMouseMoveEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseDownEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseUpEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseDoubleClick = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseEnterEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseExitEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	//esupport.onMouseScrollEvent = gxui.CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseWheelEvent = gxui.CreateEvent(func(*WheelEvent) {})

//	esupport.onKeyPressEvent = gxui.CreateEvent(func(*KeyEvent) {})
//	esupport.onKeyReleaseEvent = gxui.CreateEvent(func(*KeyEvent) {})
//	esupport.onKeyCharEvent = gxui.CreateEvent(func(*KeyEvent) {})
//	esupport.onSelectEvent = gxui.CreateEvent(func(*SelectEvent) {})
//	esupport.onModifiedEvent = gxui.CreateEvent(func(*ModifiedEvent) {})
//	esupport.onPaintEvent = gxui.CreateEvent(func(*PaintEvent) {})
//	esupport.onMovedEvent = gxui.CreateEvent(func(*MovedEvent) {})
//	esupport.onFocusEvent = gxui.CreateEvent(func(*FocusEvent) {})
//}

func (this *EventSupport) OnMouseMove(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseMoveEvent == nil {
		this.onMouseMoveEvent = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseMoveEvent.Listen(f)
}

func (this *EventSupport) OnMouseDown(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseDownEvent == nil {
		this.onMouseDownEvent = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseDownEvent.Listen(f)
}

func (this *EventSupport) OnMouseUp(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseUpEvent == nil {
		this.onMouseUpEvent = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseUpEvent.Listen(f)
}

func (this *EventSupport) OnMouseDoubleClick(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseDoubleClick == nil {
		this.onMouseDoubleClick = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseDoubleClick.Listen(f)
}

func (this *EventSupport) OnMouseEnter(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseEnterEvent == nil {
		this.onMouseEnterEvent = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseEnterEvent.Listen(f)
}

func (this *EventSupport) OnMouseExit(f func(*MouseEvent)) gxui.EventSubscription {
	if this.onMouseExitEvent == nil {
		this.onMouseExitEvent = gxui.CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseExitEvent.Listen(f)
}

func (this *EventSupport) OnMouseWheel(f func(*WheelEvent)) gxui.EventSubscription {
	if this.onMouseWheelEvent == nil {
		this.onMouseWheelEvent = gxui.CreateEvent(func(*WheelEvent) {})
	}
	return this.onMouseWheelEvent.Listen(f)
}

func (this *EventSupport) OnKeyPress(f func(*KeyEvent)) gxui.EventSubscription {
	if this.onKeyPressEvent == nil {
		this.onKeyPressEvent = gxui.CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyPressEvent.Listen(f)
}

func (this *EventSupport) OnKeyRelease(f func(*KeyEvent)) gxui.EventSubscription {
	if this.onKeyReleaseEvent == nil {
		this.onKeyReleaseEvent = gxui.CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyReleaseEvent.Listen(f)
}

func (this *EventSupport) OnKeyRepeat(f func(*KeyEvent)) gxui.EventSubscription {
	if this.onKeyRepeatEvent == nil {
		this.onKeyRepeatEvent = gxui.CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyRepeatEvent.Listen(f)
}

func (this *EventSupport) OnKeyChar(f func(*KeyEvent)) gxui.EventSubscription {
	if this.onKeyCharEvent == nil {
		this.onKeyCharEvent = gxui.CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyCharEvent.Listen(f)
}
func (this *EventSupport) OnSelect(f func(*SelectEvent)) gxui.EventSubscription {
	if this.onSelectEvent == nil {
		this.onSelectEvent = gxui.CreateEvent(func(*SelectEvent) {})
	}
	return this.onSelectEvent.Listen(f)
}
func (this *EventSupport) OnModified(f func(*ModifiedEvent)) gxui.EventSubscription {
	if this.onModifiedEvent == nil {
		this.onModifiedEvent = gxui.CreateEvent(func(*ModifiedEvent) {})
	}
	return this.onModifiedEvent.Listen(f)
}
func (this *EventSupport) OnPaint(f func(*PaintEvent)) gxui.EventSubscription {
	if this.onPaintEvent == nil {
		this.onPaintEvent = gxui.CreateEvent(func(*PaintEvent) {})
	}
	return this.onPaintEvent.Listen(f)
}
func (this *EventSupport) OnMoved(f func(*MovedEvent)) gxui.EventSubscription {
	if this.onMovedEvent == nil {
		this.onMovedEvent = gxui.CreateEvent(func(*MovedEvent) {})
	}
	return this.onMovedEvent.Listen(f)
}
func (this *EventSupport) OnFocus(f func(*FocusEvent)) gxui.EventSubscription {
	if this.onFocusEvent == nil {
		this.onFocusEvent = gxui.CreateEvent(func(*FocusEvent) {})
	}
	return this.onFocusEvent.Listen(f)
}

func (this *EventSupport) OnDestroy(f func(*DestroyEvent)) gxui.EventSubscription {
	if this.onDestroyEvent == nil {
		this.onDestroyEvent = gxui.CreateEvent(func(*DestroyEvent) {})
	}
	return this.onDestroyEvent.Listen(f)
}

//func (this *EventSupport) OnMouseScroll(f func(*MouseEvent)) gxui.EventSubscription {
//	return this.onMouseScrollEvent.Listen(f)
//}

//func (this *EventSupport) fireDestroyEvent(e DestroyEvent) {
//	if this.onDestroyEvent != nil {
//		this.onDestroyEvent.Fire(e)
//	}
//}

func (this *EventSupport) FireEvent(e event.IEvent) bool {
	switch e.GetType() {
	case MOUSE_MOVE_EVENT_TYPE:
		if this.onMouseMoveEvent != nil {
			this.onMouseMoveEvent.Fire(e)
		}
	case MOUSE_PRESS_EVENT_TYPE:
		if this.onMouseDownEvent != nil {
			this.onMouseDownEvent.Fire(e)
		}
	case MOUSE_RELEASE_EVENT_TYPE:
		if this.onMouseUpEvent != nil {
			this.onMouseUpEvent.Fire(e)
		}
	case MOUSE_ENTER_EVENT_TYPE:
		if this.onMouseEnterEvent != nil {
			this.onMouseEnterEvent.Fire(e)
		}
	case MOUSE_LEAVE_EVENT_TYPE:
		if this.onMouseExitEvent != nil {
			this.onMouseExitEvent.Fire(e)
		}
	case PAINT_EVENT_TYPE:
		if this.onPaintEvent != nil {
			this.onPaintEvent.Fire(e)
		}
	case MODIFIED_EVENT_TYPE:
		if this.onModifiedEvent != nil {
			this.onModifiedEvent.Fire(e)
		}
	case KEY_PRESS_EVENT_TYPE:
		if this.onKeyPressEvent != nil {
			this.onKeyPressEvent.Fire(e)
		}
	case KEY_RELEASE_EVENT_TYPE:
		if this.onKeyReleaseEvent != nil {
			this.onKeyReleaseEvent.Fire(e)
		}
	case KEY_REPEAT_EVENT_TYPE:
		if this.onKeyRepeatEvent != nil {
			this.onKeyRepeatEvent.Fire(e)
		}
	case KEY_CHAR_EVENT_TYPE:
		if this.onKeyCharEvent != nil {
			this.onKeyCharEvent.Fire(e)
		}
	case SELECT_EVENT_TYPE:
		if this.onSelectEvent != nil {
			this.onSelectEvent.Fire(e)
		}
	case MOVED_EVENT_TYPE:
		if this.onMovedEvent != nil {
			this.onMovedEvent.Fire(e)
		}
	case FOCUS_EVENT_TYPE:
		if this.onFocusEvent != nil {
			this.onFocusEvent.Fire(e)
		}
	case DESTROY_EVENT_TYPE:
		if this.onDestroyEvent != nil {
			this.onDestroyEvent.Fire(e)
		}

	}
	return false
}
