package eui

import (
	"github.com/tryor/commons/event"
)

type IEventSupport interface {
	FireEvent(e event.IEvent) bool

	OnMouseMove(func(*MouseEvent)) EventSubscription
	OnMouseDown(func(*MouseEvent)) EventSubscription
	OnMouseUp(func(*MouseEvent)) EventSubscription
	OnMouseDoubleClick(func(*MouseEvent)) EventSubscription
	OnMouseEnter(func(*MouseEvent)) EventSubscription
	OnMouseExit(func(*MouseEvent)) EventSubscription
	OnMouseWheel(func(*WheelEvent)) EventSubscription
	OnMouseOutsideDown(func(*MouseEvent)) EventSubscription

	OnKeyPress(func(*KeyEvent)) EventSubscription
	OnKeyRelease(func(*KeyEvent)) EventSubscription
	OnKeyRepeat(func(*KeyEvent)) EventSubscription
	OnKeyChar(func(*KeyEvent)) EventSubscription

	OnSelect(func(*SelectEvent)) EventSubscription
	OnModified(func(*ModifiedEvent)) EventSubscription
	OnPaint(func(*PaintEvent)) EventSubscription
	OnMoved(func(*MovedEvent)) EventSubscription
	OnFocus(func(*FocusEvent)) EventSubscription
	OnVisible(func(*VisibleEvent)) EventSubscription

	OnDestroy(func(*DestroyEvent)) EventSubscription

	//OnMouseScroll(func(*MouseEvent)) EventSubscription
}

type EventSupport struct {
	//	Self interface{}

	onMouseMoveEvent   Event
	onMouseDownEvent   Event
	onMouseUpEvent     Event
	onMouseDoubleClick Event
	onMouseEnterEvent  Event
	onMouseExitEvent   Event
	onMouseWheelEvent  Event
	//	onMouseScrollEvent Event
	onMouseOutsideDownEvent Event

	onKeyPressEvent   Event
	onKeyReleaseEvent Event
	onKeyRepeatEvent  Event
	onKeyCharEvent    Event
	onSelectEvent     Event
	onModifiedEvent   Event
	onPaintEvent      Event
	onMovedEvent      Event
	onFocusEvent      Event
	onDestroyEvent    Event
	onVisibleEvent    Event
}

//func initEvent(esupport *EventSupport) {
//	esupport.onMouseMoveEvent = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseDownEvent = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseUpEvent = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseDoubleClick = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseEnterEvent = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseExitEvent = CreateEvent(func(*MouseEvent) {})
//	//esupport.onMouseScrollEvent = CreateEvent(func(*MouseEvent) {})
//	esupport.onMouseWheelEvent = CreateEvent(func(*WheelEvent) {})

//	esupport.onKeyPressEvent = CreateEvent(func(*KeyEvent) {})
//	esupport.onKeyReleaseEvent = CreateEvent(func(*KeyEvent) {})
//	esupport.onKeyCharEvent = CreateEvent(func(*KeyEvent) {})
//	esupport.onSelectEvent = CreateEvent(func(*SelectEvent) {})
//	esupport.onModifiedEvent = CreateEvent(func(*ModifiedEvent) {})
//	esupport.onPaintEvent = CreateEvent(func(*PaintEvent) {})
//	esupport.onMovedEvent = CreateEvent(func(*MovedEvent) {})
//	esupport.onFocusEvent = CreateEvent(func(*FocusEvent) {})
//}

func (this *EventSupport) OnMouseMove(f func(*MouseEvent)) EventSubscription {
	if this.onMouseMoveEvent == nil {
		this.onMouseMoveEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseMoveEvent.Listen(f)
}

func (this *EventSupport) OnMouseDown(f func(*MouseEvent)) EventSubscription {
	if this.onMouseDownEvent == nil {
		this.onMouseDownEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseDownEvent.Listen(f)
}

func (this *EventSupport) OnMouseUp(f func(*MouseEvent)) EventSubscription {
	if this.onMouseUpEvent == nil {
		this.onMouseUpEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseUpEvent.Listen(f)
}

func (this *EventSupport) OnMouseDoubleClick(f func(*MouseEvent)) EventSubscription {
	if this.onMouseDoubleClick == nil {
		this.onMouseDoubleClick = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseDoubleClick.Listen(f)
}

func (this *EventSupport) OnMouseEnter(f func(*MouseEvent)) EventSubscription {
	if this.onMouseEnterEvent == nil {
		this.onMouseEnterEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseEnterEvent.Listen(f)
}

func (this *EventSupport) OnMouseExit(f func(*MouseEvent)) EventSubscription {
	if this.onMouseExitEvent == nil {
		this.onMouseExitEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseExitEvent.Listen(f)
}

func (this *EventSupport) OnMouseWheel(f func(*WheelEvent)) EventSubscription {
	if this.onMouseWheelEvent == nil {
		this.onMouseWheelEvent = CreateEvent(func(*WheelEvent) {})
	}
	return this.onMouseWheelEvent.Listen(f)
}

func (this *EventSupport) OnMouseOutsideDown(f func(*MouseEvent)) EventSubscription {
	if this.onMouseOutsideDownEvent == nil {
		this.onMouseOutsideDownEvent = CreateEvent(func(*MouseEvent) {})
	}
	return this.onMouseOutsideDownEvent.Listen(f)
}

func (this *EventSupport) OnKeyPress(f func(*KeyEvent)) EventSubscription {
	if this.onKeyPressEvent == nil {
		this.onKeyPressEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyPressEvent.Listen(f)
}

func (this *EventSupport) OnKeyRelease(f func(*KeyEvent)) EventSubscription {
	if this.onKeyReleaseEvent == nil {
		this.onKeyReleaseEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyReleaseEvent.Listen(f)
}

func (this *EventSupport) OnKeyRepeat(f func(*KeyEvent)) EventSubscription {
	if this.onKeyRepeatEvent == nil {
		this.onKeyRepeatEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyRepeatEvent.Listen(f)
}

func (this *EventSupport) OnKeyChar(f func(*KeyEvent)) EventSubscription {
	if this.onKeyCharEvent == nil {
		this.onKeyCharEvent = CreateEvent(func(*KeyEvent) {})
	}
	return this.onKeyCharEvent.Listen(f)
}
func (this *EventSupport) OnSelect(f func(*SelectEvent)) EventSubscription {
	if this.onSelectEvent == nil {
		this.onSelectEvent = CreateEvent(func(*SelectEvent) {})
	}
	return this.onSelectEvent.Listen(f)
}
func (this *EventSupport) OnModified(f func(*ModifiedEvent)) EventSubscription {
	if this.onModifiedEvent == nil {
		this.onModifiedEvent = CreateEvent(func(*ModifiedEvent) {})
	}
	return this.onModifiedEvent.Listen(f)
}
func (this *EventSupport) OnPaint(f func(*PaintEvent)) EventSubscription {
	if this.onPaintEvent == nil {
		this.onPaintEvent = CreateEvent(func(*PaintEvent) {})
	}
	return this.onPaintEvent.Listen(f)
}
func (this *EventSupport) OnMoved(f func(*MovedEvent)) EventSubscription {
	if this.onMovedEvent == nil {
		this.onMovedEvent = CreateEvent(func(*MovedEvent) {})
	}
	return this.onMovedEvent.Listen(f)
}
func (this *EventSupport) OnFocus(f func(*FocusEvent)) EventSubscription {
	if this.onFocusEvent == nil {
		this.onFocusEvent = CreateEvent(func(*FocusEvent) {})
	}
	return this.onFocusEvent.Listen(f)
}

func (this *EventSupport) OnDestroy(f func(*DestroyEvent)) EventSubscription {
	if this.onDestroyEvent == nil {
		this.onDestroyEvent = CreateEvent(func(*DestroyEvent) {})
	}
	return this.onDestroyEvent.Listen(f)
}

func (this *EventSupport) OnVisible(f func(*VisibleEvent)) EventSubscription {
	if this.onVisibleEvent == nil {
		this.onVisibleEvent = CreateEvent(func(*VisibleEvent) {})
	}
	return this.onVisibleEvent.Listen(f)
}

//OnVisible(func(*VisibleEvent)) EventSubscription

//func (this *EventSupport) OnMouseScroll(f func(*MouseEvent)) EventSubscription {
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
	case MOUSE_DOUBLE_CLICK_EVENT_TYPE:
		if this.onMouseDoubleClick != nil {
			this.onMouseDoubleClick.Fire(e)
		}

	case MOUSE_ENTER_EVENT_TYPE:
		if this.onMouseEnterEvent != nil {
			this.onMouseEnterEvent.Fire(e)
		}
	case MOUSE_LEAVE_EVENT_TYPE:
		if this.onMouseExitEvent != nil {
			this.onMouseExitEvent.Fire(e)
		}
	case MOUSE_WHEEL_EVENT_TYPE:
		if this.onMouseWheelEvent != nil {
			this.onMouseWheelEvent.Fire(e)
		}

	case MOUSE_OUTSIDE_EVENT_TYPE:
		if this.onMouseOutsideDownEvent != nil {
			this.onMouseOutsideDownEvent.Fire(e)
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
	case VISIBLE_EVENT_TYPE:
		if this.onVisibleEvent != nil {
			this.onVisibleEvent.Fire(e)
		}

	}
	return false
}
