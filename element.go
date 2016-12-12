package eui

import (
	"fmt"
	"image"
	"log"
	"reflect"
	"sort"

	"github.com/tryor/util/event"
)

type IElement interface {
	//event.IDispatcher
	IEventSupport
	ICursorSupport
	IMoveSupport

	SetId(id string)
	GetId() string
	SetParent(e IElement)
	GetParent() IElement
	GetSelf() IElement
	IsModified() bool
	SetModified(b bool)
	IsVisible() bool
	SetVisible(b bool)
	IsSelected() bool
	SetSelected(b bool)
	MouseIsHovering() bool
	GetMouseHoveringElement() IElement //如果有子元素，返回鼠标正悬停的子元素
	Destroy()

	IsEventEnabled() bool
	SetEventEnabled(e bool)
	Draw(ge IGraphicsEngine) //调用此方法绘制此元素，ge为绘图引擎
	TrackEvent(e event.IEvent) bool
	fireMouseEnterEvent(e IMouseEvent) bool //鼠标进入元素事件
	fireMouseLeaveEvent(e IMouseEvent) bool //鼠标离开元素事件
	fireFocusEvent(focus bool) bool
	firePaintEvent() bool

	Intersects(x, y int) bool //检查是否与点(x, y)相交
	IntersectsWith(rect *image.Rectangle) bool
	IntersectsWiths(rects ...*image.Rectangle) bool
	IntersectsElement(el IElement) bool
	IntersectsElements(els ...IElement) bool

	//IntersectsSnapshotWith(rect *image.Rectangle) bool

	SetAlignment(a Alignment) //设置水平和垂直对齐方式
	GetAlignment() Alignment
	SetAnchorPoint(x, y REAL)
	GetAnchorPoint() (x, y REAL)
	SetCoordinate(x, y int)
	GetCoordinate() (x, y int)      //返回在层或父元素中的相对坐标
	GetWorldCoordinate() (x, y int) //返回世界坐标，即在层中的坐标
	X() int
	Y() int
	Width() int
	Height() int
	SetWidth(w int)
	SetHeight(h int)

	SetRedraw(b bool)
	//SetChildrenRedraw(b bool)
	IsRedraw() bool
	RedrawIntersection(excludeds ...IElement)
	RedrawAll(excludeds ...IElement)
	SetLayer(layer ILayer)
	GetLayer() ILayer
	CreatePath()
	CreateBoundRect() *image.Rectangle
	GetBoundRect(adjustVal ...int) *image.Rectangle
	//MakeBoundRectSnapshot()
	//GetBoundRectSnapshot(adjustVal ...int) *image.Rectangle

	IsFocus() bool
	GetClipRegionAdjustValue() int

	IsActive() bool
	SetActive(b bool)

	SetObstacle(b int8)
	GetObstacle() int8
	IsObstacle() bool
	SetTag(tag string)
	GetTag() string
	//	SetDesc(desc string)
	//	GetDesc() string
	SetType(typ ElementType)
	GetType() ElementType

	HasChild(e IElement) bool
	GetChildIndex(e IElement) int
	DrawChildren(ge IGraphicsEngine)
	AddChild(e IElement, idx ...int) error
	AddChildren(els []IElement, pos ...int) error
	GetChildren() IElements
	GetChildrenCount() int
	GetChildrenFocusElement() IElement
	RemoveChild(e IElement) bool
	GetChild(id string) IElement
	ExistChild(id string) bool
	ClearChildren()
	ClearChildrenSortFlag()
	SetOrderZ(idx int)
	GetOrderZ() int

	IsDestroyed() bool
}

type Element struct {
	EventSupport
	//	VisibleSupport
	CursorSupport
	ModifiedSupport
	MoveSupport

	id       string
	parent   IElement
	children *Elements
	orderZ   int //在父节点或在层中的Z位置

	Self IElement //由于go不是全面支持OOP， 定义此属性用于在子对象中将自己传到父对象中使用

	anchorPoint PointF    //锚点，将通过此值计算x,y坐标
	alignment   Alignment //在父元素中的对齐方式,对齐方式的参考值为anchorPoint值，其实就是针对anchorPoint坐标点的对齐
	x, y        int       // x,y坐标，w,h高度宽度((x,y)coordinates, (w, h)height and width)
	w, h        int

	Layer ILayer //元素所在层

	boundRect         *image.Rectangle //边界矩形, 此区域不一定等于X, Y, W, H
	boundRectSnapshot *image.Rectangle //边界矩形快照，用于渲染时用,渲染时会创建边界矩形快照

	ClipRegionAdjustValue int //被剪裁区域的调整值
	//redrawElements        *Elements

	//mLeaveEvent   *MouseEvent    //定义鼠标移出事件
	//mEnterEvent   *MouseEvent    //定义鼠标移入事件
	//	selectEvent   *SelectEvent   //定义选择事件
	//	modifiedEvent *ModifiedEvent //定义改变被修改状态事件
	//	paintEvent *PaintEvent //定义绘制元素事件
	movedEvent *MovedEvent
	//	focusEvent *FocusEvent

	obstacle int8
	tag      string
	typ      ElementType

	autoDrawChildren bool //指示是否自动画子元素, 默认true
	destroyed        bool //指示是否已经被销毁
	visible          bool //是否可见，默认true
	redrawflag       bool //仅重画自己标记
	mouseHovering    bool //鼠标正在此元素上移动时，设置此标记
	selected         bool //选中标记
	active           bool //是否活动，当元素不在可视区域时，会自动设置为不活动
	eventEnabled     bool //是否能响应事件

}

func NewElement() *Element {
	e := newElement()
	e.children = NewElements()
	return e
}

func newElement() *Element {
	e := &Element{eventEnabled: true, ClipRegionAdjustValue: 1}
	e.Self = e
	e.id = fmt.Sprintf("%p", e)
	e.anchorPoint = PointF{0.5, 0.5}
	e.visible = true
	e.active = true
	e.autoDrawChildren = true
	e.SetModified(true)

	//	e.Dispatcher = event.NewDispatcher()
	//e.mLeaveEvent = NewMouseEvent(MOUSE_LEAVE_EVENT_TYPE, e, 0, 0, MButton_No, ModNone)
	//e.mEnterEvent = NewMouseEvent(MOUSE_ENTER_EVENT_TYPE, e, 0, 0, MButton_No, ModNone)
	//	e.selectEvent = NewSelectEvent(e, false)
	//	e.modifiedEvent = NewModifiedEvent(e, &e.ModifiedSupport)
	//	e.paintEvent = NewPaintEvent(e)
	//e.movedEvent = NewMovedEvent(e)
	//	e.focusEvent = NewFocusEvent(e, false)

	e.OnMoved(func(me *MovedEvent) {
		if e.children.Size() > 0 {
			for _, el := range e.children.CloneElements() {
				el.FireEvent(me)
			}
			//			e.children.ForEach(func(i int, el IElement) bool {
			//				el.FireEvent(me)
			//				return true
			//			})
		}

	})

	return e
}

func (this *Element) Destroy() {
	if this.destroyed {
		log.Println("Element " + this.id + " is already destroyed")
		return
	}

	this.destroyed = true
	//this.fireDestroyEvent(NewDestroyEvent(this.Self))
	this.Self.FireEvent(NewDestroyEvent(this.Self))

	//	for _, child := range this.children.GetElements() {
	//		child.Destroy()
	//	}
	this.children.ForEach(func(i int, el IElement) bool {
		el.Destroy()
		return true
	})

	this.children.Clear()

	//this.children.Destroy()
	//	if this.redrawElements != nil {
	//		this.redrawElements.Destroy()
	//	}

}

func (this *Element) IsAutoDrawChildren() bool {
	return this.autoDrawChildren
}

func (this *Element) SetAutoDrawChildren(b bool) {
	this.autoDrawChildren = b
}

func (this *Element) IsVisible() bool {
	return this.visible
}

func (this *Element) SetVisible(b bool) {
	this.visible = b
}

func (this *Element) IsDestroyed() bool {
	return this.destroyed
}

func (this *Element) SetId(id string) {
	this.id = id
}

func (this *Element) GetId() string {
	return this.id
}

func (this *Element) SetParent(e IElement) {
	this.parent = e
}

func (this *Element) GetParent() IElement {
	return this.parent
}

func (this *Element) SetType(typ ElementType) {
	this.typ = typ
}

func (this *Element) GetType() ElementType {
	return this.typ
}

//func (this *Element) SetDesc(desc string) {
//	//this.desc = desc
//}

//func (this *Element) GetDesc() string {
//	return ""
//}

func (this *Element) SetTag(tag string) {
	this.tag = tag
}
func (this *Element) GetTag() string {
	return this.tag
}

func (this *Element) SetObstacle(b int8) {
	this.obstacle = b
}

func (this *Element) GetObstacle() int8 {
	return this.obstacle
}

func (this *Element) IsObstacle() bool {
	return this.obstacle > ObstacleNo
}

func (this *Element) IsActive() bool {
	return this.active
}

func (this *Element) SetActive(b bool) {
	this.active = b
	//	for _, e := range this.children.GetElements() {
	//		e.SetActive(b)
	//	}
	this.children.ForEach(func(i int, el IElement) bool {
		el.SetActive(b)
		return true
	})
}

//func (this *Element) SetActive(b bool) {
//	this.active = b
//	//	if b {
//	//		atomic.StoreInt32(&this.active, 1)
//	//	} else {
//	//		atomic.StoreInt32(&this.active, 0)
//	//	}
//}

func (this *Element) GetSelf() IElement {
	return this.Self.(IElement)
}

//func (this *Element) ClearSelf() {
//	this.Self = nil
//}

func (this *Element) SetRedraw(b bool) {
	this.redrawflag = b
	p := this.parent
	if p != nil {
		p.SetRedraw(b)
	}
}

//func (this *Element) SetChildrenRedraw(b bool) {
//	this.redrawflag = b
//	for _, e := range this.children.GetElements() {
//		e.SetChildrenRedraw(b)
//	}
//}

func (this *Element) IsRedraw() bool {
	return this.redrawflag
}

//func (this *Element) setLayer(layer ILayer) {
//	this.Layer = layer
//	this.Self.(IElement).RedrawIntersection()
//}

func (this *Element) SetLayer(l ILayer) {
	this.Layer = l
	//this.Self.(IElement).RedrawIntersection()
	this.children.SetLayer(l)
}

func (this *Element) GetLayer() ILayer {
	return this.Layer
}

func (this *Element) MouseIsHovering() bool {
	return this.mouseHovering
}

func (this *Element) SetEventEnabled(e bool) {
	this.eventEnabled = e
}

func (this *Element) IsEventEnabled() bool {
	return this.eventEnabled
}

func (this *Element) GetMouseHoveringElement() IElement {
	return this.children.GetMouseHoveringElement()
}

func (this *Element) GetClipRegionAdjustValue() int {
	return this.ClipRegionAdjustValue
}

func (this *Element) Draw(ge IGraphicsEngine) {
	if !this.destroyed && this.autoDrawChildren && !this.children.Empty() {
		//this.DrawChildren(ge)
		this.children.Draw(ge)
	}
}

////设置水平对齐方式
//func (this *Element) SetHorizontalAlignment(at AlignmentType) {
//	this.alignment.Horizontal = at
//}

////设置垂直对齐方式
//func (this *Element) SetVerticalAlignment(at AlignmentType) {
//	this.alignment.Vertical = at
//}

//设置水平或垂直对齐方式
func (this *Element) SetAlignment(a Alignment) {
	this.alignment = a
}

func (this *Element) GetAlignment() Alignment {
	return this.alignment
}

func (this *Element) SetAnchorPoint(x, y REAL) {
	this.anchorPoint.X = x
	this.anchorPoint.Y = y
}

func (this *Element) GetAnchorPoint() (x, y REAL) {
	return this.anchorPoint.X, this.anchorPoint.Y
}

func (this *Element) SetWidth(w int) {
	//	w_ := int32(w)
	//	if this.w != w_ {
	//		//atomic.StoreInt32(&this.w, w_)
	//		this.Self.(IElement).SetModified(true)
	//	}
	if this.w != w {
		this.w = w
		this.Self.(IElement).SetModified(true)
	}
}

func (this *Element) SetHeight(h int) {
	//	h_ := int32(h)
	//	if this.h != h_ {
	//		//		this.h = h
	//		atomic.StoreInt32(&this.h, h_)
	//		//		this.Self.(IElement).ClearBoundRect()
	//		this.Self.(IElement).SetModified(true)
	//	}
	if this.h != h {
		this.h = h
		this.Self.(IElement).SetModified(true)
	}
}

func (this *Element) SetCoordinate(x, y int) {
	if this.x == x && this.y == y {
		return
	}
	this.x, this.y = x, y
	this.Self.(IElement).SetModified(true)

	//	atomic.StoreInt32(&this.x, int32(x))
	//	atomic.StoreInt32(&this.y, int32(y))
	//	this.Self.(IElement).SetModified(true)

}

func (this *Element) GetCoordinate() (x, y int) {
	//x, y = int(atomic.LoadInt32(&this.x)), int(atomic.LoadInt32(&this.y))
	x, y = this.x, this.y
	return
}

//func (this *Element) GetWorldCoordinate() (x, y int) {
//	x, y = int(atomic.LoadInt32(&this.x)), int(atomic.LoadInt32(&this.y))
//	p := this.parent
//	if p != nil {
//		x1, y1 := p.GetWorldCoordinate()
//		x += x1
//		y += y1
//	}
//	return
//}

//func (this *Element) getBoundRectMinPoint() (x, y int) {
//	x, y = this.Self.(IElement).GetCoordinate()
//	if this.anchorPointX != 0.0 {
//		x -= int(float32(this.Width()) * this.anchorPointX)
//	}
//	if this.anchorPointY != 0.0 {
//		y -= int(float32(this.Height()) * this.anchorPointY)
//	}
//	return
//}

func (this *Element) GetWorldCoordinate() (x, y int) {
	x, y = this.Self.(IElement).GetCoordinate()
	w, h := this.Width(), this.Height()
	pw, ph := 0, 0
	p := this.parent
	if p != nil {
		pw, ph = p.Width(), p.Height()
		x1, y1 := p.GetWorldCoordinate()
		x += x1
		y += y1
	} else {
		l := this.Layer
		if l != nil {
			pw, ph = l.Width(), l.Height()
		} else if layer, ok := this.Self.(ILayer); ok {
			page := layer.GetDrawPage()
			pw, ph = page.Width(), page.Height()
		}
	}

	hw := 0
	if this.anchorPoint.X != 0.0 {
		hw = int(REAL(w) * this.anchorPoint.X)
		x -= hw
	}
	hy := 0
	if this.anchorPoint.Y != 0.0 {
		hy = int(REAL(h) * this.anchorPoint.Y)
		y -= hy
	}

	switch this.alignment.Horizontal {
	case AlignmentNear:
	case AlignmentCenter:
		x += pw / 2
	case AlignmentFar:
		x += pw
	}

	switch this.alignment.Vertical {
	case AlignmentNear:
	case AlignmentCenter:
		y += ph / 2
	case AlignmentFar:
		y += ph
	}

	return
}

func (this *Element) X() int {
	//return int(atomic.LoadInt32(&this.x))
	return this.x
}
func (this *Element) Y() int {
	//return int(atomic.LoadInt32(&this.y))
	return this.y
}

func (this *Element) Width() int {
	//	return int(atomic.LoadInt32(&this.w))
	return this.w
}
func (this *Element) Height() int {
	//	return int(atomic.LoadInt32(&this.h))
	return this.h
}

/**
 * 增量移动到目标位置
 */
func (this *Element) MoveBy(dx, dy int, angle float32) {
	//	if this.X()+dx > 4000 || this.Y()+dy > 4000 {
	//	println("Element.MoveBy:", this, this.Self, this.X()+dx, this.Y()+dy, dx, dy)
	//	}

	this.Self.(IElement).SetCoordinate(this.X()+dx, this.Y()+dy)
	this.fireMovedEvent(dx, dy, angle)
}

/**
 * 移动到目标坐标
 */
func (this *Element) MoveTo(x, y int, angle float32) {
	this.Self.(IElement).MoveBy(x-int(this.ReferencePointX()), y-int(this.ReferencePointY()), angle)
	this.Self.(IElement).PrepareTransform(x, y)
}

/**
 * 构建路径
 */
func (this *Element) CreatePath() {
	if !this.destroyed {
		//this.Self.(IElement).CreateBoundRect()
		if !this.children.Empty() {
			this.children.CreatePath()
		}
	}
}

func (this *Element) SetModified(b bool) {
	p := this.parent
	if p != nil {
		p.SetModified(b)
	}

	if this.Layer != nil && b {
		this.Layer.SetModified(b)
	}

	if this.IsModified() == b {
		return
	}
	this.modified = b
	//this.ModifiedSupport.SetModified(b)
}

/**
 * 设置选中状态, 参数fireEvent指示是否产生SelectEvent事件
 */
func (this *Element) SetSelected(b bool) {
	if this.selected != b {
		this.selected = b
		this.Self.(IElement).SetModified(true)
	}
}

func (this *Element) IsSelected() bool {
	return this.selected
}

/**
 * 跟踪响应事件
 */
func (this *Element) TrackEvent(ie event.IEvent) bool {
	this.children.TrackEvent(ie)

	if !this.eventEnabled {
		return false
	}
	if IsMouseEvent(ie.GetType()) {
		//如果是鼠标事件
		return this.handleMouseEvent(ie.(IMouseEvent))
	} else {
		return this.Self.(IElement).FireEvent(ie)
	}

	return true
}

//func (this *Widget) TrackEvent(e event.IEvent) bool {
//	this.children.TrackEvent(e)
//	return this.trackEvent(e)
//}

/**
 * 重画与此元素相交的其它元素, 排除<code>excludeds</code>中的元素
 */
func (this *Element) RedrawIntersection(excludeds ...IElement) {
	if this.Layer != nil && this.Layer.GetDrawMode() == DrawMode_Region {
		redrawElements := this.prepareRedrawElements()
		//println("Element.RedrawIntersection.redrawElements.Size()", this.redrawElements.Size())
		redrawElements.RedrawIntersection(excludeds...)
		//redrawElements.Destroy()
	}
}

/**
 * 重画与此元素相交的元素, 包括相交元素和自己, 排除<code>excludeds</code>中的元素
 */
func (this *Element) RedrawAll(excludeds ...IElement) {
	if this.Layer != nil && this.Layer.GetDrawMode() == DrawMode_Region {
		redrawElements := this.prepareRedrawElements()
		redrawElements.RedrawAll(excludeds...)
		//redrawElements.Destroy()
	}
}

/**
 * 检查是否与点(x, y)相交
 */
func (this *Element) Intersects(x, y int) bool {

	if this.IsDestroyed() {
		//println("Element.Intersects.IsDestroyed", this.IsDestroyed())
		return false
	}
	//	if this.Self == nil {
	//		println("Element.Intersects.this.Self", this.Self)
	//	}
	r := this.Self.(IElement).GetBoundRect()
	if r.Dx() < 0 || r.Dy() < 0 {
		r = FormatRect(r)
	}
	return r.Min.X <= x && x < r.Max.X &&
		r.Min.Y <= y && y < r.Max.Y
}

/**
 * 检查是否与参数矩形区相交
 */
func (this *Element) IntersectsWith(rect *image.Rectangle) bool {
	br := this.Self.(IElement).GetBoundRect()
	if br != nil && rect != nil {
		//		if Intersect(br, rect) != (br.Intersect(*rect) != image.ZR) {
		//			println("Element.IntersectsWith", br.String(), rect.String(), Intersect(br, rect), br.Intersect(*rect).String())
		//		}

		//		r := this.BoundRect.Intersect(*rect)
		//		return !(r.Min.X == 0 && r.Min.Y == 0 && r.Max.X == 0 && r.Max.Y == 0)
		//		return (*br).Intersect(*rect) != image.ZR
		return IsIntersect(br, rect)
	}
	return false
}

/**
 * 检查是否与参数矩形区相交, 与rects中任何一个矩形相交就满足条件
 */
func (this *Element) IntersectsWiths(rects ...*image.Rectangle) bool {
	for _, r := range rects {
		if this.Self.(IElement).IntersectsWith(r) {
			return true
		}
	}
	return false
}

func (this *Element) IntersectsElement(el IElement) bool {
	return this.Self.(IElement).IntersectsWith(el.GetBoundRect())
}

func (this *Element) IntersectsElements(els ...IElement) bool {
	for _, el := range els {
		if this.Self.(IElement).IntersectsElement(el) {
			return true
		}
	}
	return false
}

///**
// * 检查边界快照是否与参数矩形区相交
// */
//func (this *Element) IntersectsSnapshotWith(rect *image.Rectangle) bool {
//	br := this.Self.(IElement).GetBoundRectSnapshot()
//	if br != nil && rect != nil {
//		return IsIntersect(br, rect)
//	}
//	return false
//}

/**
 * 创建包含元素所有可视内容的最大矩形区域
 */
func (this *Element) CreateBoundRect() *image.Rectangle {
	x, y := this.Self.(IElement).GetWorldCoordinate() //this.Self.(IElement).GetWorldCoordinate() //this.GetWorldCoordinate()
	w, h := this.Self.(IElement).Width(), this.Self.(IElement).Height()

	boundRect := this.boundRect
	if boundRect == nil ||
		!(boundRect.Min.X == x && boundRect.Min.Y == y && boundRect.Dx() == w && boundRect.Dy() == h) {
		boundRect =
			&image.Rectangle{Min: image.Point{x, y},
				Max: image.Point{x + w, y + h}}
		this.boundRect = boundRect
	}

	//	if !this.children.Empty() {
	//		childrenBR := this.children.CreateBoundRect()
	//		boundRect = Union(childrenBR, boundRect)
	//	}

	return boundRect
}

//func (this *Element) ClearBoundRect() {
//	//	this.boundRectLocker.Lock()
//	//	defer this.boundRectLocker.Unlock()
//	this.boundRect = nil

//	//	for _, e := range this.children.GetElements() {
//	//		e.ClearBoundRect()
//	//	}
//}

//func (this *Element) SetBoundRect(boundRect *image.Rectangle) {
//	this.boundRect = boundRect //&image.Rectangle{Min: boundRect.Min, Max: boundRect.Max}
//}

//func (this *Element) getBoundRect() *image.Rectangle {
//	//	this.boundRectLocker.RLock()
//	//	defer this.boundRectLocker.RUnlock()
//	return this.boundRect
//}

/**
 * 在outBoundRect中返回包含元素所有可视内容的最大矩形, adjustVal为调整值，如果为正值，将放大矩形，负值缩小矩形
 */
func (this *Element) GetBoundRect(adjustVal ...int) *image.Rectangle {

	//	br := this.getBoundRect()
	//	if br == nil {
	//		br = this.Self.(IElement).CreateBoundRect()
	//	}
	//	br := this.getBoundRect()
	br := this.boundRect
	if br == nil {
		br = this.Self.(IElement).CreateBoundRect()
	}
	x, y, x1, y1 := br.Min.X, br.Min.Y, br.Max.X, br.Max.Y
	if len(adjustVal) > 0 && adjustVal[0] != 0 {

		v := adjustVal[0]
		if x < x1 {
			x -= v
			x1 += v
		} else {
			x += v
			x1 -= v
		}
		if y < y1 {
			y -= v
			y1 += v
		} else {
			y += v
			y1 -= v
		}

		return &image.Rectangle{Min: image.Point{x, y},
			Max: image.Point{x1, y1}}
	} else {
		return br
	}

	//	r := image.Rect(x, y, x1, y1)
	//	return &r
}

//func (this *Element) MakeBoundRectSnapshot() {
//	this.boundRectSnapshot = this.GetBoundRect()
//	if this.children.Size() > 0 {
//		this.children.MakeBoundRectSnapshot()
//	}
//}

//func (this *Element) GetBoundRectSnapshot(adjustVal ...int) *image.Rectangle {
//	br := this.boundRectSnapshot
//	if br == nil {
//		br = this.Self.(IElement).CreateBoundRect()
//		//panic("boundrect snapshot not exist!")
//	}
//	x, y, x1, y1 := br.Min.X, br.Min.Y, br.Max.X, br.Max.Y
//	if len(adjustVal) > 0 && adjustVal[0] != 0 {

//		v := adjustVal[0]
//		if x < x1 {
//			x -= v
//			x1 += v
//		} else {
//			x += v
//			x1 -= v
//		}
//		if y < y1 {
//			y -= v
//			y1 += v
//		} else {
//			y += v
//			y1 -= v
//		}

//		return &image.Rectangle{Min: image.Point{x, y},
//			Max: image.Point{x1, y1}}
//	} else {
//		return br
//	}
//}

//func (this *Element) isFocus() bool {
//	l := this.Layer
//	return l != nil && l.GetFocusElement() == this.Self
//}

func (this *Element) IsFocus() bool {
	if l := this.Layer; l != nil && l.GetFocusElement() == this.Self {
		return true
	}
	p := this.parent
	return p != nil && p.GetChildrenFocusElement() == this.Self
}

/**
 * 发送鼠标离开事件
 */
func (this *Element) fireMouseLeaveEvent(e IMouseEvent) bool {
	//println("Element.fireMouseLeaveEvent ", this.Self, this.mouseHovering)

	this.children.fireMouseLeaveEvent(e)

	if !this.mouseHovering {
		return true
	}

	if l := this.Layer; l != nil && l.GetMouseHoveringElement() == this.Self {
		this.mouseHovering = false
		l.ClearMouseHoveringElement()
	}

	mLeaveEvent := NewMouseEvent(MOUSE_LEAVE_EVENT_TYPE, e.GetSource(), e.X(), e.Y(), e.GetButtons(), e.GetModifier())
	//	this.mLeaveEvent.x = e.X()
	//	this.mLeaveEvent.y = e.Y()
	//	this.mLeaveEvent.Buttons = e.GetButtons()
	//	this.mLeaveEvent.Modifier = e.GetModifier()
	mLeaveEvent.KeySequence = e.GetKeySequence()

	//向在此层上监听的对象发送鼠标离开事件
	return this.Self.(IElement).FireEvent(mLeaveEvent)
}

/**
 * 发送鼠标进入元素事件
 */
func (this *Element) fireMouseEnterEvent(e IMouseEvent) bool {
	//	println("Element.fireMouseEnterEvent ", this.Self, this.MouseHovering)
	this.mouseHovering = true
	mEnterEvent := NewMouseEvent(MOUSE_ENTER_EVENT_TYPE, e.GetSource(), e.X(), e.Y(), e.GetButtons(), e.GetModifier())
	//	this.mEnterEvent.x = e.X()
	//	this.mEnterEvent.y = e.Y()
	//	this.mEnterEvent.Buttons = e.GetButtons()
	//	this.mEnterEvent.Modifier = e.GetModifier()
	mEnterEvent.KeySequence = e.GetKeySequence()
	return this.Self.(IElement).FireEvent(mEnterEvent)
}

func (this *Element) fireFocusEvent(focus bool) bool {
	this.children.fireFocusEvent(focus)
	//focusEvent = NewFocusEvent(this.Self, focus)
	//this.focusEvent.Focus = focus
	//this.focusEvent.Source = this.Self
	return this.Self.(IElement).FireEvent(NewFocusEvent(this.Self, focus))
}

/**
 * 处理鼠标事件
 */
func (this *Element) handleMouseEvent(me IMouseEvent) bool {
	//println("Element.handleMouseEvent")
	//设置鼠标悬停状态
	this.mouseHovering = true
	//在本元素中分发鼠标事件
	return this.Self.(IElement).FireEvent(me)
}

/**
 * 抛出绘制元素事件
 */
func (this *Element) firePaintEvent() bool {
	if this.onPaintEvent != nil {
		return this.Self.(IElement).FireEvent(NewPaintEvent(this.Self))
	}
	return false
}

func (this *Element) fireMovedEvent(dx, dy int, angle float32) bool {
	if this.movedEvent == nil {
		this.movedEvent = NewMovedEvent(this.Self)
	}
	this.movedEvent.Dx = dx
	this.movedEvent.Dy = dy
	this.movedEvent.Angle = angle
	return this.Self.(IElement).FireEvent(this.movedEvent)
}

func (this *Element) prepareRedrawElements() *Elements {
	//	if this.redrawElements == nil {
	//		this.redrawElements = NewElements(this.Self.(IElement))
	//		//this.redrawElements.Add(this.Self.(IElement))
	//		this.redrawElements.SetLayer(this.Layer)
	//	}
	redrawElements := NewElements(this.Self.(IElement))
	redrawElements.SetLayer(this.Layer)
	return redrawElements
}

func (this *Element) HasChild(e IElement) bool {
	return this.children.Contains(e)
}

func (this *Element) GetChildIndex(e IElement) int {
	return this.children.GetIndex(e)
}

func (this *Element) DrawChildren(ge IGraphicsEngine) {
	if !this.children.Empty() {
		this.children.Draw(ge)
	}
}

func (this *Element) AddChild(e IElement, idx ...int) error {
	//	err := this.children.Add(e, idx...)
	//	if err == nil {
	//		e.SetParent(this.Self.(IElement))
	//		if this.Layer != nil {
	//			e.SetLayer(this.Layer)
	//		}
	//		return nil
	//	}
	//	return err
	err := this.addChild(e, idx...)
	if err == nil {
		e.SetParent(this.Self.(IElement))
		return nil
	}
	return err
}

func (this *Element) addChild(e IElement, idx ...int) error {
	err := this.children.Add(e, idx...)
	if err == nil {
		//e.SetParent(this.Self.(IElement))
		if this.Layer != nil {
			e.SetLayer(this.Layer)
		}
		return nil
	}
	return err
}

//func (this *Element) AddChildren(es ...IElement) error {
func (this *Element) AddChildren(es []IElement, pos ...int) error {
	for _, e := range es {
		//		e.SetParent(this.Self.(IElement))
		err := this.Self.(IElement).AddChild(e, pos...)
		if err != nil {
			return err
		}
	}
	return nil
	//	err := this.children.Adds(es...)
	//	if err == nil {
	//		if this.Layer != nil {
	//			this.children.SetLayer(this.Layer)
	//		}
	//		return nil
	//	}
	//	return err
}

func (this *Element) GetChildren() IElements {
	return this.children
}

func (this *Element) GetChildrenCount() int {
	return this.children.Size()
}

func (this *Element) GetChildrenFocusElement() IElement {
	return this.children.GetFocusElement()
}

func (this *Element) RemoveChild(e IElement) bool {
	e.RedrawIntersection()
	return this.children.Remove(e)
}

func (this *Element) GetChild(id string) IElement {
	return this.children.GetById(id)
}

func (this *Element) ExistChild(id string) bool {
	return this.children.Exist(id)
}

func (this *Element) ClearChildren() {
	this.children.RedrawIntersection()
	this.children.Clear()
	//	this.Self.(IElement).SetModified(true)
}

func (this *Element) ClearChildrenSortFlag() {
	this.children.ClearSortFlag()
}

func (this *Element) SetOrderZ(z int) {
	if this.orderZ == z {
		return
	}
	this.orderZ = z
	if this.parent != nil {
		this.parent.ClearChildrenSortFlag()
	}

}
func (this *Element) GetOrderZ() int {
	return this.orderZ
}

func (this *Element) SetMoving(b bool) {
	if !this.children.Empty() {
		this.children.SetMoving(b)
	}
	this.MoveSupport.SetMoving(b)
}

func (this *Element) PrepareTransform(x, y int) {
	if !this.children.Empty() {
		this.children.PrepareTransform(x, y)
	}
	this.MoveSupport.PrepareTransform(x, y)
}

func (this *Element) GetGraphicsEngine() IGraphicsEngine {
	if layer := this.GetLayer(); layer != nil {
		return layer.GetGraphicsEngine()
	}
	return nil
}

type ElementSortor struct {
	elements []IElement
	comparer func(elements []IElement, i, j int) bool
}

func NewElementSortor(elements []IElement, comparer func(elements []IElement, i, j int) bool) *ElementSortor {
	if comparer == nil {
		comparer = func(elements []IElement, i, j int) bool {
			return reflect.ValueOf(elements[i]).Pointer() < reflect.ValueOf(elements[j]).Pointer()
		}
	}
	return &ElementSortor{elements: elements, comparer: comparer}
}

func (this *ElementSortor) SetElements(elements []IElement) {
	this.elements = elements
}

func (this *ElementSortor) GetElements() []IElement {
	return this.elements
}

func (this *ElementSortor) Len() int {
	return len(this.elements)
}

func (this *ElementSortor) Swap(i, j int) {
	this.elements[i], this.elements[j] = this.elements[j], this.elements[i]
}

func (this *ElementSortor) Less(i, j int) bool {
	return this.comparer(this.elements, i, j)
}

func (this *ElementSortor) Sort() {
	sort.Sort(this)
}

func (this *ElementSortor) Search(e IElement) int {
	return searchElement(this.elements, e)
}

func (this *ElementSortor) Contains(e IElement) bool {
	idx := searchElement(this.elements, e)
	return idx >= 0 && idx < len(this.elements)
}

func searchElement(es []IElement, e IElement) int {
	return sort.Search(len(es), func(i int) bool { return es[i] == e })
}
