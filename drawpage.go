package eui

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"sync/atomic"
	"time"

	//	"github.com/google/gxui"
	"github.com/tryor/commons/event"
	"github.com/tryor/gdiplus"
)

type IDrawPage interface {
	IEventSupport

	IsModified() bool
	SetModified(b bool)
	IsDrawing() bool   //是否正在绘制图层
	IsRendering() bool //是否正在渲染图层和绘制图层
	IsDestroyed() bool

	GetFocusLayer() ILayer
	SetFocusLayer(l ILayer)
	GetMouseHoveringLayer() ILayer
	ClearMouseHoveringLayer()
	SetGraphicsEngine(ge IGraphicsEngine)
	GetGraphicsEngine() IGraphicsEngine
	SortLayers()
	ClearLayersSortFlag()
	GetLayers() []ILayer
	GetLayersCount() int
	GetGeometry() *image.Rectangle
	SetGeometry(r *image.Rectangle)
	X() int
	Y() int
	Width() int
	Height() int
	SetVisibleRegion(x, y, w, h int)
	GetVisibleRegion() image.Rectangle
	Draw()
	Render(...bool)
	SetRenderEnable(b bool)
	//isFinish false:渲染开始， true表示渲染完成
	OnRender(func(page IDrawPage, isFinish bool)) EventSubscription

	AddLayer(l ILayer, idx ...int) error
	RemoveLayer(l ILayer) bool
	//HandleEvent(e event.IEvent) bool
	GetCursor() ICursor
	//RefreshCursor() bool
	TrackEvent(e event.IEvent) bool
	TrackMouseEvent(t event.Type, source interface{}, x, y int, buttons MButton, modifier KeyboardModifier, kseq *KeySequence)

	//	InitLayerGraphicsEngines()
	Init() error
	//	GetSelf() IDrawPage

	GetVisibleRegionCoordinate() (x, y int)
	//GetVisibleRegionSnapshotCoordinate() (x, y int)
	//MakeVisibleRegionSnapshot()

	AdjustVisibleRegion(dx, dy REAL)

	SetShowStatInfoEnable(b bool)
	GetDrawSpendTime() int64
	GetFPS(statTimes int) (float64, bool)
	ShowFPS(statTimes int)

	StartStatTimer()
	StopStatTimer()
}

type DrawPage struct {
	ModifiedSupport
	EventSupport

	Self   IDrawPage
	layers *Elements

	FocusLayer         ILayer //当前焦点层
	MouseHoveringLayer ILayer //鼠标悬停层

	MouseEventMode  MouseEventMode //鼠标事件响应模式
	ManualModeLayer ILayer         //MEventMode_Manual模式下指定层

	Geometry      image.Rectangle //画布位置
	VisibleRegion image.Rectangle //可视矩形区

	GraphicsEngine IGraphicsEngine //图形引擎

	modifiedEvent *ModifiedEvent //定义改变被修改状态事件
	paintEvent    *PaintEvent
	onRenderEvent Event

	drawing       bool
	renderEnable  bool //仅当renderEnable为true才能继续进行渲染，否则将被暂停渲染
	rendering     int32
	lockRendering int32
	destroyed     bool //指示是否已经被销毁

	showStatInfoEnable bool
	spendTime          int64 //最近一次绘制用时，单位：纳秒
	statRenderCount    int
	statRenderTime     int64
	statFps            float64
	actualFps          int //实际最近一秒渲染帧数
	statActualFps      int32
	actualFpsTimer     *Ticker
}

func NewDrawPage(rect image.Rectangle) *DrawPage {
	page := &DrawPage{layers: NewElements(), MouseEventMode: MEventMode_Hovering}
	//	page.Dispatcher = event.NewDispatcher()
	page.Self = page
	page.Geometry = rect

	page.modifiedEvent = NewModifiedEvent(page, &page.ModifiedSupport)
	page.paintEvent = NewPaintEvent(page)

	page.renderEnable = true

	return page
}

func (this *DrawPage) Destroy() {
	if this.destroyed {
		log.Println("DrawPage is already destroyed")
		return
	}

	this.destroyed = true

	this.Self.FireEvent(NewDestroyEvent(this.Self))

	for _, l := range this.Self.GetLayers() {
		l.Destroy()
	}
	//this.layers.Destroy()
}

func (this *DrawPage) IsDestroyed() bool {
	return this.destroyed
}

//启动统计用Timer
func (this *DrawPage) StartStatTimer() {
	if this.actualFpsTimer != nil {
		return
	}
	this.actualFpsTimer = NewTicker(time.Second, func(t time.Time) {
		this.actualFps = int(atomic.LoadInt32(&this.statActualFps))
		atomic.StoreInt32(&this.statActualFps, 0)
		//log.Println(this.actualFps)
	})
}

func (this *DrawPage) StopStatTimer() {
	if this.actualFpsTimer != nil {
		this.actualFpsTimer.Close()
		this.actualFpsTimer = nil
	}
}

func (this *DrawPage) AdjustVisibleRegion(dx, dy REAL) {
	this.adjustVisibleRegion(dx, dy)
}

func (this *DrawPage) adjustVisibleRegion(dx, dy REAL) {
	if dx == 0 && dy == 0 {
		return
	}
	fdx, fdy := dx, dy
	var scrollrateX, scrollrateY REAL
	for _, layer := range this.Self.(IDrawPage).GetLayers() {
		//fmt.Println("DrawPage) adjustVisibleRegion:", layer.GetId(), layer.GetBoundRect(), layer.GetVisibleRegion())
		vr := layer.GetVisibleRegionF()
		if int(vr.W) == layer.Width() && int(vr.H) == layer.Height() {
			//fmt.Println("DrawPage) adjustVisibleRegion:", layer.GetId())
			continue
		}
		scrollrateX = REAL(layer.GetScrollrateX())
		scrollrateY = REAL(layer.GetScrollrateY())
		if scrollrateX == 0.0 && scrollrateY == 0.0 {
			continue
		}
		if scrollrateX > 0 || scrollrateY > 0 {
			if scrollrateX > 0 {
				fdx = dx * scrollrateX
			} else {
				fdx = 0
			}
			if scrollrateY > 0 {
				fdy = dy * scrollrateY
			} else {
				fdy = 0
			}
			//if layer.GetId() == "sky" || layer.GetId() == "earth" {
			//	fmt.Printf("DrawPage.layer.scrollrate:%v, %v, %v, %0.2f, %0.2f, %0.2f, %0.2f\n", layer.GetId(), scrollrateX, scrollrateY, fdx, fdy, dx, dy)
			//}
			layer.AdjustVisibleRegion(fdx, fdy)
		}
	}
}

//func (this *DrawPage) adjustVisibleRegion(dx, dy int) {
//	fdx, fdy := float32(dx), float32(dy)
//	for _, layer := range this.Self.(IDrawPage).GetLayers() {
//		scrollrateX := layer.GetScrollrateX()
//		scrollrateY := layer.GetScrollrateY()
//		if scrollrateX > 0 || scrollrateY > 0 {
//			if scrollrateX > 0 {
//				dx = int(fdx * scrollrateX)
//			} else {
//				dx = 0
//			}
//			if scrollrateY > 0 {
//				dy = int(fdy * scrollrateY)
//			} else {
//				dy = 0
//			}
//			//fmt.Printf("DrawPage.layer.scrollrate:%v, %v, %v, %v, %v, %v\n", scrollrateX, scrollrateY, dx, dy, fdx, fdy)
//			layer.AdjustVisibleRegion(dx, dy)
//		}
//	}
//}

func (this *DrawPage) GetVisibleRegionCoordinate() (x, y int) {
	return -this.VisibleRegion.Min.X, -this.VisibleRegion.Min.Y
}

//func (this *DrawPage) GetVisibleRegionSnapshotCoordinate() (x, y int) {
//	return 0, 0
//}

//func (this *DrawPage) MakeVisibleRegionSnapshot() {
//	this.layers.ForEach(func(i int, layer IElement) bool {
//		layer.(ILayer).MakeVisibleRegionSnapshot()
//		return true
//	})
//}

/**
 * 返回当前焦点层
 */
func (this *DrawPage) GetFocusLayer() ILayer {
	return this.FocusLayer
}

/**
 * 设置焦点层
 */
func (this *DrawPage) SetFocusLayer(l ILayer) {
	if l != this.FocusLayer {
		if this.FocusLayer != nil {
			focusElement := this.FocusLayer.GetFocusElement()
			if focusElement != nil {
				focusElement.fireFocusEvent(false)
				this.FocusLayer.setFocusElement(nil)
				//this.FocusLayer.SetFocusElement(nil) //TODO at 20190712
			}
		}
		this.FocusLayer = l
	}
}

func (this *DrawPage) GetMouseHoveringLayer() ILayer {
	return this.MouseHoveringLayer
}

func (this *DrawPage) ClearMouseHoveringLayer() {
	this.MouseHoveringLayer = nil
}

func (this *DrawPage) SetGraphicsEngine(ge IGraphicsEngine) {
	this.GraphicsEngine = ge
	ge.SetVisibleRegion(this.Self.(IDrawPage))
}

func (this *DrawPage) GetGraphicsEngine() IGraphicsEngine {
	return this.GraphicsEngine
}

func (this *DrawPage) Init() error {
	if this.GraphicsEngine == nil {
		return errors.New("Graphics Engine is nil")
	}
	//	for _, layer := range this.layers.GetElements() {
	//		this.GraphicsEngine.GetLayerMerger().InitLayerGraphicsEngine(layer.(ILayer))
	//		layer.(ILayer).Init()
	//	}
	this.layers.ForEach(func(i int, layer IElement) bool {
		this.GraphicsEngine.GetLayerMerger().InitLayerGraphicsEngine(layer.(ILayer))
		layer.(ILayer).Init()
		return true
	})

	return nil
}

/**
 * 设置画布位置，(x,y)为坐标，(w, h)为宽度和高度
 */
//func (this *DrawPage) SetGeometry(x, y, w, h int) {
//	this.Geometry = image.Rect(x, y, x+w, y+h)
//}
func (this *DrawPage) SetGeometry(r *image.Rectangle) {
	this.Geometry = *r
}

func (this *DrawPage) GetGeometry() *image.Rectangle {
	return &this.Geometry
}

func (this *DrawPage) X() int {
	return this.Geometry.Min.X
}

func (this *DrawPage) Y() int {
	return this.Geometry.Min.Y
}

func (this *DrawPage) Width() int {
	return this.Geometry.Dx()
}

func (this *DrawPage) Height() int {
	return this.Geometry.Dy()
}

func (this *DrawPage) SetVisibleRegion(x, y, w, h int) {
	this.VisibleRegion = image.Rect(x, y, x+w, y+h)
}

func (this *DrawPage) GetVisibleRegion() image.Rectangle {
	return this.VisibleRegion
}

func (this *DrawPage) AddLayer(l ILayer, idx ...int) error {
	err := this.layers.Add(l, idx...)
	if err == nil {
		l.SetDrawPage(this.Self.(IDrawPage))
		this.Self.(IDrawPage).SetModified(true)
		return nil
	}
	return err
}

func (this *DrawPage) RemoveLayer(l ILayer) bool {
	if l == nil {
		return false
	}

	if this.layers.Remove(l) {
		//		l.RemoveListener(this)
		this.Self.(IDrawPage).SetModified(true)
		return true
	}
	return false
}

func (this *DrawPage) SortLayers() {
	this.layers.Sort()
}

func (this *DrawPage) ClearLayersSortFlag() {
	this.layers.ClearSortFlag()
}

func (this *DrawPage) GetLayers() []ILayer {
	layers := make([]ILayer, 0, this.layers.Size())
	//	for _, l := range this.layers.GetElements() {
	//		layers = append(layers, l.(ILayer))
	//	}
	this.layers.ForEach(func(i int, layer IElement) bool {
		layers = append(layers, layer.(ILayer))
		return true
	})
	return layers
}

func (this *DrawPage) GetLayersCount() int {
	return this.layers.Size()
}

func (this *DrawPage) SetModified(b bool) {
	//	if this.IsModified() == b {
	//		return
	//	}
	this.ModifiedSupport.SetModified(b)

	//处理事件方法
	if b {
		this.Self.(IEventSupport).FireEvent(this.modifiedEvent)
	}
}

//func (this *DrawPage) HandleEvent(e event.IEvent) bool {
//	//	if e.GetType() == MODIFIED_EVENT_TYPE {
//	//		me := e.(*ModifiedEvent)
//	//		//		fmt.Println("DrawPage.HandleEvent.MODIFIED_EVENT_TYPE", me.ModifiedSupport.IsModified())
//	//		if me.ModifiedSupport.IsModified() {
//	//			this.SetModified(true)
//	//			return true
//	//		}
//	//	}
//	return false
//}

func (this *DrawPage) TrackEvent(e event.IEvent) bool {
	if IsMouseEvent(e.GetType()) {
		me := e.(IMouseEvent)
		x := me.X() - this.X()
		y := me.Y() - this.Y()
		if me.GetType() == MOUSE_WHEEL_EVENT_TYPE {
			nme := *(me.(*WheelEvent))
			nme.SetX(x)
			nme.SetY(y)
			nme.Source = this
			nme.PageEvent = &nme
			return this.handleMouseEvent(&nme)
		} else {
			nme := *(me.(*MouseEvent))
			nme.SetX(x)
			nme.SetY(y)
			nme.Source = this
			nme.PageEvent = &nme
			return this.handleMouseEvent(&nme)
		}
	} else {
		if this.FireEvent(e) {
			return true
		}
		if this.FocusLayer != nil && this.FocusLayer.IsEventEnabled() {
			return this.FocusLayer.TrackEvent(e)
		}
		return false
	}
	return false
}

func (this *DrawPage) handleMouseEvent(me IMouseEvent) bool {
	//在画布中分发鼠标事件
	this.Self.(IEventSupport).FireEvent(me)

	x, y := me.X(), me.Y()

	//向层转发鼠标事件
	var layer ILayer
	switch this.MouseEventMode {
	case MEventMode_Hovering:
		layer = this.getHoveringLayer(x, y)
		if this.MouseHoveringLayer != nil && this.MouseHoveringLayer != layer {
			this.MouseHoveringLayer.fireMouseLeaveEvent(me) //离开层
			this.MouseHoveringLayer = nil
		}

		if layer != nil {
			if this.MouseHoveringLayer != layer {
				this.MouseHoveringLayer = layer
				this.MouseHoveringLayer.fireMouseEnterEvent(me)
			}
			return this.layerTrackEvent(layer, me)
		}
		break
	case MEventMode_Top:
		//layer := //this.layers.GetTopElement()
		//layers := this.layers.GetElements()
		//if len(this.layers) > 0 {
		if this.layers.Size() > 0 {
			this.layers.Sort()
			layerE := this.layers.At(this.layers.Size() - 1) //(layers[len(layers)-1]).(ILayer)
			if layerE != nil {
				layer := layerE.(ILayer)
				if layer.IsEventEnabled() && layer.IsVisible() && layer.Intersects(x, y) {
					if this.MouseHoveringLayer != layer {
						this.MouseHoveringLayer = layer
						this.MouseHoveringLayer.fireMouseEnterEvent(me)
					}
					return this.layerTrackEvent(layer, me)
				} else {
					layer.fireMouseLeaveEvent(me) //离开层
				}
			}
		}
		break
	case MEventMode_Manual:
		if this.ManualModeLayer != nil &&
			this.ManualModeLayer.IsEventEnabled() &&
			this.ManualModeLayer.IsVisible() &&
			this.ManualModeLayer.Intersects(x, y) {
			if this.MouseHoveringLayer != this.ManualModeLayer {
				this.MouseHoveringLayer = this.ManualModeLayer
				this.MouseHoveringLayer.fireMouseEnterEvent(me)
			}
			return this.layerTrackEvent(this.MouseHoveringLayer, me)
		} else {
			if this.MouseHoveringLayer != nil {
				this.MouseHoveringLayer.fireMouseLeaveEvent(me) //离开层
				this.MouseHoveringLayer = nil
			}
		}
	}
	return false
}

func (this *DrawPage) getHoveringLayer(x, y int) ILayer {
	var layer ILayer
	this.layers.ForEachLast(func(i int, l IElement) bool {
		//log.Println("getHoveringLayer:", i, l.GetId(), l.IsEventEnabled(), l.IsVisible(), l.Intersects(x, y), x, y)
		if l.IsEventEnabled() && l.IsVisible() && l.Intersects(x, y) {
			layer = l.(ILayer)
			return false
		}
		return true
	})
	return layer
	//	layers := this.layers.GetElements()
	//	for i := len(layers) - 1; i >= 0; i-- {
	//		l := layers[i]
	//		if l.IsEventEnabled() && l.IsVisible() && l.Intersects(x, y) {
	//			return l.(ILayer)
	//		}
	//	}
	//	return nil
}

func (this *DrawPage) layerTrackEvent(layer ILayer, me event.IEvent) bool {
	if layer == nil {
		return false
	}
	vrx, vry := layer.GetVisibleRegionCoordinate()
	if me.GetType() == MOUSE_WHEEL_EVENT_TYPE {
		nme := *(me.(*WheelEvent))
		x, y := layer.GetWorldCoordinate() //GetCoordinate()
		nme.SetX(nme.X() - x + vrx)
		nme.SetY(nme.Y() - y + vry)
		nme.Source = layer
		nme.LayerEvent = &nme
		return layer.TrackEvent(&nme)
	} else {
		nme := *(me.(*MouseEvent))
		x, y := layer.GetWorldCoordinate()
		//vr := layer.GetVisibleRegion()
		nme.SetX(nme.X() - x + vrx)
		nme.SetY(nme.Y() - y + vry)
		nme.Source = layer
		nme.LayerEvent = &nme
		return layer.TrackEvent(&nme)
	}
}

/**
 * 画所有图层, 将从最低层图层开始绘制
 */
func (this *DrawPage) Draw() {
	this.drawing = true
	this.SetModified(false)
	this.GraphicsEngine.GetLayerMerger().Merge(this.Self, func() {
		this.drawing = false
		this.Self.FireEvent(this.paintEvent)
	})
}

func (this *DrawPage) IsDrawing() bool {
	return this.drawing
}

func (this *DrawPage) SetRenderEnable(b bool) {
	this.renderEnable = b
}

func (this *DrawPage) setRendering(b bool) {
	if b {
		atomic.StoreInt32(&this.rendering, 1)
	} else {
		atomic.StoreInt32(&this.rendering, 0)
	}
}

func (this *DrawPage) IsRendering() bool {
	return atomic.LoadInt32(&this.rendering) == 1
}

/**
 * 返回最近一次绘制页使用时间
 */
func (this *DrawPage) GetDrawSpendTime() int64 {
	return this.spendTime
}

//拷贝缓冲区内容到画布
//drawing 指示是否调用Draw(ge IGraphicsEngine)方法绘制
func (this *DrawPage) render() {
	//this.setRendering(true)
	timespender := NewTimespender("DrawPage.Render")
	defer func() {
		//this.setRendering(false)
		this.spendTime = timespender.Spendtime()
		this.statRenderCount++
		this.statRenderTime += this.spendTime
		//timespender.Print()
	}()

	this.GraphicsEngine.Clear()
	this.Draw()
	if this.showStatInfoEnable {
		this.Self.ShowFPS(30)
	}
	this.GraphicsEngine.SwapBuffers()
	//this.setRendering(false)
	//log.Println("DrawPage) render:", this.IsRendering())
}

func (this *DrawPage) ShowFPS(statTimes int) {
	fps, _ := this.GetFPS(statTimes)
	//if isstat {
	fontfamily, _ := this.GraphicsEngine.NewFontFamily("Arial")
	brush, _ := this.GraphicsEngine.NewBrush(color.NRGBA{255, 255, 0, 255})
	font, _ := this.GraphicsEngine.NewFont(fontfamily, 15, IFontStyle(gdiplus.FontStyleRegular), IUnit(gdiplus.UnitPixel))
	rect := image.Rect(10, 10, 150, 26)
	//this.GraphicsEngine.DrawRectangle(&rect)
	this.GraphicsEngine.DrawString(fmt.Sprintf("%000.f  %###d / FPS", fps, this.actualFps), font, &rect, nil, brush)

	//brush, _ = this.GraphicsEngine.NewBrush(color.NRGBA{255, 255, 255, 255})
	//rect = image.Rect(11, 11, 110, 26)
	//this.GraphicsEngine.DrawString(fmt.Sprintf("%0.2f / FPS", fps), font, &rect, nil, brush)

	rect = image.Rect(10, 35, 150, 60)
	spendtime := float64(this.GetDrawSpendTime()) / 1000000.0

	this.GraphicsEngine.DrawString(fmt.Sprintf("%vMS / F", fmt.Sprintf("%0.2f0000", spendtime)[0:5]), font, &rect, nil, brush)

	//this.SetModified(true)
}

func (this *DrawPage) OnRender(f func(IDrawPage, bool)) EventSubscription {
	if this.onRenderEvent == nil {
		this.onRenderEvent = CreateEvent(func(IDrawPage, bool) {})
	}
	return this.onRenderEvent.Listen(f)
}

func (this *DrawPage) Render(necessary ...bool) {
	if atomic.LoadInt32(&this.lockRendering) > 0 {
		return
	}
	if this.IsRendering() || !this.renderEnable {
		atomic.AddInt32(&this.lockRendering, 1)
		//if !this.renderEnable {
		//	log.Println("DrawPage.Render.renderEnable:", this.renderEnable)
		//}
		for this.IsRendering() || !this.renderEnable {
			time.Sleep(time.Millisecond)
		}
		this.setRendering(true) //必须放在atomic.AddInt32(&this.lockRendering, -1)之前
		atomic.AddInt32(&this.lockRendering, -1)
	} else {
		this.setRendering(true)
	}

	if this.onRenderEvent != nil {
		this.onRenderEvent.Fire(this.Self, false)
	}
	defer func() {
		this.setRendering(false)
		atomic.AddInt32(&this.statActualFps, 1)
		if this.onRenderEvent != nil {
			this.onRenderEvent.Fire(this.Self, true)
		}
	}()

	if len(necessary) > 0 && necessary[0] {
		this.render()
	} else if this.IsModified() || this.showStatInfoEnable {
		this.render()
	}
}

func (this *DrawPage) SetShowStatInfoEnable(b bool) {
	this.showStatInfoEnable = true
}

//当统计次数大于limit时，返回true
func (this *DrawPage) GetFPS(limit int) (float64, bool) {
	if this.statRenderCount >= limit {
		this.statFps = float64(this.statRenderCount) / (float64(this.statRenderTime) / 1000000000)
		this.statRenderCount = 0
		this.statRenderTime = 0
		//		fmt.Printf("this.statFps:%f\n", this.statFps)
		return this.statFps, true
	}
	return this.statFps, false
}

func (this *DrawPage) TrackMouseEvent(t event.Type, source interface{}, x, y int, buttons MButton, modifier KeyboardModifier, kseq *KeySequence) {
	if source == nil {
		source = this
	}
	me := NewMouseEvent(t, source, x, y, buttons, modifier)
	me.KeySequence = kseq
	this.TrackEvent(me)
}

func (this *DrawPage) GetCursor() ICursor {
	layer := this.GetMouseHoveringLayer()
	if layer != nil {
		e := layer.GetMouseHoveringElement()
		var c ICursor
		if e != nil {
			e1 := e.GetMouseHoveringElement()
			if e1 != nil {
				c = e1.GetCursor()
			}
			if c == nil {
				c = e.GetCursor()
			}
		}
		if c == nil {
			c = layer.GetCursor()
		}

		return c
	}
	return nil
}
