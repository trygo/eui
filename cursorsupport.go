package eui

type ICursorSupport interface {
	GetCursor() ICursor
	SetCursor(c ICursor)
	ResetCursor()
}

type CursorSupport struct {
	cursor ICursor
}

func (this *CursorSupport) GetCursor() ICursor {
	return this.cursor
}

func (this *CursorSupport) ResetCursor() {
	this.cursor = nil
}

func (this *CursorSupport) SetCursor(c ICursor) {
	this.cursor = c
}
