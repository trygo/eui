package eui

import (
//	"github.com/trygo/util/event"
)

type ISelection interface {
	IElements
	//	event.IListener
}
type Selection struct {
	*Elements
}

func NewSelection() *Selection {
	s := &Selection{Elements: NewElements()}
	s.Self = s
	return s
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
	//	for _, e := range this.GetElements() {
	//		e.MoveBy(dx, dy, angle)
	//	}
	this.ForEach(func(i int, el IElement) bool {
		el.MoveBy(dx, dy, angle)
		return true
	})
}

func (this *Selection) Add(e IElement, idx ...int) error {
	//	if e == nil {
	//		return false
	//	}
	if !e.IsSelected() {
		e.SetSelected(true)
		e.FireEvent(NewSelectEvent(e, true))
	}
	return this.Elements.Add(e, idx...)
}

//func (this *Selection) AddsAndExclude(els []IElement, excluded IElement) error {
//	for _, el := range els {
//		if !el.IsSelected() {
//			el.SetSelected(true)
//			el.FireEvent(NewSelectEvent(el, true))
//		}
//	}
//	return this.Elements.AddsAndExclude(els, excluded)
//}

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

//func (this *Selection) Clear() {
//	this.Self.(ISelection).RemovesAndExcludes(nil, nil)
//}

//func (this *Selection) Removes() []IElement {
//	return this.Self.(ISelection).RemovesAndExcludes(nil, nil)
//}

//func (this *Selection) RemovesAndExcludes(iterator func(e IElement), excludeds IElements) []IElement {
//	removeds := this.Elements.RemovesAndExcludes(iterator, excludeds)
//	for _, el := range removeds {
//		el.SetSelected(false)
//		el.FireEvent(NewSelectEvent(el, false))
//	}
//	return removeds
//}
