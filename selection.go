package eui

type ISelection interface {
	IMoveSupport
	IElements
	//	event.IListener
}
type Selection struct {
	MoveSupport
	*Elements
}

func NewSelection() *Selection {
	s := &Selection{Elements: NewElements()}
	s.Self = s
	return s
}

func (this *Selection) SetMoving(b bool) {
	this.Elements.SetMoving(b)
	this.MoveSupport.SetMoving(b)
}

func (this *Selection) PrepareTransform(x, y int) {
	this.Elements.PrepareTransform(x, y)
	this.MoveSupport.PrepareTransform(x, y)
}

/**
 * 增量移动到目标位置
 */
func (this *Selection) MoveBy(dx, dy int, angle float32) {
	if this.Empty() || this.Layer == nil {
		return
	}
	//绘制与此组中元素相交的其它元素, 自己不重画
	this.Self.(IElements).RedrawIntersection()
	//在目标区绘图
	this.ForEach(func(i int, el IElement) bool {
		el.MoveBy(dx, dy, angle)
		return true
	})
}

func (this *Selection) MoveTo(x, y int, angle float32) {
	mover := this.Self.(IMoveSupport)
	mover.MoveBy(x-int(mover.ReferencePointX()), y-int(mover.ReferencePointY()), angle)
	mover.PrepareTransform(x, y)
}

func (this *Selection) Add(e IElement, idx ...int) error {
	if !e.IsSelected() {
		e.SetSelected(true)
		e.FireEvent(NewSelectEvent(e, true))
	}
	return this.Elements.Add(e, idx...)
}

func (this *Selection) Remove(el IElement) bool {
	if el == nil {
		return false
	}
	if el.IsSelected() {
		el.SetSelected(false)
		el.FireEvent(NewSelectEvent(el, false))
	}
	return this.Elements.Remove(el)
}

func (this *Selection) Clear() {
	for _, el := range this.elements {
		if el.IsSelected() {
			el.SetSelected(false)
			el.FireEvent(NewSelectEvent(el, false))
		}
	}
	this.Elements.Clear()
}
