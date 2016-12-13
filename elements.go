package eui

import (
	"errors"
	"image"

	"log"
	"sync"

	"github.com/tryor/util/event"
)

var emptyElements = make([]IElement, 0)

type IElements interface {
	//IElement
	Empty() bool
	Size() int
	GetIndex(e IElement) int
	Contains(e IElement) bool
	Exist(id string) (ret bool)

	Add(e IElement, idx ...int) error
	//Adds(es ...IElement) error
	//AddsAndExclude(es []IElement, excluded IElement) error
	Remove(e IElement) bool
	Clear()
	//RemovesAndExcludes(iterator func(e IElement), excludeds IElements) []IElement
	//Removes() []IElement
	//Destroy()

	//	GetElements() []IElement
	CloneElements() []IElement
	At(idx int) IElement
	GetById(id string) IElement
	Sort()
	ClearSortFlag()
	//	GetMouseHoveringElement() IElement
	SetFocusElement(e IElement)
	GetFocusElement() IElement

	ForEach(f func(idx int, el IElement) (continue_ bool))
	ForEachLast(f func(idx int, el IElement) (continue_ bool))

	RedrawAll(excludeds ...IElement)
	RedrawIntersection(excludeds ...IElement)
	SetLayer(layer ILayer)

	CreateBoundRect() *image.Rectangle
}

type Elements struct {
	//*Element
	//*CallSupport

	Self  IElements
	Layer ILayer //元素所在层

	elementsmap    map[string]IElement //Key为元素ID
	elements       []IElement
	elementsLocker sync.RWMutex

	sortFlag             bool     //已经被排序标记
	mouseHoveringElement IElement //鼠标悬停元素
	focusElement         IElement //当前焦点元素

	sortor *ElementSortor
}

func NewElements(els ...IElement) *Elements {
	//es := &Elements{Element: newElement()}
	es := &Elements{}
	es.Self = es
	//es.CallSupport = NewCallSupport()

	es.elements = els //make([]IElement, 0)

	es.elementsmap = make(map[string]IElement)
	es.sortor = NewElementSortor(es.elements, func(elements []IElement, i, j int) bool {
		return elements[i].GetOrderZ() < elements[j].GetOrderZ()
	})
	//	es.pendingCall = make(chan func(), 256)

	//	go es.pendingCallLoop()

	return es
}

//func NewElementsBy(els ...IElement) (*Elements, error) {
//	es := NewElements()
//	//	err := es.Adds(els...)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	return es, nil
//}

//func (this *Elements) Destroy() {
//	this.CallSupport.Destroy()
//}

func (this *Elements) SetFocusElement(e IElement) {
	this.focusElement = e
}

func (this *Elements) GetFocusElement() IElement {
	return this.focusElement
}

//func (this *Elements) Intersects(x, y int) (ret bool) {
//	this.elementsLocker.RLock()
//	defer this.elementsLocker.RUnlock()
//	for _, e := range this.elements { //this.CloneElements() {
//		if e.Intersects(x, y) {
//			return true
//		}
//	}
//	return false
//}

//func (this *Elements) IntersectsWith(rect *image.Rectangle) (ret bool) {
//	//	this.elementsLocker.RLock()
//	//	defer this.elementsLocker.RUnlock()
//	//	this.syncCall(func() {
//	//		for _, e := range this.elements {
//	//			if e.IntersectsWith(rect) {
//	//				ret = true
//	//			}
//	//		}
//	//	})
//	//	return ret

//	for _, e := range this.CloneElements() {
//		if e.IntersectsWith(rect) {
//			return true
//		}
//	}
//	return false
//}

//func (this *Elements) MakeBoundRectSnapshot() {
//	this.elementsLocker.RLock()
//	defer this.elementsLocker.RUnlock()
//	for _, e := range this.elements {
//		e.MakeBoundRectSnapshot()
//	}
//}

func (this *Elements) Empty() (ret bool) {
	//	this.elementsLocker.RLock()
	//	defer this.elementsLocker.RUnlock()
	//	if this.IsCalling() {
	//		this.syncCall(func() {
	//			ret = len(this.elements) == 0
	//		})
	//	} else {
	//		ret = len(this.elements) == 0
	//	}
	//	return
	return this.Size() == 0
}

func (this *Elements) Size() (ret int) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//	if this.IsCalling() {
	//		this.syncCall(func() {
	//			ret = len(this.elements)
	//		})
	//	} else {
	//		ret = len(this.elements)
	//	}

	return len(this.elements)
}

func (this *Elements) Exist(id string) (ret bool) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//	if this.IsCalling() {
	//		this.syncCall(func() {
	//			_, ret = this.elementsmap[id]
	//		})
	//	} else {
	//		_, ret = this.elementsmap[id]
	//	}
	//println("Elements) Exist:", ret)
	_, ret = this.elementsmap[id]
	return
}

func (this *Elements) Add(e IElement, idx ...int) error {
	if e.IsDestroyed() {
		return errors.New("the element id is destroyed")
	}

	id := e.GetId()
	if id == "" {
		return errors.New("the element id is nil")
	}

	//	this.elementsLocker.Lock()
	//	defer this.elementsLocker.Unlock()

	//if ee, ok := this.elementsmap[e.GetId()]; ok {
	ee := this.GetById(id)
	if ee != nil && ee != e {
		return errors.New("the element already exists, id is" + e.GetId())
	}

	if len(idx) > 0 {
		e.SetOrderZ(idx[0])
		this.sortFlag = false
	}

	//	if this.contains(e) {
	//		return nil
	//	}
	//if this.Contains(e) {
	if ee != nil {
		return nil
	}

	//	this.syncCall(func() {
	//		this.elements = append(this.elements, e)
	//		this.elementsmap[e.GetId()] = e
	//		this.sortFlag = false
	//	})

	this.elementsLocker.Lock()
	defer this.elementsLocker.Unlock()
	this.elements = append(this.elements, e)
	this.elementsmap[e.GetId()] = e
	this.sortFlag = false

	return nil
}

//func (this *Elements) Add_old(e IElement, idx ...int) error {
//	//	if e == nil {
//	//		return false
//	//	}

//	if e.GetId() == "" {
//		return errors.New("the element id is nil")
//	}

//	this.elementsLocker.Lock()
//	defer this.elementsLocker.Unlock()

//	if ee, ok := this.elementsmap[e.GetId()]; ok {
//		if ee != e {
//			return errors.New("the element already exists, id is" + e.GetId())
//		}
//	}

//	elsSize := len(this.elements)
//	pos := elsSize
//	if len(idx) > 0 {
//		pos = idx[0]
//	}

//	if pos < 0 {
//		pos = 0
//	} else if pos > elsSize {
//		pos = elsSize
//	}

//	i := this.getIndex(e)
//	if i > -1 {
//		if pos < elsSize && this.elements[i] == this.elements[pos] {
//			//return false
//			return nil
//		}
//		this.elements = removeElement(this.elements, i, i+1)
//		if i < pos {
//			pos--
//		}
//		this.elements = insertElement(this.elements, pos, e)

//	} else {
//		if pos >= elsSize {
//			this.elements = append(this.elements, e)
//		} else {
//			this.elements = insertElement(this.elements, pos, e)
//		}
//		this.elementsmap[e.GetId()] = e
//	}
//	this.sortFlag = false

//	return nil
//}

/**
 * 返回元素索引位置, 如果元素不存在，返回-1
 */
func (this *Elements) GetIndex(e IElement) (ret int) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	return this.getIndex(e)
	//	this.syncCall(func() {
	//		ret = this.getIndex(e)
	//	})
	//	return
}

//func (this *Elements) GetChildIndex(e IElement) int {
//	return this.Self.(IElements).GetIndex(e)
//}

/**
 * 返回元素索引位置, 如果元素不存在，返回-1
 */
func (this *Elements) getIndex(e IElement) int {
	//	if this.sortFlag {
	//		idx := sort.Search(len(this.elements), func(i int) bool {
	//			return this.elements[i].GetOrderZ() == e.GetOrderZ() && this.elements[i] == e
	//		})
	//		if idx >= len(this.elements) {
	//			idx = -1
	//		}
	//		return idx
	//	} else {

	for i, v := range this.elements {
		if v == e {
			return i
		}
	}
	return -1

	//	return searchElement(this.elements, e)
	//	}
}

//func (this *Elements) AddChildren(els []IElement, pos ...int) error {
//	for _, e := range els {
//		err := this.Self.(IElements).AddChild(e, pos...)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func (this *Elements) Adds(es ...IElement) error {
//	return this.Self.(IElements).AddsAndExclude(es, nil)
//}

//func (this *Elements) AddsAndExclude(els []IElement, excluded IElement) error {
//	//	if len(els) == 0 {
//	//		return false
//	//	}
//	//	this.elementsLocker.Lock()
//	//	defer this.elementsLocker.Unlock()

//	//	retrs := false

//	for _, e := range els {
//		if e == excluded {
//			continue
//		}
//		if e.GetId() == "" {
//			return errors.New("the element id is nil")
//		}
//		if ee, ok := this.elementsmap[e.GetId()]; ok {
//			if ee != e {
//				return errors.New("the element already exists, id is" + e.GetId())
//			}
//		} else {
//			this.elements = append(this.elements, e)
//			this.elementsmap[e.GetId()] = e
//		}

//		//		if !this.contains(e) {
//		//			this.elements = append(this.elements, e)
//		//			this.sortFlag = false
//		//			if !retrs {
//		//				retrs = true
//		//			}
//		//		}
//	}
//	this.sortFlag = false
//	return nil
//}

//func (this *Elements) HasChild(e IElement) bool {
//	return this.Self.(IElements).Contains(e)
//}

func (this *Elements) Contains(e IElement) (ret bool) {
	//	this.elementsLocker.RLock()
	//	defer this.elementsLocker.RUnlock()
	//	this.syncCall(func() {
	//		ret = this.getIndex(e) > -1
	//	})
	//	return
	//	return this.getIndex(e) > -1
	el := this.GetById(e.GetId())
	if el == nil {
		return false
	}
	if el != e {
		log.Println("Elements.Contains, el != e")
	}
	return true
}

//func (this *Elements) ContainsBy(id string) bool {
//	this.elementsLocker.RLock()
//	defer this.elementsLocker.RUnlock()
//	_, ok := this.elementsmap[id]
//	return ok
//}

func (this *Elements) contains(e IElement) bool {
	//return this.getIndex(e) > -1

	el := this.GetById(e.GetId())
	if el == nil {
		return false
	}
	if el != e {
		log.Println("Elements.Contains, el != e")
	}
	return true

}

//func (this *Elements) RemoveChild(e IElement) bool {
//	return this.Self.(IElements).Remove(e)
//}

func (this *Elements) Remove(e IElement) (ret bool) {
	if e == nil {
		return false
	}
	this.elementsLocker.Lock()
	defer this.elementsLocker.Unlock()

	//	this.syncCall(func() {
	//		idx := this.getIndex(e)
	//		if idx > -1 {
	//			ret = this.remove(idx)
	//		} else {
	//			ret = false
	//		}
	//	})
	//	return
	idx := this.getIndex(e)
	if idx > -1 {
		return this.remove(idx)
	}
	return false
}

func (this *Elements) remove(idx int) bool {
	if idx < len(this.elements) && idx >= 0 {
		re := this.elements[idx]

		if re == this.mouseHoveringElement {
			this.mouseHoveringElement = nil
		}

		if re == this.focusElement {
			this.focusElement = nil
		}

		delete(this.elementsmap, re.GetId())
		this.elements = removeElement(this.elements, idx, idx+1)
		this.sortFlag = false
		return true
	}
	return false
}

//func (this *Elements) ClearChildren() {
//	this.Self.(IElements).Clear()
//}

func (this *Elements) Clear() {
	this.elementsLocker.Lock()
	defer this.elementsLocker.Unlock()
	//	this.syncCall(func() {
	this.mouseHoveringElement = nil
	this.focusElement = nil
	this.elements = this.elements[0:0]
	this.elementsmap = make(map[string]IElement)
	this.sortFlag = false
	//	})

}

//func (this *Elements) Removes() []IElement {
//	return this.Self.(IElements).RemovesAndExcludes(nil, nil)
//}

//每删除一个元素都会调用iterator函数，iterator可以为nil
//移除参数<code>excludeds</code>元素以外的元素, excludeds指示不被移除的元素
//返回被移除的元素
//func (this *Elements) RemovesAndExcludes(iterator func(e IElement), excludeds IElements) []IElement {
//	res := make([]IElement, 0)
//	if this.Empty() {
//		return res
//	}
//	//	this.elementsLocker.Lock()
//	//	defer this.elementsLocker.Unlock()
//	els := make([]IElement, 0)
//	count := len(this.elements)
//	for _, el := range this.elements {
//		if excludeds == nil || !excludeds.Contains(el) {
//			res = append(res, el)
//			delete(this.elementsmap, el.GetId())

//			if el == this.mouseHoveringElement {
//				this.mouseHoveringElement = nil
//			}
//			if el == this.focusElement {
//				this.focusElement = nil
//			}
//			if iterator != nil {
//				iterator(el)
//			}
//		} else {
//			els = append(els, el)
//		}
//	}
//	this.elements = els
//	if len(this.elements) != count {
//		this.sortFlag = false
//	}
//	return res
//}

//func (this *Elements) CreateBoundRect() *image.Rectangle {
//	//	this.locker.RLock()
//	//	defer this.locker.RUnlock()
//	for _, el := range this.GetElements() {
//		el.CreateBoundRect()
//	}
//	return this.Element.CreateBoundRect()
//}

/**
 * 创建边界矩形
 */
func (this *Elements) CreateBoundRect() *image.Rectangle {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()

	var l, t, r, b int   //左，上，右，下
	els := this.elements //this.GetElements()
	if len(els) > 0 {
		firstRect := els[0].GetBoundRect()
		l = firstRect.Min.X
		t = firstRect.Min.Y
		r = firstRect.Max.X
		b = firstRect.Max.Y
		for _, el := range els[0:] {
			rect := el.GetBoundRect()
			if rect.Min.X < l {
				l = rect.Min.X
			}
			if rect.Min.Y < t {
				t = rect.Min.Y
			}
			if rect.Max.X > r {
				r = rect.Max.X
			}
			if rect.Max.Y > b {
				b = rect.Max.Y
			}
		}
	}

	return &image.Rectangle{Min: image.Point{l, t}, Max: image.Point{r, b}}

	//	this.Self.(IElement).SetCoordinate(l, t)
	//	this.Self.(IElement).SetWidth(r - l)
	//	this.Self.(IElement).SetHeight(b - t)
	//	return this.Element.CreateBoundRect()

}

func (this *Elements) SetLayer(layer ILayer) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	for _, e := range this.elements { //this.CloneElements() {
		e.SetLayer(layer)
	}
	this.Layer = layer
}

func (this *Elements) IsModified() bool {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	for _, e := range this.elements { //this.CloneElements() {
		if e.IsModified() {
			return true
		}
	}
	return false
}

//从第一个元素开始
func (this *Elements) ForEach(f func(idx int, el IElement) (continue_ bool)) {
	//	this.elementsLocker.RLock()
	//	defer this.elementsLocker.RUnlock()
	els := this.CloneElements()
	//	els := this.Clone()
	for i, el := range els {
		if !f(i, el) {
			break
		}
	}
}

//从最后面元素开始
func (this *Elements) ForEachLast(f func(idx int, el IElement) (continue_ bool)) {
	//	this.elementsLocker.RLock()
	//	defer this.elementsLocker.RUnlock()
	els := this.CloneElements()
	//	els := this.Clone()
	for i := len(els) - 1; i >= 0; i-- {
		if !f(i, els[i]) {
			break
		}
	}
}

//func (this *Elements) GetElements() []IElement {
//	this.elementsLocker.RLock()
//	defer this.elementsLocker.RUnlock()
//	//	result := make([]IElement, len(this.elements))
//	//	copy(result, this.elements)
//	return this.elements
//}

func (this *Elements) CloneElements() []IElement {
	if this.Size() == 0 {
		return emptyElements
	}
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//return this.elements

	result := make([]IElement, len(this.elements))
	copy(result, this.elements)
	//	if len(this.elements) > 10 {
	//		println("Elements.CloneElements:", len(this.elements))
	//	}
	return result
}

//func (this *Elements) GetChildren() []IElement {
//	return this.Self.(IElements).GetElements()
//}

func (this *Elements) At(idx int) IElement {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	if idx < len(this.elements) && idx >= 0 {
		return this.elements[idx]
	}
	return nil
}

func (this *Elements) GetById(id string) IElement {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	if e, ok := this.elementsmap[id]; ok {
		return e
	} else {
		return nil
	}
}

//func (this *Elements) DrawChildren(ge IGraphicsEngine) {
//	this.Self.(IElements).Draw(ge)
//}

func (this *Elements) ClearSortFlag() {
	this.sortFlag = false
}

/**
 * 构建路径
 */
func (this *Elements) CreatePath() {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	for _, el := range this.elements { //this.CloneElements() {
		el.CreatePath()
	}
}

func (this *Elements) Draw(ge IGraphicsEngine) {
	this.Sort()
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	for _, el := range this.elements { //this.SortAndCloneElements() {
		if el.IsVisible() {
			el.SetModified(false)
			el.SetRedraw(false)
			el.Draw(ge)
			el.firePaintEvent()
		}
	}
}

/**
 * 原地重画与此组里所有元素相交的其它元素, 自己不重画，排除<code>excludeds</code>中的元素
 */
func (this *Elements) RedrawIntersection(excludeds ...IElement) {

	if !(this.Layer != nil && this.Layer.GetDrawMode() == DrawMode_Region && !this.Empty()) {
		return
	}
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	if len(excludeds) == 0 {
		this.redrawAll(this.elements...)
		//		this.Self.(IElements).RedrawAll(this.CloneElements()...)
	} else {
		//		allExcludeds := insertElement(this.CloneElements(), this.Size(), excludeds...)
		//this.Self.(IElements).RedrawAll(allExcludeds...)
		allExcludeds := insertElement(this.elements, len(this.elements), excludeds...)
		this.redrawAll(allExcludeds...)
	}

}

func (this *Elements) RedrawAll(excludeds ...IElement) {
	if !(this.Layer != nil && this.Layer.GetDrawMode() == DrawMode_Region && !this.Empty()) {
		return
	}
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	this.redrawAll(excludeds...)
}

/**
 * 原地重画此组里所有元素, 包括相交元素和自己, 排除<code>excludeds</code>中的元素
 */
func (this *Elements) redrawAll(excludeds ...IElement) {
	//if this.Layer != nil && this.Layer.GetDrawMode() == DrawMode_Region && !this.Empty() {
	//		allExcludeds, err := NewElements(excludeds...)
	//		if err != nil {
	//			return
	//		}
	//allExcludeds := NewElements(excludeds...)
	//allExcludeds.Sort()
	sortor := NewElementSortor(excludeds, nil)
	sortor.Sort()
	//		els := this.GetElements()

	//		this.elementsLocker.RLock()
	//		defer this.elementsLocker.RUnlock()

	for _, el := range this.elements { //this.CloneElements() {
		//清除绘图区
		this.Layer.UnionDrawRegion(el.GetBoundRect(el.GetClipRegionAdjustValue()))
		iels := this.Layer.GetIntersectionsByElement(el, sortor, nil)
		if len(iels) > 0 {
			//通知绘制相交元素,
			for _, iel := range iels {
				iel.SetRedraw(true)
			}
			//allExcludeds.Adds(iels...)
			//allExcludeds.Sort()
			excludeds = append(excludeds, iels...)
			sortor.SetElements(excludeds)
			sortor.Sort()
		}
	}
	//}

}

//func (this *Elements) FireEvent(e event.IEvent, asyn ...bool) bool {
//	ret := false
//	for _, el := range this.elements {
//		if el.FireEvent(e, asyn...) && !ret {
//			ret = true
//		}
//	}
//	return ret
//}

/**
 * 增量移动到目标位置
 */
//func (this *Elements) MoveBy(dx, dy int, angle float64) {
//	if this.Empty() || this.Layer == nil {
//		return
//	}
//	//绘制与此组中元素相交的其它元素, 自己不重画
//	this.Self.(IElements).RedrawIntersection()
//	//在目标区绘图
//	for _, e := range this.GetElements() {
//		e.MoveBy(dx, dy, angle)
//	}
//	//	this.Self.(IElements).ClearBoundRect()
//	//	this.Self.(IElements).CreateBoundRect()

//	this.Element.MoveBy(dx, dy, angle)
//}

func (this *Elements) SetMoving(b bool) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//log.Println("Elements) SetMoving:", b)
	for _, e := range this.elements { //this.CloneElements() {
		e.SetMoving(b)
	}
	//	this.MoveSupport.SetMoving(b)
}

func (this *Elements) PrepareTransform(x, y int) {
	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//log.Println("Elements) PrepareTransform:", x, y)
	for _, e := range this.elements { //this.CloneElements() {
		e.PrepareTransform(x, y)
	}
	//	this.MoveSupport.PrepareTransform(x, y)
}

/**
 * 对元素进行排序
 */
func (this *Elements) Sort() {
	if this.sortFlag {
		return
	}
	this.elementsLocker.Lock()
	defer this.elementsLocker.Unlock()

	this.sortor.SetElements(this.elements)
	this.sortor.Sort()
	this.sortFlag = true
}

func (this *Elements) SortAndCloneElements() []IElement {
	if this.sortFlag {
		return this.CloneElements()
	}
	this.elementsLocker.Lock()
	defer this.elementsLocker.Unlock()

	this.sortor.SetElements(this.elements)
	this.sortor.Sort()
	this.sortFlag = true

	result := make([]IElement, len(this.elements))
	copy(result, this.elements)
	return result
}

//func (this *Elements) Len() int {
//	return len(this.elements)
//}

//func (this *Elements) Swap(i, j int) {
//	this.elements[i], this.elements[j] = this.elements[j], this.elements[i]
//}

//func (this *Elements) Less(i, j int) bool {
//	//	return reflect.ValueOf(this.elements[i]).Pointer() < reflect.ValueOf(this.elements[j]).Pointer()
//	//	return uintptr(unsafe.Pointer(&this.elements[i])) < uintptr(unsafe.Pointer(&this.elements[j]))
//	//	return uintptr(this.elements[i]) < uintptr(this.elements[j])
//	return this.elements[i].GetOrderZ() < this.elements[j].GetOrderZ()
//}

func (this *Elements) GetMouseHoveringElement() IElement {
	return this.mouseHoveringElement
}

//func (this *Elements) setMouseHoveringElement(e IElement) {
//	this.mouseHoveringElement = e
//}

//func (this *Elements) ClearMouseHoveringElement() {
//	this.mouseHoveringElement = nil
//}

//func (this *Elements) clearDifferentMouseHoveringElement(el IElement) {
//	if this.mouseHoveringElement != el {
//		this.mouseHoveringElement = nil
//	}
//}

/**
 * 跟踪响应事件
 */
func (this *Elements) TrackEvent(e event.IEvent) bool {
	//	if !this.eventEnabled {
	//		return false
	//	}

	focusElement := this.focusElement
	if IsMouseEvent(e.GetType()) {
		//如果是鼠标事件
		ret := this.handleMouseEvent(e.(IMouseEvent))
		if e.GetType() == MOUSE_PRESS_EVENT_TYPE {
			mhoveringElement := this.GetMouseHoveringElement()
			if mhoveringElement != nil {
				if focusElement != mhoveringElement {
					if focusElement != nil {
						focusElement.fireFocusEvent(false)
					}
					focusElement = mhoveringElement
					this.focusElement = focusElement //this.MouseHoveringElement
					focusElement.fireFocusEvent(true)
				}
			} else if focusElement != nil {
				focusElement.fireFocusEvent(false)
				this.focusElement = nil
			}
		}
		return ret
	} else {
		//		if this.Self.(IElements).FireEvent(e) {
		//			return true
		//		}
		if focusElement != nil && focusElement.IsEventEnabled() {
			return focusElement.TrackEvent(e)
		}
		return false
	}
}

/**
 * 处理鼠标事件
 */
func (this *Elements) handleMouseEvent(me IMouseEvent) bool {
	x := me.X()
	y := me.Y()
	//向组中元素转发鼠标事件
	var manualModeElement IElement
	var e IElement
	var els []IElement
	//有可能没有层,那默认使用MEventMode_Hovering模式
	var mode MouseEventMode
	if this.Layer != nil {
		mode = this.Layer.GetMouseEventMode()
	} else {
		mode = MEventMode_Hovering
	}
	//	println("Elements.handleMouseEvent.Layer:", this.Layer, mode)
	switch mode {
	case MEventMode_Hovering:
		mhoveringElement := this.GetMouseHoveringElement()
		e = this.findByEventEnabled(x, y, nil)
		if mhoveringElement != nil {
			if mhoveringElement != e {
				mhoveringElement.fireMouseLeaveEvent(me)
				this.mouseHoveringElement = nil
			}
		}
		if e != nil {
			if mhoveringElement != e {
				this.mouseHoveringElement = e
				mhoveringElement = e
				mhoveringElement.fireMouseEnterEvent(me)
			}
			mhoveringElement.TrackEvent(me)
			//			if this.MouseHoveringElement != nil {
			//				return this.MouseHoveringElement.TrackEvent(me)
			//			} else {
			//				println("Elements.handleMouseEvent MouseHoveringElement is nil")
			//			}
		}
	case MEventMode_Top:
		els = this.elements //getAbstractElements();
		if len(els) > 0 {
			mhoveringElement := this.GetMouseHoveringElement()
			e = els[len(els)-1]
			if e.IsEventEnabled() && e.IsVisible() && e.Intersects(x, y) {
				if mhoveringElement != e {
					mhoveringElement = e
					this.mouseHoveringElement = e
					mhoveringElement.fireMouseEnterEvent(me)
				}
				return e.TrackEvent(me)
			} else if mhoveringElement != nil {
				//mouseLeaveElement(MouseHoveringElement)
				mhoveringElement.fireMouseLeaveEvent(me)
				this.mouseHoveringElement = nil
			}
		}
	case MEventMode_Manual:
		mhoveringElement := this.GetMouseHoveringElement()
		manualModeElement = this.Layer.GetManualModeElement()
		if manualModeElement != nil && manualModeElement.IsEventEnabled() &&
			manualModeElement.IsVisible() &&
			manualModeElement.Intersects(x, y) {
			if mhoveringElement != manualModeElement {
				mhoveringElement = manualModeElement
				this.mouseHoveringElement = manualModeElement
				mhoveringElement.fireMouseEnterEvent(me)
			}
			return mhoveringElement.TrackEvent(me)
		} else {
			if mhoveringElement != nil {
				mhoveringElement.fireMouseLeaveEvent(me)
				this.mouseHoveringElement = nil
			}
		}
	}

	return false
}

/**
 * 返回指定坐标 x,y上的最顶层可视并能响应事件的元素, 如果没有，返回NULL
 *
 * @param x X坐标
 * @param y Y坐标
 * @param excluded 被排除的元素
 *
 */
func (this *Elements) findByEventEnabled(x, y int, excluded IElement) IElement {
	this.Sort()

	this.elementsLocker.RLock()
	defer this.elementsLocker.RUnlock()
	//	for _, e := range this.elements {
	//		if e != excluded && e.IsEventEnabled() && e.IsVisible() && e.Intersects(x, y) {
	//			return e
	//		}
	//	}
	//els := this.SortAndCloneElements() //this.elements //this.elements.GetElements()
	els := this.elements
	for i := len(els) - 1; i >= 0; i-- {
		e := els[i]
		if e != excluded && e.IsEventEnabled() && e.IsVisible() && e.Intersects(x, y) {
			return e
		}
	}

	return nil
}

func (this *Elements) fireMouseLeaveEvent(e IMouseEvent) bool {
	//如果鼠标最近一次悬停的元素存在，发送离开元素事件
	mhoveringElement := this.GetMouseHoveringElement()
	if mhoveringElement != nil && mhoveringElement.MouseIsHovering() {
		mhoveringElement.fireMouseLeaveEvent(e)
	}
	this.mouseHoveringElement = nil

	//	return this.Element.fireMouseLeaveEvent(e)
	return false
}

func (this *Elements) fireFocusEvent(focus bool) bool {
	focusElement := this.focusElement
	if focusElement != nil {
		if focus {
			return focusElement.fireFocusEvent(focus)
		} else {
			//			fel := this.FocusElement
			this.focusElement = nil
			return focusElement.fireFocusEvent(focus)
		}
	}
	//	return this.Element.fireFocusEvent(focus)
	return false
}

//func searchElement(es []IElement, e IElement) int {
//	for i, v := range es {
//		if v == e {
//			return i
//		}
//	}
//	return -1
//}

func insertElement(slice []IElement, index int, insertion ...IElement) []IElement {
	result := make([]IElement, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

func removeElement(slice []IElement, start, end int) []IElement {
	return append(slice[:start], slice[end:]...)
}
