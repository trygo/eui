package eui

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/tryor/commons/event"
	//	"log"
	"sort"

	. "github.com/tryor/winapi"
	//	"sync/atomic"
)

type ILayer interface {
	IElement

	SetLayerType(typ LayerType)
	GetLayerType() LayerType

	Init()
	SetFocusElement(e IElement)
	GetFocusElement() IElement

	ClearMouseHoveringElement()
	GetMouseEventMode() MouseEventMode
	GetManualModeElement() IElement
	GetDrawPage() IDrawPage
	SetDrawPage(page IDrawPage)
	GetSelection() ISelection
	SetDrawMode(dm DrawMode)
	GetDrawMode() DrawMode
	UnionDrawRegion(rect *image.Rectangle)
	ClipDrawRegion(ge IGraphicsEngine)
	ResetClipedDrawRegion(ge IGraphicsEngine)

	GetIntersections(intersects func(e IElement) bool, excludeds *ElementSortor, appendeds *ElementSortor) []IElement
	GetIntersectionsByElement(el IElement, excludeds *ElementSortor, appendeds *ElementSortor) []IElement
	GetIntersectionsByElements(els []IElement, excludeds *ElementSortor, appendeds *ElementSortor) []IElement
	GetIntersectionsByRect(rect *image.Rectangle, excludeds *ElementSortor, appendeds *ElementSortor) []IElement
	GetIntersectionsByRects(rects []*image.Rectangle, excludeds *ElementSortor, appendeds *ElementSortor) []IElement
	GetIntersectsWithPointInVisibleRegion(x, y int, excludeds ...IElement) []IElement //返回在可视区域与指定点相交的元素
	GetIntersectsWithRectInVisibleRegion(rect *image.Rectangle, excludeds ...IElement) []IElement
	IsRedrawAll() bool
	SetRedrawAll()

	GetBackground() color.Color
	SetBackground(c color.Color)

	//	SortElements(ordereds []IElement) []IElement

	Adjust(e IElement, indexBy int, target IElement)
	SwitchPosition(e IElement) //交换位置，根据元素低部坐标交换位置

	GetDrawElementsCount() int
	GetDrawSpendTime() int64
	Find(x, y int, excluded IElement) IElement

	SetGraphicsEngine(ge IGraphicsEngine)
	GetGraphicsEngine() IGraphicsEngine
	//	IsLayerGraphicsEngine() bool //IsLayerGraphicsEngine or IsPageGraphicsEngine

	GetVisibleRegionCoordinate() (x, y int)
	GetVisibleRegion(clearflag ...bool) (*image.Rectangle, *image.Rectangle)
	GetVisibleRegionF(clearflag ...bool) RectF
	SetVisibleRegion(vr image.Rectangle)
	AdjustVisibleRegion(dx, dy REAL)
	IsVisibleRegionModified() bool //返回true, 如果视口区域被调整
	Scrollable() bool              //同IsVisibleRegionModified()
	GetLastVisibleRegion() RectF

	//	MakeVisibleRegionSnapshot()
	//	GetVisibleRegionSnapshot() (vr *image.Rectangle, lvr *image.Rectangle, vrModified bool)
	//	GetVisibleRegionSnapshotCoordinate() (x, y int)

	SetLastVisibleRegionElements(els []IElement)
	GetLastVisibleRegionElements() []IElement

	SetScrollrateX(v float32)
	GetScrollrateX() float32
	SetScrollrateY(v float32)
	GetScrollrateY() float32

	//IsRendering() bool
}

type Layer struct {
	*Element

	//	id        string
	layerType LayerType

	focusElement         IElement //当前焦点元素
	mouseHoveringElement IElement //鼠标悬停元素

	MouseEventMode    MouseEventMode //鼠标事件响应模式
	ManualModeElement IElement       //MEventMode_Manual模式下指定元素

	Page          IDrawPage  //绘画页
	Selection     ISelection //选择器
	DrawMode      DrawMode   //元素绘制模式
	scrollMode    ScrollMode //层卷动格式
	RedrawAllFlag bool       //重画所有元素标记, 一旦检查了此标记后，此标记将自动复位

	Background color.Color //背景颜色

	//rendering   int32
	drawedCount int   //最近一次绘制元素的数量
	spendTime   int64 //最近一次绘制层用时，单位：纳秒

	clipedRegion IRegion
	clipedRects  []*image.Rectangle //重画区域

	GraphicsEngine      IGraphicsEngine
	layerGraphicsEngine bool

	visibleRegionModified bool
	visibleRegion         RectF //可视区域
	lastVisibleRegion     RectF //上一次可视区域
	//	visibleRegionSnapshot     *image.Rectangle //可视区域快照,用于在渲染时锁定可高区域坐标位置
	//	lastVisibleRegionSnapshot *image.Rectangle //可视区域快照,用于在渲染时锁定可高区域坐标位置
	//	vrModifiedSnapshot        bool

	lastVisibleRegionElements []IElement //最近一次在可视区域的元素

	scrollrateX float32
	scrollrateY float32

	sortor *ElementSortor

	mergeRegions []*Rect //层合并时，可合并区域
	tmpx, tmpy   int
}

func NewLayer(page IDrawPage, selct ISelection) *Layer {
	l := &Layer{Element: NewElement(), Page: page, Selection: selct, MouseEventMode: MEventMode_Hovering,
		DrawMode: DrawMode_All, scrollMode: ScrollMode_Auto, Background: nil}
	l.Self = l
	l.anchorPoint = PointF{0.0, 0.0} //X = 0.0
	l.clipedRects = make([]*image.Rectangle, 0)
	if l.Selection == nil {
		l.Selection = NewSelection()
		l.Selection.SetLayer(l)
	}
	l.lastVisibleRegionElements = make([]IElement, 0)

	l.sortor = NewElementSortor(nil, func(elements []IElement, i, j int) bool {
		return elements[i].GetOrderZ() < elements[j].GetOrderZ()
	})

	l.mergeRegions = make([]*Rect, 0)

	return l
}

func (this *Layer) Destroy() {
	//if this.Selection != nil {
	//	this.Selection.Destroy()
	//}
	this.Element.Destroy()
}

func (this *Layer) GetMergeRegions() []*Rect {
	return this.mergeRegions
}

func (this *Layer) SetLayerType(typ LayerType) {
	this.layerType = typ
}

func (this *Layer) GetLayerType() LayerType {
	return this.layerType
}

func (this *Layer) SetScrollrateX(v float32) {
	if v > 1 {
		this.scrollrateX = 1
	} else if v < 0 {
		this.scrollrateX = 0
	} else {
		this.scrollrateX = v
	}
}

func (this *Layer) GetScrollrateX() float32 {
	return this.scrollrateX
}

func (this *Layer) SetScrollrateY(v float32) {
	//	if v > 1 {
	//		this.scrollrateY = 1
	//	} else if v < 0 {
	//		this.scrollrateY = 0
	//	} else {
	//		this.scrollrateY = v
	//	}
	this.scrollrateY = v
}

func (this *Layer) GetScrollrateY() float32 {
	return this.scrollrateY
}

func (this *Layer) SetLastVisibleRegionElements(els []IElement) {
	//	this.lastVisibleRegionElementsLocker.Lock()
	this.lastVisibleRegionElements = els
	//	this.lastVisibleRegionElementsLocker.Unlock()
}

func (this *Layer) GetLastVisibleRegionElements() []IElement {
	//	this.lastVisibleRegionElementsLocker.RLock()
	//	defer this.lastVisibleRegionElementsLocker.RUnlock()
	return this.lastVisibleRegionElements
}

/**
 * 检查是否与点(x, y)相交
 */
//func (this *Layer) Intersects(x, y int) bool {
//	thisX, thisY := this.GetCoordinate()
//	return thisX <= x && x < thisX+this.Width() &&
//		thisY <= y && y < thisY+this.Height()
//}

func (this *Layer) Intersects(x, y int) bool {
	r := this.Self.(ILayer).GetVisibleRegionF()
	//thisX, thisY := this.GetCoordinate()
	thisX, thisY := this.GetWorldCoordinate()
	return thisX <= x && x < thisX+int(r.W) &&
		thisY <= y && y < thisY+int(r.H)
}

//func (this *Layer) MakeVisibleRegionSnapshot() {
//	this.vrModifiedSnapshot = this.visibleRegionModified
//	this.visibleRegionSnapshot, this.lastVisibleRegionSnapshot = this.GetVisibleRegion(true)
//	//统计可视区域元素，@TODO 还未实现
//	for _, el := range this.GetChildren().CloneElements() {
//		el.MakeBoundRectSnapshot()
//	}
//	//this.ClearVisibleRegionModified()
//}

//func (this *Layer) GetVisibleRegionSnapshot() (vr *image.Rectangle, lvr *image.Rectangle, vrModified bool) {
//	return this.visibleRegionSnapshot, this.lastVisibleRegionSnapshot, this.vrModifiedSnapshot
//}

//func (this *Layer) GetVisibleRegionSnapshotCoordinate() (x, y int) {
//	return this.visibleRegionSnapshot.Min.X, this.visibleRegionSnapshot.Min.Y
//}

func (this *Layer) GetVisibleRegionCoordinate() (x, y int) {
	//	this.visibleRegionLocker.RLock()
	//	defer this.visibleRegionLocker.RUnlock()
	x, y = int(this.visibleRegion.X+0.5), int(this.visibleRegion.Y+0.5)
	//	if this.GetId() == "earth" {
	//		ax, _ := this.GetDrawPage().GetFocusLayer().GetVisibleRegionCoordinate()
	//		if ax != x {
	//			fmt.Println("Layer.GetVisibleRegionCoordinate:", this.GetId(), x, ax)
	//		}
	//	}
	return
	//return this.tmpx, this.tmpy
}

func (this *Layer) GetVisibleRegionF(clearflag ...bool) RectF {
	if this.visibleRegion == ZRF { //image.ZR {
		r := this.Self.GetBoundRect()
		this.visibleRegion = RectF{W: REAL(r.Dx()), H: REAL(r.Dy())}
	}
	if len(clearflag) > 0 && clearflag[0] {
		this.visibleRegionModified = false
	}
	return this.visibleRegion
}

func (this *Layer) GetVisibleRegion_old(clearflag ...bool) (vr *image.Rectangle, lvr *image.Rectangle) {
	if this.visibleRegion == ZRF {
		r := this.Self.GetBoundRect()
		this.visibleRegion = RectF{W: REAL(r.Dx()), H: REAL(r.Dy())}
		this.lastVisibleRegion = this.visibleRegion
	}
	vr = &image.Rectangle{Min: image.Point{int(this.visibleRegion.X), int(this.visibleRegion.Y)}, Max: image.Point{int(this.visibleRegion.X + this.visibleRegion.W), int(this.visibleRegion.Y + this.visibleRegion.H)}}                          //this.visibleRegion
	lvr = &image.Rectangle{Min: image.Point{int(this.lastVisibleRegion.X), int(this.lastVisibleRegion.Y)}, Max: image.Point{int(this.lastVisibleRegion.X + this.lastVisibleRegion.W), int(this.lastVisibleRegion.Y + this.lastVisibleRegion.H)}} //this.lastVisibleRegion

	if len(clearflag) > 0 && clearflag[0] {
		this.visibleRegionModified = false
	}

	return

}

func (this *Layer) GetVisibleRegion(clearflag ...bool) (vr *image.Rectangle, lvr *image.Rectangle) {
	//	this.visibleRegionLocker.Lock()
	//	defer this.visibleRegionLocker.Unlock()

	if this.visibleRegion == ZRF { //image.ZR {
		r := this.Self.GetBoundRect()
		this.visibleRegion = RectF{W: REAL(r.Dx()), H: REAL(r.Dy())}
		this.lastVisibleRegion = this.visibleRegion
	}

	vrX := int(this.visibleRegion.X + 0.5)
	vrY := int(this.visibleRegion.Y + 0.5)
	vrMX := int(this.visibleRegion.X + this.visibleRegion.W + 0.5)
	vrMY := int(this.visibleRegion.Y + this.visibleRegion.H + 0.5)

	lvrX := int(this.lastVisibleRegion.X + 0.5)
	lvrY := int(this.lastVisibleRegion.Y + 0.5)
	lvrMX := int(this.lastVisibleRegion.X + this.lastVisibleRegion.W + 0.5)
	lvrMY := int(this.lastVisibleRegion.Y + this.lastVisibleRegion.H + 0.5)

	vr = &image.Rectangle{Min: image.Point{vrX, vrY}, Max: image.Point{vrMX, vrMY}}
	lvr = &image.Rectangle{Min: image.Point{lvrX, lvrY}, Max: image.Point{lvrMX, lvrMY}}

	//vr = &image.Rectangle{Min: image.Point{int(this.visibleRegion.X), int(this.visibleRegion.Y)}, Max: image.Point{int(this.visibleRegion.X + this.visibleRegion.W), int(this.visibleRegion.Y + this.visibleRegion.H)}}                          //this.visibleRegion
	//lvr = &image.Rectangle{Min: image.Point{int(this.lastVisibleRegion.X), int(this.lastVisibleRegion.Y)}, Max: image.Point{int(this.lastVisibleRegion.X + this.lastVisibleRegion.W), int(this.lastVisibleRegion.Y + this.lastVisibleRegion.H)}} //this.lastVisibleRegion

	if len(clearflag) > 0 && clearflag[0] {
		this.tmpx, this.tmpy = vrX, vrY
		this.visibleRegionModified = false
	}

	return
}

//func (this *Layer) CreateBoundRect() *image.Rectangle {
//	x, y := this.Self.(IElement).GetWorldCoordinate()
//	w, h := this.Self.(IElement).Width(), this.Self.(IElement).Height()

//	boundRect := this.boundRect
//	if boundRect == nil ||
//		!(boundRect.Min.X == x && boundRect.Min.Y == y && boundRect.Dx() == w && boundRect.Dy() == h) {
//		boundRect =
//			&image.Rectangle{Min: image.Point{x, y},
//				Max: image.Point{x + w, y + h}}
//		this.boundRect = boundRect
//	}

//	return boundRect
//}

func (this *Layer) AdjustVisibleRegion(dx, dy REAL) {
	//	this.visibleRegionLocker.Lock()
	//	defer this.visibleRegionLocker.Unlock()
	if dx == 0.0 && dy == 0.0 {
		return
	}

	if !this.visibleRegionModified {
		this.lastVisibleRegion = this.visibleRegion
	}

	vr := this.visibleRegion
	x1, x2 := vr.X+dx, vr.X+vr.W+dx
	if x1 < 0 {
		x1 = 0
		x2 = vr.W
	} else {
		layerWidth := REAL(this.Width())
		if x2 > layerWidth {
			x1 = layerWidth - vr.W
			x2 = layerWidth
		}
	}

	y1, y2 := vr.Y+dy, vr.Y+vr.H+dy
	if y1 < 0 {
		y1 = 0
		y2 = vr.H
	} else {
		layerHeight := REAL(this.Height())
		if y2 > layerHeight {
			y1 = layerHeight - vr.H
			y2 = layerHeight
		}
	}

	this.visibleRegion.X = x1
	this.visibleRegion.W = x2 - x1
	this.visibleRegion.Y = y1
	this.visibleRegion.H = y2 - y1

	if !this.visibleRegionModified && this.visibleRegion != this.lastVisibleRegion {
		this.visibleRegionModified = true
		//this.Self.(ILayer).SetModified(true)
	}

	this.GetDrawPage().SetModified(true)
}

func (this *Layer) GetLastVisibleRegion() RectF {
	return this.lastVisibleRegion
}

//func (this *Layer) AdjustVisibleRegion(dx, dy float32) {
//	//	this.visibleRegionLocker.Lock()
//	//	defer this.visibleRegionLocker.Unlock()

//	this.Self.(ILayer).SetModified(true)

//	x1, x2 := this.visibleRegion.Min.X+dx, this.visibleRegion.Max.X+dx
//	if x1 < 0 {
//		x1 = 0
//		x2 = this.visibleRegion.Dx()
//	} else {
//		layerWidth := this.Width()
//		if x2 > layerWidth {
//			x1 = layerWidth - this.visibleRegion.Dx()
//			x2 = layerWidth
//		}
//	}
//	this.visibleRegion.Min.X = x1
//	this.visibleRegion.Max.X = x2

//	y1, y2 := this.visibleRegion.Min.Y+dy, this.visibleRegion.Max.Y+dy
//	if y1 < 0 {
//		y1 = 0
//		y2 = this.visibleRegion.Dy()
//	} else {
//		layerHeight := this.Height()
//		if y2 > layerHeight {
//			y1 = layerHeight - this.visibleRegion.Dy()
//			y2 = layerHeight
//		}
//	}
//	this.visibleRegion.Min.Y = y1
//	this.visibleRegion.Max.Y = y2

//	this.visibleRegionModified = true
//}

func (this *Layer) IsVisibleRegionModified() bool {
	//	this.visibleRegionLocker.RLock()
	//	defer this.visibleRegionLocker.RUnlock()
	return this.visibleRegionModified
}

func (this *Layer) Scrollable() bool {
	return this.visibleRegionModified
}

func (this *Layer) ClearVisibleRegionModified() {
	this.visibleRegionModified = false
}

func (this *Layer) SetVisibleRegion(vr image.Rectangle) {
	//	this.visibleRegionLocker.Lock()
	//	defer this.visibleRegionLocker.Unlock()
	this.visibleRegion = *NewRectF(REAL(vr.Min.X), REAL(vr.Min.Y), REAL(vr.Dx()), REAL(vr.Dy()))
	//fmt.Println("Layer.SetVisibleRegion", this.GetId(), this.GetTag())
	//this.visibleRegionModified = true
}

//返回在可视区域与指定区域相交的元素
func (this *Layer) GetIntersectsWithRectInVisibleRegion(rect *image.Rectangle, excludeds ...IElement) []IElement {
	lastVisibleEls := this.GetLastVisibleRegionElements()
	intersectElements := make([]IElement, 0)

	excludedsLen := len(excludeds)
	var excluded IElement
	var sortor *ElementSortor
	if excludedsLen > 1 {
		sortor := NewElementSortor(excludeds, nil)
		sortor.Sort()
	} else if excludedsLen == 1 {
		excluded = excludeds[0]
	}

	for _, el := range lastVisibleEls {

		if excluded != nil {
			if excluded == el {
				continue
			}
		} else if sortor != nil {
			idx := sort.Search(excludedsLen, func(i int) bool {
				return excludeds[i] == el
			})
			if idx > -1 {
				continue
			}
		}
		if el.IntersectsWith(rect) {
			intersectElements = append(intersectElements, el)
		}
	}
	return intersectElements
}

//返回在可视区域与指定点相交的元素
func (this *Layer) GetIntersectsWithPointInVisibleRegion(x, y int, excludeds ...IElement) []IElement {
	lastVisibleEls := this.GetLastVisibleRegionElements()
	intersectElements := make([]IElement, 0)

	excludedsLen := len(excludeds)
	var excluded IElement
	var sortor *ElementSortor
	if excludedsLen > 1 {
		sortor := NewElementSortor(excludeds, nil)
		sortor.Sort()
	} else if excludedsLen == 1 {
		excluded = excludeds[0]
	}

	for _, el := range lastVisibleEls {

		if excluded != nil {
			if excluded == el {
				continue
			}
		} else if sortor != nil {
			idx := sort.Search(excludedsLen, func(i int) bool {
				return excludeds[i] == el
			})
			if idx > -1 {
				continue
			}
		}
		if el.Intersects(x, y) {
			intersectElements = append(intersectElements, el)
		}
	}
	return intersectElements
}

//func (this *Layer) IsLayerGraphicsEngine() bool {
//	return this.layerGraphicsEngine
//}

//func (this *Layer) SetGraphicsEngine(ge IGraphicsEngine, isLayerGraphicsEngine bool) {
func (this *Layer) SetGraphicsEngine(ge IGraphicsEngine) {
	this.GraphicsEngine = ge
	ge.SetVisibleRegion(this.Self.(ILayer))

	//	if ge != nil {
	//		this.GraphicsEngine = ge
	//		this.layerGraphicsEngine = isLayerGraphicsEngine
	//		if isLayerGraphicsEngine {
	//			ge.SetVisibleRegion(this.Self.(ILayer))
	//		} else {
	//			ge.SetVisibleRegion(this.Page.GetSelf())
	//			this.SetDrawMode(DrawMode_All)
	//		}
	//	}
}
func (this *Layer) GetGraphicsEngine() IGraphicsEngine {
	return this.GraphicsEngine
}

func (this *Layer) SetFocusElement(e IElement) {
	this.focusElement = e
}

func (this *Layer) GetFocusElement() IElement {
	return this.focusElement
}

func (this *Layer) GetMouseHoveringElement() IElement {
	return this.mouseHoveringElement
}

func (this *Layer) setMouseHoveringElement(e IElement) {
	this.mouseHoveringElement = e
}

func (this *Layer) ClearMouseHoveringElement() {
	this.mouseHoveringElement = nil
}

func (this *Layer) GetMouseEventMode() MouseEventMode {
	return this.MouseEventMode
}

func (this *Layer) GetManualModeElement() IElement {
	return this.ManualModeElement
}

func (this *Layer) GetDrawPage() IDrawPage {
	return this.Page
}

func (this *Layer) SetDrawPage(page IDrawPage) {
	this.Page = page
}

func (this *Layer) SetBackground(c color.Color) {
	this.Background = c
}

func (this *Layer) GetBackground() color.Color {
	return this.Background
}

func (this *Layer) SetModified(b bool) {
	if this.Page != nil && b {
		this.Page.SetModified(b)
	}

	if this.IsModified() == b {
		return
	}
	this.ModifiedSupport.SetModified(b)

}

/**
 * 加入元素, 也可用于调整元素位置
 */
func (this *Layer) AddChild(e IElement, idx ...int) error {
	if _, ok := e.(ILayer); ok {
		return errors.New("element is layer, not supported")
	}
	err := this.addChild(e, idx...)
	if err == nil {
		e.SetLayer(this.Self.(ILayer))
		this.Self.(ILayer).SetModified(true)
		return nil
	}
	return err
}

//切换位置，根据元素低部坐标切换位置
func (this *Layer) SwitchPosition(e IElement) {
	layer := this.Self.(ILayer)
	lastVisibleEls := layer.GetLastVisibleRegionElements()
	if len(lastVisibleEls) < 2 {
		return
	}

	//	intersectElements := lastVisibleEls
	intersectElements := make([]IElement, len(lastVisibleEls))
	copy(intersectElements, lastVisibleEls)

	NewElementSortor(intersectElements, func(elements []IElement, i, j int) bool {
		return elements[i].GetBoundRect().Max.Y < elements[j].GetBoundRect().Max.Y
	}).Sort()

	iels := intersectElements
	e = iels[0]
	for _, el := range iels[1:] {
		//println(e.GetId(), e.Y(), el.Y(), e.Y()+e.Height(), el.Y()+el.Height(), layer.GetChildIndex(e), layer.GetChildIndex(el), e.GetOrderZ(), el.GetOrderZ())
		//println(e.GetId(), e.GetBoundRect().Max.Y, layer.GetChildIndex(e), e.GetOrderZ())
		this.adjustPosition(el, e, true)
		e = el
	}
}

/**
 * 将选中元素移动到指定元素的上面或下面
 */
func (this *Layer) adjustPosition(el IElement, target IElement, up bool) {
	if el == nil {
		return
	}
	layer := this.Self.(ILayer)
	if up {
		layer.Adjust(el, 1, target)
	} else {
		layer.Adjust(el, -1, target)
	}
	//	el.RedrawAll()
}

/**
 * 调整元素位置,target为目标参照元素，indexBy小于等于0，调整到在目标元素之前
 * （底层方向），大于0调整到目标元素之后(顶层方向), 如果不指定目标参照元素，默认为自己
 */
func (this *Layer) Adjust(e IElement, indexBy int, target IElement) {
	layer := this.Self.(ILayer)
	if !layer.HasChild(e) {
		return
	}
	if target == nil {
		target = e
	}

	e.SetOrderZ(target.GetOrderZ() + indexBy)

	//	idx := layer.GetChildIndex(target)
	//	if idx != -1 {
	//		newIndex = idx + indexBy
	//	}
	//	this.Element.addChild(e, newIndex)
	//this.Self.(ILayer).AddChild(e, newIndex)

}
func (this *Layer) Adjust_old(e IElement, indexBy int, target IElement) {
	layer := this.Self.(ILayer)
	if !layer.HasChild(e) {
		return
	}
	newIndex := indexBy
	if target == nil {
		target = e
	}
	idx := layer.GetChildIndex(target)
	if idx != -1 {
		newIndex = idx + indexBy
	}
	//	this.Element.addChild(e, newIndex)
	this.Self.(ILayer).AddChild(e, newIndex)
}

//func (this *Layer) Adjust(e IElement, indexBy int, target IElement) {
//	if !this.elements.Contains(e) {
//		return
//	}
//	newIndex := indexBy
//	if target == nil {
//		target = e
//	}
//	idx := this.elements.GetIndex(target)
//	if idx != -1 {
//		newIndex = idx + indexBy
//	}
//	this.Element.AddChild(e, newIndex)
//}

func (this *Layer) RemoveChild(e IElement) bool {
	//removedChildren = append(removedChildren, e)

	//鼠标事件响应元素
	if e == this.mouseHoveringElement {
		this.mouseHoveringElement = nil
		//		this.Self.(ILayer).ClearMouseHoveringElement()
	}

	//当前焦点元素
	if this.focusElement == e {
		this.focusElement = nil
	}

	this.Self.(ILayer).SetModified(true)

	return this.Element.RemoveChild(e)
}

/**
 * 对ordereds中元素排序，排序规则与元素在层中的顺序一致
 * 返回排序后的元素
 */
//func (this *Layer) SortElements(ordereds []IElement) []IElement {
//	if len(ordereds) == 0 {
//		return ordereds
//	}

//	tmpOrdereds := NewElements()
//	tmpOrdereds.Adds(ordereds...)
//	tmpOrdereds.Sort()

//	retrs := make([]IElement, 0)
//	//for _, e := range this.elements.GetElements() {
//	for _, e := range this.Self.(ILayer).GetChildren() {
//		if tmpOrdereds.Contains(e) {
//			retrs = append(retrs, e)
//		}
//	}

//	return retrs
//}

//func (this *Layer) Clear() {
//	//this.SetModified(true)
//	this.elements.Clear()
//	this.Self.(ILayer).SetModified(true)
//}

/**
 * 返回指定坐标 x,y上的最顶层可视并能响应事件的元素, 如果没有，返回NULL
 *
 * @param x X坐标
 * @param y Y坐标
 * @param excluded 被排除的元素
 *
 */
func (this *Layer) findByEventEnabled(x, y int) IElement {
	var e IElement
	elements := this.Self.(ILayer).GetChildren()
	elements.Sort()
	elements.ForEachLast(func(i int, el IElement) bool {
		//log.Println(i, el.GetId())
		if el != nil && el.IsEventEnabled() && el.IsVisible() && el.Intersects(x, y) {
			e = el
			return false
		}
		return true
	})
	return e

	//	els := this.Self.(ILayer).GetChildren() //this.elements.GetElements()
	//	for i := len(els) - 1; i >= 0; i-- {
	//		e := els[i]
	//		if e != nil && e.IsEventEnabled() && e.IsVisible() && e.Intersects(x, y) {
	//			return e
	//		}
	//	}
	//	return nil
}

/**
 * 返回指定坐标 x,y上的最顶层可视元素, 如果没有，返回nil
 *
 * @param x X坐标
 * @param y Y坐标
 * @param excluded 被排除的元素
 *
 */
func (this *Layer) Find(x, y int, excluded IElement) IElement {
	var e IElement
	elements := this.Self.(ILayer).GetChildren()
	elements.Sort()
	elements.ForEachLast(func(i int, el IElement) bool {
		if el != excluded && el.IsVisible() && el.Intersects(x, y) {
			e = el
			return false
		}
		return true
	})
	return e
	//	els := this.Self.(ILayer).GetChildren()
	//	for i := len(els) - 1; i >= 0; i-- {
	//		e := els[i]
	//		if e != excluded && e.IsVisible() && e.Intersects(x, y) {
	//			return e
	//		}
	//	}

}

//func (this *Layer) GetElements() []IElement {
//	return this.elements.GetElements()
//}

//func (this *Layer) GetIndex(e IElement) int {
//	return this.elements.GetIndex(e)
//}

func (this *Layer) GetSelection() ISelection {
	return this.Selection
}

///**
// * 设置焦点
// */
//func (this *Layer) SetFocus() {
//	this.Page.SetFocusLayer(this)
//}

/**
 * 检查是否焦点层
 */
//func (this *Layer) IsFocus() bool {
//	return this.Page.GetFocusLayer() == this
//}

func (this *Layer) SetDrawMode(dm DrawMode) {
	this.DrawMode = dm
}

func (this *Layer) GetDrawMode() DrawMode {
	return this.DrawMode
}

func (this *Layer) SetScrollMode(sm ScrollMode) {
	this.scrollMode = sm
}

func (this *Layer) GetScrollMode() ScrollMode {
	return this.scrollMode
}

//func (this *Layer) SetEventEnabled(e bool) {
//	this.EventEnabled = e
//}

//func (this *Layer) IsEventEnabled() bool {
//	return this.EventEnabled
//}

/**
 * 跟踪响应事件
 */
func (this *Layer) TrackEvent(e event.IEvent) bool {
	if !this.eventEnabled {
		return false
	}

	focusElement := this.focusElement
	if IsMouseEvent(e.GetType()) {
		//如果是鼠标事件
		ret := this.handleMouseEvent(e.(IMouseEvent))
		if e.GetType() == MOUSE_PRESS_EVENT_TYPE {
			mhoveringElement := this.Self.(ILayer).GetMouseHoveringElement()
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
	} else { //是键盘或其它事件
		if this.Self.(ILayer).FireEvent(e) {
			return true
		}
		//向层中当前焦点元素下发事件
		if focusElement != nil && focusElement.IsEventEnabled() {
			return focusElement.TrackEvent(e)
		}
		return false
	}
}

/**
 * 处理鼠标事件
 */
func (this *Layer) handleMouseEvent(me IMouseEvent) bool {

	x, y := me.X(), me.Y()

	//设置鼠标悬停状态
	this.mouseHovering = true

	//在本层中分发鼠标事件
	this.Self.(ILayer).FireEvent(me)

	//向层中元素转发鼠标事件
	switch this.MouseEventMode {
	case MEventMode_Hovering:
		e := this.findByEventEnabled(x, y)
		mhoveringElement := this.Self.(ILayer).GetMouseHoveringElement()
		if mhoveringElement != nil && mhoveringElement != e {
			mhoveringElement.fireMouseLeaveEvent(me) //离开元素
			//this.MouseHoveringElement = nil
			this.Self.(ILayer).ClearMouseHoveringElement()
			mhoveringElement = nil
		}
		if e != nil {
			if mhoveringElement != e {
				mhoveringElement = e
				this.setMouseHoveringElement(mhoveringElement)
				mhoveringElement.fireMouseEnterEvent(me)
			}
			return this.elementTrackEvent(mhoveringElement, me)
		}
	case MEventMode_Top:
		els := this.Self.(ILayer).GetChildren() //this.elements.GetElements()
		if els.Size() > 0 {
			mhoveringElement := this.Self.(ILayer).GetMouseHoveringElement()
			//e := els[len(els)-1]
			e := els.At(els.Size() - 1)
			if e != nil {
				if e.IsEventEnabled() && e.IsVisible() && e.Intersects(x, y) {
					if mhoveringElement != e {
						mhoveringElement = e
						this.setMouseHoveringElement(mhoveringElement)
						mhoveringElement.fireMouseEnterEvent(me)
					}
					return this.elementTrackEvent(mhoveringElement, me)
				} else if mhoveringElement != nil {
					mhoveringElement.fireMouseLeaveEvent(me) //离开元素
					this.Self.(ILayer).ClearMouseHoveringElement()
				}
			}
		}
		break
	case MEventMode_Manual:
		mhoveringElement := this.Self.(ILayer).GetMouseHoveringElement()
		if this.ManualModeElement != nil &&
			this.ManualModeElement.IsEventEnabled() &&
			this.ManualModeElement.IsVisible() &&
			this.ManualModeElement.Intersects(x, y) {
			if mhoveringElement != this.ManualModeElement {
				mhoveringElement = this.ManualModeElement
				this.setMouseHoveringElement(mhoveringElement)
				mhoveringElement.fireMouseEnterEvent(me)

			}
			return this.elementTrackEvent(mhoveringElement, me)
		} else {
			if mhoveringElement != nil {
				mhoveringElement.fireMouseLeaveEvent(me) //离开元素
				this.Self.(ILayer).ClearMouseHoveringElement()
			}
		}
	}
	return false
}

func (this *Layer) elementTrackEvent(el IElement, me event.IEvent) bool {
	return el.TrackEvent(me)

	//	if me.GetType() == MOUSE_WHEEL_EVENT_TYPE {
	//		nme := *(me.(*WheelEvent))
	//		x, y := el.GetCoordinate()
	//		nme.SetX(nme.X() - int(x))
	//		nme.SetY(nme.Y() - int(y))
	//		nme.Source = el
	//		return el.TrackEvent(&nme)
	//	} else {
	//		nme := *(me.(*MouseEvent))
	//		x, y := el.GetCoordinate()
	//		nme.SetX(nme.X() - int(x))
	//		nme.SetY(nme.Y() - int(y))
	//		nme.Source = el
	//		return el.TrackEvent(&nme)
	//	}
}

/**
 * 发送离开层事件
 */
func (this *Layer) fireMouseLeaveEvent(e IMouseEvent) bool {
	if !this.mouseHovering {
		return true
	}
	//清除鼠标悬停状态
	if this.Page.GetMouseHoveringLayer() == this.Self {
		this.mouseHovering = false
		this.Page.ClearMouseHoveringLayer()
	}
	//this.MouseHovering = false
	this.Page.ClearMouseHoveringLayer()

	//如果鼠标最近一次悬停的元素存在，触发离开元素方法
	mhoveringElement := this.Self.(ILayer).GetMouseHoveringElement()
	if mhoveringElement != nil && mhoveringElement.MouseIsHovering() {
		mhoveringElement.fireMouseLeaveEvent(e)
	}

	mLeaveEvent := NewMouseEvent(MOUSE_LEAVE_EVENT_TYPE, e.GetSource(), e.X(), e.Y(), e.GetButtons(), e.GetModifier())
	//	mLeaveEvent.x = e.X()
	//	mLeaveEvent.y = e.Y()
	mLeaveEvent.KeySequence = e.GetKeySequence()

	//向在此层上监听的对象发送鼠标离开事件
	return this.Self.(ILayer).FireEvent(mLeaveEvent)
}

/**
 * 返回最近一次绘制元素的数量
 */
func (this *Layer) GetDrawElementsCount() int {
	return this.drawedCount
}

/**
 * 返回最近一次绘制层使用时间
 */
func (this *Layer) GetDrawSpendTime() int64 {
	return this.spendTime
}

func (this *Layer) UnionDrawRegion(rect *image.Rectangle) {
	//println("Layer.UnionDrawRegion.rect", rect.String(), rect.Dx(), rect.Dy(), len(this.clipedRects))
	this.clipedRects = append(this.clipedRects, rect)
	this.Self.(ILayer).SetModified(true)
}

func (this *Layer) ClipDrawRegion(ge IGraphicsEngine) {
	if this.clipedRegion == nil {
		this.clipedRegion = ge.NewRegion()
		//		fmt.Println("this.clipedRegion:", this.clipedRegion)
	} else {
		this.clipedRegion.Clear()
	}

	//println("Layer.ClipDrawRegion.clipedRects", len(this.clipedRects))
	for _, r := range this.clipedRects {
		//println("Layer.ClipDrawRegion.r", r.Dx(), r.Dy())
		this.clipedRegion.UnionWithIntersect(*r)
	}
	ge.SetClip(this.clipedRegion)

	if this.Background != nil {
		ge.SetFillColor(this.Background)
		ge.Clear()
	}
}

func (this *Layer) ResetClipedDrawRegion(ge IGraphicsEngine) {
	ge.ResetClip()
	if len(this.clipedRects) > 0 {
		this.clipedRects = this.clipedRects[0:0]
	}

	if this.clipedRegion != nil {
		this.clipedRegion.Release()
		this.clipedRegion = nil
	}
}

func (this *Layer) GetIntersectionsByElement(el IElement, excludeds *ElementSortor, appendeds *ElementSortor) []IElement {
	return this.Self.(ILayer).GetIntersections(func(e IElement) bool { return e.IntersectsElement(el) }, excludeds, appendeds)
}

func (this *Layer) GetIntersectionsByElements(els []IElement, excludeds *ElementSortor, appendeds *ElementSortor) []IElement {
	return this.Self.(ILayer).GetIntersections(func(e IElement) bool { return e.IntersectsElements(els...) }, excludeds, appendeds)
}

func (this *Layer) GetIntersectionsByRect(rect *image.Rectangle, excludeds *ElementSortor, appendeds *ElementSortor) []IElement {
	return this.Self.(ILayer).GetIntersections(func(e IElement) bool { return e.IntersectsWith(rect) }, excludeds, appendeds)
}

func (this *Layer) GetIntersectionsByRects(rects []*image.Rectangle, excludeds *ElementSortor, appendeds *ElementSortor) []IElement {
	return this.Self.(ILayer).GetIntersections(func(e IElement) bool { return e.IntersectsWiths(rects...) }, excludeds, appendeds)
}

/**
 * 返回选择区域中除<code>excludeds</code>以外的所有相交的元素, 在<code>v</code>中返回相交的元素，<code>rects</code>为指定的区域集。
 * <code>excludeds</code>更优先于<code>appendeds</code>。
 *
 *
 * @param excludeds 为被排除的元素
 * @param appendeds 中为附加的已经排序的元素，即使appendeds中的元素没有在相交区域，也将此附加元素在<code>v</code>中返回，
 *                  这样的目的主要是为了能返回拥有正确层次的元素
 *
 */
func (this *Layer) GetIntersections(intersects func(e IElement) bool, excludeds *ElementSortor, appendeds *ElementSortor) []IElement {
	els := this.Self.(ILayer).GetChildren() //this.elements.GetElements()
	retes := make([]IElement, 0)
	els.ForEach(func(i int, e IElement) bool {

		if !e.IsVisible() || excludeds != nil && excludeds.Contains(e) {
			return true
		}
		if appendeds != nil && appendeds.Contains(e) || intersects(e) {
			retes = append(retes, e)
		}
		return true
	})

	//	for _, e := range els {
	//		if !e.IsVisible() || excludeds != nil && excludeds.Contains(e) {
	//			continue
	//		}
	//		if appendeds != nil && appendeds.Contains(e) || intersects(e) {
	//			retes = append(retes, e)
	//		}
	//	}
	return retes
}

func GetIntersectionsByElement(srcs []IElement, el IElement, excludeds IElements, appendeds IElements) []IElement {
	return GetIntersections(srcs, func(e IElement) bool { return e.IntersectsElement(el) }, excludeds, appendeds)
}

func GetIntersectionsByElements(srcs []IElement, els []IElement, excludeds IElements, appendeds IElements) []IElement {
	return GetIntersections(srcs, func(e IElement) bool { return e.IntersectsElements(els...) }, excludeds, appendeds)
}

func GetIntersectionsByRect(srcs []IElement, rect *image.Rectangle, excludeds IElements, appendeds IElements) []IElement {
	return GetIntersections(srcs, func(e IElement) bool { return e.IntersectsWith(rect) }, excludeds, appendeds)
}

func GetIntersectionsByRects(srcs []IElement, rects []*image.Rectangle, excludeds IElements, appendeds IElements) []IElement {
	return GetIntersections(srcs, func(e IElement) bool { return e.IntersectsWiths(rects...) }, excludeds, appendeds)
}

func GetIntersections(srcs []IElement, intersects func(e IElement) bool, excludeds IElements, appendeds IElements) []IElement {
	//	els := this.elements.GetElements()
	retes := make([]IElement, 0)
	for _, e := range srcs {

		if !e.IsVisible() || excludeds != nil && excludeds.Contains(e) {
			continue
		}
		if appendeds != nil && appendeds.Contains(e) || intersects(e) {
			retes = append(retes, e)
		}
	}
	return retes
}

//func (this *Layer) setRendering(b bool) {
//	if b {
//		atomic.StoreInt32(&this.rendering, 1)
//	} else {
//		atomic.StoreInt32(&this.rendering, 0)
//	}
//}

//func (this *Layer) IsRendering() bool {
//	return atomic.LoadInt32(&this.rendering) == 1
//}

//func (this *Layer) getScrollRects(vr, lvr *image.Rectangle) (rightvr, bottomvr, leftvr, topvr *image.Rectangle) {
//	rightvr = &image.Rectangle{Min: image.Point{lvr.Max.X, vr.Min.Y}, Max: image.Point{vr.Max.X, vr.Max.Y}}
//	bottomvr = &image.Rectangle{Min: image.Point{vr.Min.X, vr.Max.Y}, Max: image.Point{lvr.Max.X, vr.Max.Y}}
//	leftvr = &image.Rectangle{Min: image.Point{lvr.Min.X, lvr.Min.Y}, Max: image.Point{vr.Min.X, lvr.Max.Y}}
//	topvr = &image.Rectangle{Min: image.Point{vr.Min.X, lvr.Min.Y}, Max: image.Point{lvr.Max.X, vr.Min.Y}}
//	return

//}

func (this *Layer) scroll(ge IGraphicsEngine, vr, lvr *image.Rectangle) {
	if *vr == *lvr {
		//fmt.Println("Layer.scroll *vr == *lvr:", this.id)
		return
	}
	//卡马克卷轴法
	ivr := vr.Intersect(*lvr)

	ge.SetFillColor(this.Background)
	scrollRegions := make([]*image.Rectangle, 0, 2)
	if lvr.Min.X < vr.Min.X {
		if lvr.Min.Y < vr.Min.Y {
			rightvr := image.Rectangle{Min: image.Point{lvr.Max.X, vr.Min.Y}, Max: image.Point{vr.Max.X, vr.Max.Y}}
			bottomvr := image.Rectangle{Min: image.Point{vr.Min.X, vr.Max.Y}, Max: image.Point{lvr.Max.X, vr.Max.Y}}
			scrollRegions = append(scrollRegions, &rightvr)
			scrollRegions = append(scrollRegions, &bottomvr)
			ge.CopyBuffer(0, 0, INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)

		} else if lvr.Min.Y > vr.Min.Y {
			rightvr := image.Rectangle{Min: image.Point{lvr.Max.X, vr.Min.Y}, Max: image.Point{lvr.Max.X, vr.Max.Y}}
			topvr := image.Rectangle{Min: image.Point{vr.Min.X, vr.Min.Y}, Max: image.Point{lvr.Max.X, lvr.Min.Y}}
			scrollRegions = append(scrollRegions, &rightvr)
			scrollRegions = append(scrollRegions, &topvr)
			ge.CopyBuffer(0, INT(lvr.Min.Y-vr.Min.Y), INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)

		} else {
			rightvr := image.Rectangle{Min: image.Point{lvr.Max.X, vr.Min.Y}, Max: image.Point{vr.Max.X, vr.Max.Y}}
			scrollRegions = append(scrollRegions, &rightvr)
			ge.CopyBuffer(0, 0, INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)
		}
	} else {
		if lvr.Min.Y < vr.Min.Y {
			leftvr := image.Rectangle{Min: image.Point{vr.Min.X, vr.Min.Y}, Max: image.Point{lvr.Min.X, vr.Max.Y}}
			bottomvr := image.Rectangle{Min: image.Point{lvr.Min.X, lvr.Max.Y}, Max: image.Point{vr.Max.X, vr.Max.Y}}
			scrollRegions = append(scrollRegions, &leftvr)
			scrollRegions = append(scrollRegions, &bottomvr)
			ge.CopyBuffer(INT(lvr.Dx()-ivr.Dx()), 0, INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)

		} else if lvr.Min.Y > vr.Min.Y {
			leftvr := image.Rectangle{Min: image.Point{vr.Min.X, vr.Min.Y}, Max: image.Point{lvr.Min.X, vr.Max.Y}}
			topvr := image.Rectangle{Min: image.Point{lvr.Min.X, vr.Min.Y}, Max: image.Point{vr.Max.X, lvr.Min.Y}}
			scrollRegions = append(scrollRegions, &leftvr)
			scrollRegions = append(scrollRegions, &topvr)
			ge.CopyBuffer(INT(lvr.Dx()-ivr.Dx()), INT(lvr.Dy()-ivr.Dy()), INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)
		} else {
			leftvr := image.Rectangle{Min: image.Point{vr.Min.X, vr.Min.Y}, Max: image.Point{lvr.Min.X, vr.Max.Y}}
			scrollRegions = append(scrollRegions, &leftvr)
			ge.CopyBuffer(INT(lvr.Dx()-ivr.Dx()), 0, INT(ivr.Min.X-lvr.Min.X), INT(ivr.Min.Y-lvr.Min.Y), INT(ivr.Dx()), INT(ivr.Dy()), true)
		}
	}

	//@TODO 优化来只取可视区域元素
	var vrels []IElement
	if len(scrollRegions) > 1 {
		vrels = this.Self.(ILayer).GetIntersectionsByRect(vr, nil, nil)
	} else {
		vrels = this.Self.(ILayer).GetChildren().CloneElements()
	}
	//	x, _ := this.GetVisibleRegionCoordinate()
	//	if x != vr.Min.X {
	//		fmt.Println(x, vr.Min.X)
	//	}

	ge.SwapBuffers()

	for _, sr := range scrollRegions {
		//fmt.Println("Layer.scroll sr:", this.GetId(), i, *sr, sr.Dx(), sr.Dy(), this.GetId())
		//if sr.Dx() > 0 && sr.Dy() > 0 {
		this.drawScrollRegion(ge, sr, vrels, vr)
		//ge.SetStrokeColor(color.NRGBA{255, 255, 255, 255})
		//ge.DrawRectangle(sr)
		//} else {
		//	fmt.Println("Layer.scroll sr:", this.GetId(), *sr, sr.Dx(), sr.Dy(), this.GetId())
		//}
	}
}

//vrels中为有可能被重画的元素, sr为与元素相交区域
func (this *Layer) drawScrollRegion(ge IGraphicsEngine, sr *image.Rectangle, vrels []IElement, vr *image.Rectangle) {

	sr_ := sr.Inset(-1)
	sr = &sr_

	ge.SetClipRect(sr)
	defer ge.ResetClip()
	if this.Background != nil {
		ge.SetFillColor(this.Background)
		ge.Clear()
	}

	for _, el := range vrels {
		if el.IsVisible() && el.IntersectsWith(sr) {
			//el.SetModified(false)
			el.Draw(ge)
		}
	}
}

func (this *Layer) Draw(ge IGraphicsEngine) {
	timespender := NewTimespender("Layer.Draw, " + this.GetId())
	defer func() {
		this.spendTime = timespender.Spendtime()
		//timespender.Print()
	}()

	layer := this.Self.(ILayer)
	//	if !layer.IsVisible() {
	//		return
	//	}

	scrollable := this.visibleRegionModified
	vr, lvr := layer.GetVisibleRegion(true)
	if this.scrollMode == ScrollMode_CAMAC && scrollable && !layer.IsModified() {
		//fmt.Println("Layer.Draw.scrollable", scrollable, this.GetId())
		//图层卷动
		this.scroll(ge, vr, lvr)
		return
	}
	//	if this.scrollMode == ScrollMode_CAMAC {
	//		fmt.Println("Layer.Draw.scrollable", scrollable, this.GetId(), layer.IsModified())
	//	}
	//	if !layer.IsModified() {
	//		return
	//	}

	//	scrollable := this.visibleRegionModified
	//	vr, _ := layer.GetVisibleRegion(true)
	if !layer.IsModified() && !scrollable {
		//		if this.id == "sky" || this.id == "earth" {
		//			fmt.Println("Layer.Draw", this.id)
		//		}
		//如果层没有被修改过
		return
	}

	//	drawTimespender := NewTimespender("Draw  aaa ")
	childrenCount := this.GetChildrenCount()
	if childrenCount > 3000 {
		childrenCount = 3000
	}
	drawedEls := make([]IElement, 0, childrenCount)

	for _, el := range this.children.CloneElements() {
		if !el.IsDestroyed() {
			if el.IntersectsWith(vr) {
				//if el.GetBoundRect().Intersect(vr) != image.ZR {
				el.SetActive(true)
				drawedEls = append(drawedEls, el)
			} else {
				el.SetActive(false)
			}
		} else {
			layer.RemoveChild(el)
		}
	}
	//	if this.GetLayerType() == LayerTypeActive {
	//		drawTimespender.Print()
	//	}

	if scrollable && this.DrawMode != DrawMode_All {
		panic(fmt.Sprintf("visible region Adjusted, can not use the DrawMode_Region mode, id:%v, tag:%v", this.GetId(), this.GetTag()))
	}

	this.sortor.SetElements(drawedEls)
	this.sortor.Sort()

	//	drawTimespender = NewTimespender("Draw  ccc ")
	layer.SetModified(false)
	if this.DrawMode == DrawMode_All || scrollable || layer.IsRedrawAll() {
		//绘制所有
		this.drawAll(ge, drawedEls, scrollable)
	} else {
		//区域绘制
		this.drawRegion(ge, drawedEls)
	}

	//ge.SwapBuffers()

	//	if this.GetLayerType() == LayerTypeActive {
	//		drawTimespender.Print()
	//	}
	layer.SetLastVisibleRegionElements(drawedEls)

	layer.firePaintEvent()
}

func (this *Layer) Draw_old(ge IGraphicsEngine) {

	layer := this.Self.(ILayer)
	if !layer.IsVisible() {
		return
	}

	scrollable := this.visibleRegionModified //layer.IsVisibleRegionModified()
	vr, _ := layer.GetVisibleRegion(true)
	//if !layer.IsVisible() || !(layer.IsModified() || vrm) {
	if scrollable {
		//执行地图卷动

	}
	if !layer.IsModified() && !scrollable {
		//如果层没有被修改过
		return
	}

	//	drawTimespender := NewTimespender("Draw  aaa ")
	drawedEls := make([]IElement, 0, 3000)
	//	this.Self.(ILayer).GetChildren().ForEach(func(i int, el IElement) bool {
	//		if !el.IsDestroyed() {
	//			if el.IntersectsWith(&vr) {
	//				el.SetActive(true)
	//				drawedEls = append(drawedEls, el)
	//			} else {
	//				el.SetActive(false)
	//			}
	//		} else {
	//			log.Println("Layer.Draw IsDestroyed ", el)
	//			//if this.Exist(el.GetId()) {
	//			this.RemoveChild(el)
	//			//}
	//		}
	//		return true
	//	})

	//elements := this.Self.(ILayer).GetChildren().CloneElements()
	for _, el := range this.children.CloneElements() {
		if !el.IsDestroyed() {
			if el.IntersectsWith(vr) {
				//if el.GetBoundRect().Intersect(vr) != image.ZR {
				el.SetActive(true)
				drawedEls = append(drawedEls, el)
			} else {
				el.SetActive(false)
			}
			//} else {
			//	log.Println("Layer.Draw IsDestroyed ", el)
			//this.RemoveChild(el)
		}
	}
	//	if this.GetLayerType() == LayerTypeActive {
	//		drawTimespender.Print()
	//	}

	if scrollable && this.DrawMode != DrawMode_All {
		panic(fmt.Sprintf("visible region Adjusted, can not use the DrawMode_Region mode, id:%v, tag:%v", this.GetId(), this.GetTag()))
	}

	this.sortor.SetElements(drawedEls)
	this.sortor.Sort()

	//	drawTimespender = NewTimespender("Draw  ccc ")
	layer.SetModified(false)
	if this.DrawMode == DrawMode_All || scrollable || layer.IsRedrawAll() {
		//绘制所有
		this.drawAll(ge, drawedEls, scrollable)
	} else {
		//区域绘制
		this.drawRegion(ge, drawedEls)
	}
	//	if this.GetLayerType() == LayerTypeActive {
	//		drawTimespender.Print()
	//	}
	layer.SetLastVisibleRegionElements(drawedEls)

	layer.firePaintEvent()
}

func (this *Layer) drawAll(ge IGraphicsEngine, drawedEls []IElement, vrm bool) {
	//layer := this.Self.(ILayer)
	this.drawedCount = 0
	//	stime := time.Now().UnixNano()
	//	if this.IsLayerGraphicsEngine() {
	if this.Background != nil {
		ge.SetFillColor(this.Background)
		ge.Clear()
	}
	//	}

	for _, el := range drawedEls { //this.elements.GetElements() {
		if el.IsVisible() && !el.IsDestroyed() { // && layer.EnableDisplay(el) {
			if el.IsModified() || vrm {
				el.CreatePath()
			}
			//println("Layer) drawAll:", el.GetId(), el.GetOrderZ())
			//el.CreatePath()
			this.drawElement(el, ge)
			this.drawedCount++
		}
	}
	//发送绘制层事件
	//layer.FireEvent(this.paintEvent)
	//	this.spendTime = time.Now().UnixNano() - stime
}

func (this *Layer) drawRegion(ge IGraphicsEngine, drawedEls []IElement) {
	layer := this.Self.(ILayer)
	this.drawedCount = 0
	//	stime := time.Now().UnixNano()

	//过滤出需要重绘的元素
	modifiedElements := make([]IElement, 0) //NewElements()
	redrawedElements := make([]IElement, 0) //NewElements()

	for _, e := range drawedEls {
		if !e.IsVisible() {
			continue
		}
		if e.IsModified() {
			e.CreatePath()
			modifiedElements = append(modifiedElements, e)
			layer.UnionDrawRegion(e.GetBoundRect(e.GetClipRegionAdjustValue()))
		} else if e.IsRedraw() {
			redrawedElements = append(redrawedElements, e)
		}
	}

	if len(this.clipedRects) > 0 {
		//锁定需要重绘的区域
		layer.ClipDrawRegion(ge)
		defer layer.ResetClipedDrawRegion(ge)

		if len(modifiedElements) == 0 {
			//绘制仅重画自己的元素
			for _, el := range redrawedElements {
				el.SetModified(false)
				el.SetRedraw(false)
				el.Draw(ge)
				this.drawedCount++
			}

		} else {
			//两者都重画
			redrawedElements = append(redrawedElements, modifiedElements...)
			sortor := NewElementSortor(redrawedElements, nil)
			sortor.Sort()
			allRedrawedEls := layer.GetIntersectionsByRects(this.clipedRegion.GetRects(), nil, sortor)
			for _, e := range allRedrawedEls {
				if e.IsVisible() {
					this.drawElement(e, ge)
					this.drawedCount++
				}
			}

		}
		//layer.ResetClipedDrawRegion(ge)
	}

	//	this.spendTime = time.Now().UnixNano() - stime
}

func (this *Layer) drawElement(el IElement, ge IGraphicsEngine) {
	el.SetModified(false)
	el.SetRedraw(false)
	el.Draw(ge)
	el.firePaintEvent()
}

///**
// * 检查此元素是否在可视区
// */
//func (this *Layer) EnableDisplay(el IElement) bool {
//	return true
//}

func (this *Layer) IsRedrawAll() bool {
	if this.RedrawAllFlag {
		this.RedrawAllFlag = false
		return true
	}
	return false
}

func (this *Layer) SetRedrawAll() {
	this.RedrawAllFlag = true
}

func (this *Layer) PrepareTransform(x, y int) {
	this.MoveSupport.PrepareTransform(x, y)
}

//func (this *Layer) GetElementsCount() int {
//	return this.elements.Size()
//}

func (this *Layer) Init() {

}
