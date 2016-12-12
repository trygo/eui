package graphicsengine

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"sync"
	"unsafe"
)

import (
	. "github.com/tryor/eui"
	"github.com/tryor/gdiplus"
	. "github.com/tryor/winapi"
	. "github.com/tryor/winapi/gdi"
)

/*
  IGraphicsEngine的gdiplus实现
*/

var gpToken ULONG_PTR

func GdiplusStartup() {
	gdiplus.Startup(&gpToken, nil, nil)
}

func GdiplusShutdown() {
	//	gdiplus.Shutdown(gpToken)
	status, err := gdiplus.Shutdown(gpToken)
	fmt.Println("Shutdown.status:", status)
	fmt.Println("Shutdown.err:", err)
}

func GdiplusCreateBuffer(hostHWND HANDLE, width, height int, threeBuffer bool) (hostHDC, bufferHDC, threeBufferHDC HDC) {
	hostHDC = GetDC(HWND(hostHWND))
	// 创建双缓冲
	bufferHDC = CreateCompatibleDC(HWND(hostHDC))
	hbitmap := CreateCompatibleBitmap(hostHDC, uintptr(width), uintptr(height))
	SelectObject(bufferHDC, hbitmap)
	DeleteObject(hbitmap)

	if threeBuffer {
		threeBufferHDC = CreateCompatibleDC(HWND(hostHDC))
		hbitmap := CreateCompatibleBitmap(hostHDC, uintptr(width), uintptr(height))
		SelectObject(threeBufferHDC, hbitmap)
		DeleteObject(hbitmap)
	}

	return
}

type myMutex struct {
	sync.RWMutex
}

func (this *myMutex) Lock() {
	this.RWMutex.Lock()
}

func (this *myMutex) Unlock() {
	this.RWMutex.Unlock()
}

func (this *myMutex) RLock() {
	this.RWMutex.RLock()
}

func (this *myMutex) RUnlock() {
	this.RWMutex.RUnlock()
}

type GdiplusPath struct {
	*gdiplus.GraphicsPath
	ge *GdiplusGraphicsEngine
	//	visibleRegion IVisibleRegion
	locker myMutex
}

func NewGdiplusPath(ge IGraphicsEngine) *GdiplusPath {
	gpath, err := gdiplus.NewGraphicsPath()
	if err != nil {
		panic(err)
	} else {
		p := &GdiplusPath{GraphicsPath: gpath, ge: ge.(*GdiplusGraphicsEngine)}
		//		p.visibleRegion = p
		return p
	}
}

//func (this *GdiplusPath) GetVisibleRegionCoordinate() (x, y int) {
//	return 0, 0
//}

func (this *GdiplusPath) SetVisibleRegion(vr IVisibleRegion) {
	this.locker.Lock()
	defer this.locker.Unlock()
	//	this.visibleRegion = vr
}

//func (this *GdiplusPath) getVRX() int32 {
//	return atomic.LoadInt32(&this.vrX)
//}

//func (this *GdiplusPath) getVRY() int32 {
//	return atomic.LoadInt32(&this.vrY)
//}

func (this *GdiplusPath) Reset() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.GraphicsPath.Reset()
	if this.GraphicsPath.LastResult != gdiplus.Ok {
		log.Println(this.GraphicsPath.LastError)
	}
}

func (this *GdiplusPath) Release() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.GraphicsPath.Release()
	if this.GraphicsPath.LastResult != gdiplus.Ok {
		log.Println(this.GraphicsPath.LastError)
	}
}

func (this *GdiplusPath) IsVisible(x, y float32, ge IGraphicsEngine) bool {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	return bool(this.GraphicsPath.IsVisible(gdiplus.REAL(x-float32(vx)), gdiplus.REAL(y-float32(vy)), ge.(*GdiplusGraphicsEngine).Graphics))
}

func (this *GdiplusPath) IsOutlineVisible(x, y float32, ge IGraphicsEngine) bool {
	this.locker.Lock()
	defer this.locker.Unlock()
	g := ge.(*GdiplusGraphicsEngine)
	var pen *gdiplus.Pen
	if g.renderPen == nil {
		pen, _ = gdiplus.NewPen(ToGdiplusColor(g.StrokColor), gdiplus.REAL(g.LineWidth))
		defer pen.Release()
	} else {
		pen = g.renderPen.(*gdiplus.Pen)
		g.renderPen = nil
	}
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	return bool(this.GraphicsPath.IsOutlineVisible(gdiplus.REAL(x-float32(vx)), gdiplus.REAL(y-float32(vy)), pen, g.Graphics))
}

func (this *GdiplusPath) AddEllipse(r *image.Rectangle) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	this.GraphicsPath.AddEllipse(gdiplus.REAL(r.Min.X-vx), gdiplus.REAL(r.Min.Y-vy), gdiplus.REAL(r.Dx()), gdiplus.REAL(r.Dy()))
}

func (this *GdiplusPath) AddRectangle(r *image.Rectangle) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	r = FormatRect(r)
	x, y, w, h := gdiplus.REAL(r.Min.X-vx), gdiplus.REAL(r.Min.Y-vy), gdiplus.REAL(r.Dx()), gdiplus.REAL(r.Dy())
	this.GraphicsPath.AddRectangle(&gdiplus.RectF{x, y, w, h})
}

func (this *GdiplusPath) AddPolygon(points []image.Point) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()

	ps := make([]gdiplus.Point, len(points))
	for i, p := range points {
		ps[i].X = INT(p.X - vx)
		ps[i].Y = INT(p.Y - vy)
	}
	this.GraphicsPath.AddPolygonI(ps)
}

func (this *GdiplusPath) AddLine(x1, y1, x2, y2 float32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	this.GraphicsPath.AddLine(gdiplus.REAL(x1-float32(vx)), gdiplus.REAL(y1-float32(vy)), gdiplus.REAL(x2-float32(vx)), gdiplus.REAL(y2-float32(vy)))
}

func (this *GdiplusPath) AddBezier(x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	this.GraphicsPath.AddBezier(gdiplus.REAL(x1-float32(vx)), gdiplus.REAL(y1-float32(vy)), gdiplus.REAL(x2-float32(vx)), gdiplus.REAL(y2-float32(vy)),
		gdiplus.REAL(x3-float32(vx)), gdiplus.REAL(y3-float32(vy)), gdiplus.REAL(x4-float32(vx)), gdiplus.REAL(y4-float32(vy)))
}

func (this *GdiplusPath) AddArc(r *image.Rectangle, startAngle, sweepAngle float32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	this.GraphicsPath.AddArc(gdiplus.REAL(r.Min.X-vx), gdiplus.REAL(r.Min.Y-vy), gdiplus.REAL(r.Dx()), gdiplus.REAL(r.Dy()), gdiplus.REAL(startAngle), gdiplus.REAL(sweepAngle))
}

func (this *GdiplusPath) AddString(layoutRect *image.Rectangle, text string, fontSize float32, family IFontFamily, style IFontStyle, stringFormat IStringFormat) {
	this.locker.Lock()
	defer this.locker.Unlock()
	var family_ *gdiplus.FontFamily
	var stringFormat_ *gdiplus.StringFormat
	if family != nil {
		family_ = family.(*gdiplus.FontFamily)
	}
	if stringFormat != nil {
		stringFormat_ = stringFormat.(*GdiplusStringFormat).StringFormat
	}
	vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
	lrect := &gdiplus.Rect{INT(layoutRect.Min.X - vx), INT(layoutRect.Min.Y - vy), INT(layoutRect.Dx()), INT(layoutRect.Dy())}
	this.GraphicsPath.AddStringI(text, family_, gdiplus.FontStyle(style), gdiplus.REAL(fontSize), lrect, stringFormat_)
}

func (this *GdiplusPath) LastPoint() (x, y float32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	lastPoint, status := this.GraphicsPath.GetLastPoint()
	if status == gdiplus.Ok {
		vx, vy := this.ge.visibleRegion.GetVisibleRegionCoordinate() //this.visibleRegion.GetVisibleRegionCoordinate()
		x, y = float32(lastPoint.X)+float32(vx), float32(lastPoint.Y)+float32(vy)
	}
	return
}

func (this *GdiplusPath) LastError() error {
	return this.GraphicsPath.LastError
}

type defaultVisibleRegion struct {
}

func (this *defaultVisibleRegion) GetVisibleRegionCoordinate() (x, y int) {
	return 0, 0
}

type GdiplusGraphicsEngine struct {
	*DefaultGraphicsEngine
	*gdiplus.Graphics
	Graphics2 *gdiplus.Graphics

	hostHWND   HANDLE
	hostHDC    HDC
	canvas     IImage //buffer 2
	canvas2    IImage //buffer 3
	bufferSize Size
	isLayer    bool

	renderPen   IPen
	renderBrush IBrush

	visibleRegion IVisibleRegion
	locker        myMutex
}

//if isLayer, canvas is gdiplus.IImage
//if not isLayer, canvas is winapi.HDC

//func NewGdiplusGraphicsEngine(hostHWND HANDLE, bufferSize Size, icanvas IImage, isLayer bool, layerGraphicsEngineCreater GraphicsEngineCreaterType) *GdiplusGraphicsEngine {
//	ge := &GdiplusGraphicsEngine{hostHWND: hostHWND, bufferSize: bufferSize, canvas: icanvas, isLayer: isLayer}
//	ge.DefaultGraphicsEngine = NewDefaultGraphicsEngine(nil, layerGraphicsEngineCreater)
//	if isLayer {
//		ge.Graphics, _ = gdiplus.NewGraphicsFromImage(icanvas.(gdiplus.IImage))

//	} else {
//		ge.hostHDC, ge.canvas, ge.canvas2 = GdiplusCreateBuffer(hostHWND, bufferSize.W, bufferSize.H, true)
//		ge.Graphics, _ = gdiplus.NewGraphicsFromHDC(ge.canvas.(HDC))
//	}
//	if ge.Graphics.LastResult != gdiplus.Ok {
//		panic(ge.Graphics.LastError)
//	}
//	ge.ClipedRegion = newGdiplusRegion(ge)
//	//	ge.visibleRegion = &defaultVisibleRegion{}
//	return ge
//}

//用于页绘制
func NewGdiplusGraphicsPageEngine(hostHWND HANDLE, bufferSize Size, threeBuffer bool, layerGraphicsEngineCreater GraphicsEngineCreaterType) *GdiplusGraphicsEngine {
	ge := &GdiplusGraphicsEngine{hostHWND: hostHWND, bufferSize: bufferSize}
	ge.DefaultGraphicsEngine = NewDefaultGraphicsEngine(nil, layerGraphicsEngineCreater)

	ge.hostHDC, ge.canvas, ge.canvas2 = GdiplusCreateBuffer(hostHWND, bufferSize.W, bufferSize.H, threeBuffer)
	ge.Graphics, _ = gdiplus.NewGraphicsFromHDC(ge.canvas.(HDC))
	//	ge.Graphics, _ = gdiplus.NewGraphicsFromHDC(GetDC(HWND(hostHWND)))
	if ge.Graphics.LastResult != gdiplus.Ok {
		panic(ge.Graphics.LastError)
	}

	ge.ClipedRegion = newGdiplusRegion(ge)
	return ge
}

//用于层绘制
func NewGdiplusGraphicsLayerEngine(icanvas IImage, icanvas2 IImage) *GdiplusGraphicsEngine {
	ge := &GdiplusGraphicsEngine{canvas: icanvas, canvas2: icanvas2, isLayer: true}
	ge.DefaultGraphicsEngine = NewDefaultGraphicsEngine(nil, nil)
	ge.Graphics, _ = gdiplus.NewGraphicsFromImage(icanvas.(gdiplus.IImage))
	if ge.Graphics.LastResult != gdiplus.Ok {
		panic(ge.Graphics.LastError)
	}
	if icanvas2 != nil {
		ge.Graphics2, _ = gdiplus.NewGraphicsFromImage(icanvas2.(gdiplus.IImage))
		if ge.Graphics2.LastResult != gdiplus.Ok {
			panic(ge.Graphics2.LastError)
		}
	}
	ge.ClipedRegion = newGdiplusRegion(ge)
	return ge
}

func (this *GdiplusGraphicsEngine) SetVisibleRegion(vr IVisibleRegion) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.visibleRegion = vr
}

func (this *GdiplusGraphicsEngine) Release() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.ClipedRegion.Release()
	this.Graphics.Release()
}

func (this *GdiplusGraphicsEngine) LastError() error {
	this.locker.Lock()
	defer this.locker.Unlock()
	return this.Graphics.LastError
}

func (this *GdiplusGraphicsEngine) AddPaths(p ...IPath) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.Paths = append(this.Paths, p...)
}

func (this *GdiplusGraphicsEngine) NewRegion() IRegion {
	this.ClipedRegion.Clear()
	return this.ClipedRegion
}

func (this *GdiplusGraphicsEngine) NewFontFamily(name string) (IFontFamily, error) {
	family, err := gdiplus.NewFontFamily(name, nil)
	if err != nil {
		return nil, err
	}
	return family, nil
}
func (this *GdiplusGraphicsEngine) NewStringFormat() (IStringFormat, error) {
	format, err := gdiplus.NewStringFormat()
	if err != nil {
		return nil, err
	}
	return newGdiplusStringFormat(format), nil
}

func (this *GdiplusGraphicsEngine) NewFont(family IFontFamily, emSize float32,
	style IFontStyle, unit IUnit) (IFont, error) {

	var family_ *gdiplus.FontFamily
	if family != nil {
		family_ = family.(*gdiplus.FontFamily)
	}

	font, err := gdiplus.NewFont(family_, gdiplus.REAL(emSize), gdiplus.FontStyle(style), gdiplus.Unit(unit))
	if err != nil {
		return nil, err
	}

	return font, nil
}

func (this *GdiplusGraphicsEngine) NewBrush(c color.Color) (IBrush, error) {
	brush, err := gdiplus.NewSolidBrush(ToGdiplusColor(c))
	if err != nil {
		return nil, err
	}
	return brush, nil
}

func (this *GdiplusGraphicsEngine) NewPen(c color.Color, width ...float32) (IPen, error) {
	var width_ gdiplus.REAL
	if len(width) > 0 {
		width_ = gdiplus.REAL(width[0])
	} else {
		if this.LineWidth > 0 {
			width_ = gdiplus.REAL(this.LineWidth)
		} else {
			width_ = 1.0
		}
	}
	pen, err := gdiplus.NewPen(ToGdiplusColor(c), width_)
	if err != nil {
		return nil, err
	}
	return pen, nil
}
func (this *GdiplusGraphicsEngine) NewPen2(brush IBrush, width ...float32) (IPen, error) {
	var width_ gdiplus.REAL
	if len(width) > 0 {
		width_ = gdiplus.REAL(width[0])
	} else {
		if this.LineWidth > 0 {
			width_ = gdiplus.REAL(this.LineWidth)
		} else {
			width_ = 1.0
		}
	}

	var brush_ gdiplus.IBrush
	if brush != nil {
		brush_ = brush.(gdiplus.IBrush)
	}

	pen, err := gdiplus.NewPen2(brush_, width_)
	if err != nil {
		return nil, err
	}
	return pen, nil
}

func (this *GdiplusGraphicsEngine) SetClipRect(rect *image.Rectangle, adjustPos ...bool) {
	if len(adjustPos) == 0 || adjustPos[0] {
		vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
		this.Graphics.SetClipRectI(&gdiplus.Rect{INT(rect.Min.X - vx), INT(rect.Min.Y - vy), INT(rect.Dx()), INT(rect.Dy())})
	} else {
		this.Graphics.SetClipRectI(&gdiplus.Rect{INT(rect.Min.X), INT(rect.Min.Y), INT(rect.Dx()), INT(rect.Dy())})
	}
}

func (this *GdiplusGraphicsEngine) SetClip(cliped IRegion) {
	vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
	rgn := cliped.(*GdiplusRegion)
	for _, rect := range rgn.rects {
		rgn.Union(&gdiplus.RectF{gdiplus.REAL(rect.Min.X - vx), gdiplus.REAL(rect.Min.Y - vy), gdiplus.REAL(rect.Dx()), gdiplus.REAL(rect.Dy())})
	}
	this.Graphics.SetClipRegion(rgn.Region)

}

func (this *GdiplusGraphicsEngine) ResetClip() {
	this.Graphics.ResetClip()

	rects := this.ClipedRegion.GetRects()
	if rects != nil {
		drawRegionStrokeLine, regionStrokeLineColor := DrawRegionStrokeLineEnable()
		if drawRegionStrokeLine {
			strokColor := this.StrokColor
			op := this.Op

			this.Op = draw.Over
			this.SetStrokeColor(regionStrokeLineColor)
			this.SetLineWidth(1)
			for _, r := range rects {
				this.DrawRectangle(r)
			}

			this.Op = op
			this.StrokColor = strokColor
		}
	}
}

func (this *GdiplusGraphicsEngine) Clear() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.Graphics.Clear(ToGdiplusColor(this.FillColor))
}

func (this *GdiplusGraphicsEngine) DrawImage(src IImage, x, y, srcx, srcy, srcwidth, srcheight int, ops ...draw.Op) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
	//	if vx > 0 || vy > 0 {
	//		println("GdiplusGraphicsEngine.DrawImage ", vx, vy)
	//	}
	switch src_ := src.(type) {
	case IBitmap:
		var sw, sh INT
		if srcwidth == 0 {
			sw = INT(src_.Width())
		} else {
			sw = INT(srcwidth)
		}
		if srcheight == 0 {
			sh = INT(src_.Height())
		} else {
			sh = INT(srcheight)
		}
		//this.Graphics.DrawImageI6(src.(*GdiplusBitmap), INT(x), INT(y), INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel)
		this.Graphics.DrawImageI7(src.(*GdiplusBitmap), &gdiplus.Rect{INT(x - vx), INT(y - vy), sw, sh}, INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel, nil, nil, 0)
	case *gdiplus.Image:
		var sw, sh INT
		if srcwidth == 0 {
			sw = INT(src_.GetWidth())
		} else {
			sw = INT(srcwidth)
		}
		if srcheight == 0 {
			sh = INT(src_.GetHeight())
		} else {
			sh = INT(srcheight)
		}
		//this.Graphics.DrawImageI6(src_, INT(x), INT(y), INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel)
		this.Graphics.DrawImageI7(src_, &gdiplus.Rect{INT(x - vx), INT(y - vy), sw, sh}, INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel, nil, nil, 0)
	case *gdiplus.Bitmap:
		var sw, sh INT
		if srcwidth == 0 {
			sw = INT(src_.GetWidth())
		} else {
			sw = INT(srcwidth)
		}
		if srcheight == 0 {
			sh = INT(src_.GetHeight())
		} else {
			sh = INT(srcheight)
		}
		//this.Graphics.DrawImageI6(src_, INT(x), INT(y), INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel)
		this.Graphics.DrawImageI7(src_, &gdiplus.Rect{INT(x - vx), INT(y - vy), sw, sh}, INT(srcx), INT(srcy), sw, sh, gdiplus.UnitPixel, nil, nil, 0)
	case gdiplus.IImage:
		this.Graphics.DrawImageI(src_, INT(x-vx), INT(y-vy))
	case *gdiplus.CachedBitmap:
		this.Graphics.DrawCachedBitmap(src_, INT(x-vx), INT(y-vy))
	}
}

func (this *GdiplusGraphicsEngine) MeasureString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat) (boundingBox *image.Rectangle, codepointsFitted, linesFilled int) {
	this.locker.Lock()
	defer this.locker.Unlock()
	var font_ *gdiplus.Font
	if font != nil {
		font_ = font.(*gdiplus.Font)
	}

	var stringFormat_ *gdiplus.StringFormat
	if stringFormat != nil {
		stringFormat_ = stringFormat.(*GdiplusStringFormat).StringFormat
	}

	boundingBox_, codepointsFitted_, linesFilled_, _ := this.Graphics.MeasureString(text, font_,
		&gdiplus.RectF{gdiplus.REAL(layoutRect.Min.X), gdiplus.REAL(layoutRect.Min.Y), gdiplus.REAL(layoutRect.Dx()), gdiplus.REAL(layoutRect.Dy())},
		stringFormat_)

	boundingBox = &image.Rectangle{Min: image.Point{int(boundingBox_.X), int(boundingBox_.Y)}, Max: image.Point{int(boundingBox_.X + boundingBox_.Width), int(boundingBox_.Y + boundingBox_.Height)}}
	codepointsFitted = int(codepointsFitted_)
	linesFilled = int(linesFilled_)
	return
}

func (this *GdiplusGraphicsEngine) DrawString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat, brush IBrush) {
	this.locker.Lock()
	defer this.locker.Unlock()
	var font_ *gdiplus.Font
	if font != nil {
		font_ = font.(*gdiplus.Font)
	}

	var stringFormat_ *gdiplus.StringFormat
	if stringFormat != nil {
		stringFormat_ = stringFormat.(*GdiplusStringFormat).StringFormat
	}

	var brush_ gdiplus.IBrush
	if brush != nil {
		brush_ = brush.(gdiplus.IBrush)
	}

	vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
	//	if vx > 0 || vy > 0 {
	//		println("GdiplusGraphicsEngine.DrawString ", vx, vy)
	//	}
	this.Graphics.DrawString(text, font_,
		&gdiplus.RectF{gdiplus.REAL(layoutRect.Min.X - vx), gdiplus.REAL(layoutRect.Min.Y - vy), gdiplus.REAL(layoutRect.Dx()), gdiplus.REAL(layoutRect.Dy())},
		stringFormat_, brush_)

}

func (this *GdiplusGraphicsEngine) DrawRectangle(rect *image.Rectangle) {
	this.locker.Lock()
	defer this.locker.Unlock()
	vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
	x, y, w, h := gdiplus.REAL(rect.Min.X-vx), gdiplus.REAL(rect.Min.Y-vy), gdiplus.REAL(rect.Dx()), gdiplus.REAL(rect.Dy())
	pen, _ := gdiplus.NewPen(ToGdiplusColor(this.StrokColor), gdiplus.REAL(this.LineWidth))
	defer pen.Release()
	if pen.LastResult == gdiplus.Ok {
		this.Graphics.DrawRectangle(pen, x, y, w, h)
	}
}

func (this *GdiplusGraphicsEngine) DrawLine(x1, y1, x2, y2 float32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	pen, _ := gdiplus.NewPen(ToGdiplusColor(this.StrokColor), gdiplus.REAL(this.LineWidth))
	defer pen.Release()
	if pen.LastResult == gdiplus.Ok {
		vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
		this.Graphics.DrawLine(pen, gdiplus.REAL(x1-float32(vx)), gdiplus.REAL(y1-float32(vy)), gdiplus.REAL(x2-float32(vx)), gdiplus.REAL(y2-float32(vy)))
	}
}

func (this *GdiplusGraphicsEngine) SetRenderPen(pen IPen) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.renderPen = pen
}

func (this *GdiplusGraphicsEngine) SetRenderBrush(brush IBrush) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.renderBrush = brush
}

func (this *GdiplusGraphicsEngine) Render(mode ...RenderMode) {
	this.locker.Lock()
	defer this.locker.Unlock()
	if len(this.Paths) == 0 {
		return
	}

	//println(len(this.Paths))

	var m RenderMode
	if len(mode) > 0 {
		m = mode[0]
	} else {
		m = FillStroke
	}

	var pen *gdiplus.Pen
	if (m & OnlyStroke) > 0 {
		if this.renderPen == nil {
			pen, _ = gdiplus.NewPen(ToGdiplusColor(this.StrokColor), gdiplus.REAL(this.LineWidth))
			//			println("NewPen():", pen)
			defer pen.Release()
		} else {
			pen = this.renderPen.(*gdiplus.Pen)
		}
	}

	var brush gdiplus.IBrush
	if (m & OnlyFill) > 0 {
		if this.renderBrush == nil {
			brush, _ = gdiplus.NewSolidBrush(ToGdiplusColor(this.FillColor))
			//			println("NewSolidBrush():", brush)
			defer brush.Release()
		} else {
			brush = this.renderBrush.(gdiplus.IBrush)
		}
	}

	for _, p := range this.Paths {
		if p != nil {
			gpath := p.(*GdiplusPath).GraphicsPath

			switch m {
			case FillStroke:
				this.Graphics.FillPath(brush, gpath)
				this.Graphics.DrawPath(pen, gpath)
			case OnlyStroke:
				this.Graphics.DrawPath(pen, gpath)
			case OnlyFill:
				this.Graphics.FillPath(brush, gpath)
			}
		}
	}

	this.Paths = this.Paths[0:0]
	this.renderBrush = nil
	this.renderPen = nil

}

//将当前画布内容复制到第二画布
func (this *GdiplusGraphicsEngine) CopyBuffer(x, y INT, srcx, srcy, srcw, srch INT, clear ...bool) {

	// gdiplus.UnitPixel, nil, nil, 0
	//this.Graphics2.DrawImageI6(this.canvas.(gdiplus.IImage), x, y, srcx, srcy, srcw, srch, gdiplus.UnitPixel)
	//vx, vy := this.visibleRegion.GetVisibleRegionCoordinate()
	drect := &gdiplus.Rect{x, y, srcw, srch}
	if len(clear) > 0 && clear[0] {
		this.Graphics2.SetClipRectI(drect)
		defer this.Graphics2.ResetClip()
		this.Graphics2.Clear(ToGdiplusColor(this.FillColor))
	}
	this.Graphics2.DrawImageI7(this.canvas.(gdiplus.IImage), drect, srcx, srcy, srcw, srch, gdiplus.UnitPixel, nil, nil, 0)
}

func (this *GdiplusGraphicsEngine) GetBuffer() IImage {
	return this.canvas
}

func (this *GdiplusGraphicsEngine) SwapBuffers() {
	if this.isLayer {
		if this.Graphics2 != nil {
			this.Graphics, this.Graphics2 = this.Graphics2, this.Graphics
			this.canvas, this.canvas2 = this.canvas2, this.canvas
		}
	} else {
		if this.canvas != nil {
			BitBlt(this.hostHDC, 0, 0, this.bufferSize.W, this.bufferSize.H, this.canvas.(HDC), 0, 0, SRCCOPY)
			if this.canvas2 != nil {
				BitBlt(this.canvas2.(HDC), 0, 0, this.bufferSize.W, this.bufferSize.H, this.canvas.(HDC), 0, 0, SRCCOPY)
			}
		}
	}
}

//threeBuffer
func (this *GdiplusGraphicsEngine) SwapThreeBuffers() {
	BitBlt(this.hostHDC, 0, 0, this.bufferSize.W, this.bufferSize.H, this.canvas2.(HDC), 0, 0, SRCCOPY)
}

func ToGdiplusColor(c color.Color) gdiplus.Color {
	switch rgba := c.(type) {
	case color.NRGBA:
		return gdiplus.NewColor3(BYTE(rgba.A), BYTE(rgba.R), BYTE(rgba.G), BYTE(rgba.B))
	case color.RGBA:
		return gdiplus.NewColor3(BYTE(rgba.A), BYTE(rgba.R), BYTE(rgba.G), BYTE(rgba.B))
	}
	r, g, b, a := c.RGBA() //ARGB
	return gdiplus.NewColor3(BYTE(a), BYTE(r), BYTE(g), BYTE(b))
}

func (this *GdiplusGraphicsEngine) CloneImage(bitmap IBitmap, x, y, width, height int) IBitmap {
	this.locker.Lock()
	defer this.locker.Unlock()
	img := bitmap.(*GdiplusBitmap).Bitmap.CloneI(INT(x), INT(y), INT(width), INT(height), gdiplus.PixelFormat32bppARGB)
	if img != nil {
		return newGdiplusBitmap2(img)
	} else {
		return nil
	}
}

func (this *GdiplusGraphicsEngine) CacheImage(bitmap IBitmap) (ICachedBitmap, error) {
	this.locker.Lock()
	defer this.locker.Unlock()
	return gdiplus.NewCachedBitmap(bitmap.(*GdiplusBitmap).Bitmap, this.Graphics)
}

func (this *GdiplusGraphicsEngine) LoadImage(filename string) (IBitmap, error) {
	this.locker.Lock()
	defer this.locker.Unlock()
	return newGdiplusBitmap(filename)
}

type GdiplusBitmap struct {
	*gdiplus.Bitmap
	frameCount  uint
	frameIndex  uint
	frameDelays []int
}

func newGdiplusBitmap(filename string) (*GdiplusBitmap, error) {
	gdiplusbitmap, err := gdiplus.NewBitmap(filename)
	if err != nil {
		return nil, err
	}
	return newGdiplusBitmap2(gdiplusbitmap), nil
}

func newGdiplusBitmap2(gdiplusbitmap *gdiplus.Bitmap) *GdiplusBitmap {

	bitmap := &GdiplusBitmap{Bitmap: gdiplusbitmap, frameIndex: 1}

	fdcount := gdiplusbitmap.GetFrameDimensionsCount()
	if fdcount > 0 {
		dimensionIDs, status := gdiplusbitmap.GetFrameDimensionsList(fdcount)
		if status == gdiplus.Ok {
			if len(dimensionIDs) > 0 {
				bitmap.frameCount = uint(gdiplusbitmap.GetFrameCount(&dimensionIDs[0]))
				size := gdiplusbitmap.GetPropertyItemSize(gdiplus.PropertyTagFrameDelay)
				if size > 0 {
					pitem, status := gdiplusbitmap.GetPropertyItem(gdiplus.PropertyTagFrameDelay, size)
					if status == gdiplus.Ok {
						bitmap.frameDelays = make([]int, bitmap.frameCount)
						delays := (*int32)(unsafe.Pointer(pitem.Value))
						for i := uint(0); i < bitmap.frameCount; i++ {
							bitmap.frameDelays[i] = int(*delays)
							delaysptr := uintptr(unsafe.Pointer(delays)) + uintptr(pitem.Type)
							delays = (*int32)(unsafe.Pointer(delaysptr))
						}
					}
				}
			}
		}
	}

	return bitmap
}

func (this *GdiplusBitmap) GetFrameCount() uint {
	return this.frameCount
}

func (this *GdiplusBitmap) GetFrameDelays() []int {
	return this.frameDelays
}

func (this *GdiplusBitmap) SelectActiveFrame(frameIndex uint) {
	if frameIndex < this.frameCount {
		this.frameIndex = frameIndex
		this.Bitmap.SelectActiveFrame(gdiplus.FrameDimensionTime, UINT(this.frameIndex+1))
	}
}

func (this *GdiplusBitmap) GetFrameIndex() uint {
	return this.frameIndex
}

func (this *GdiplusBitmap) GetFrameDelay() int {
	return this.frameDelays[this.frameIndex]
}

func (this *GdiplusBitmap) SelectNextFrame() {
	if this.frameIndex >= this.frameCount-1 {
		this.frameIndex = 0
	} else {
		this.frameIndex++
	}
	this.Bitmap.SelectActiveFrame(gdiplus.FrameDimensionTime, UINT(this.frameIndex+1))
}

func (this *GdiplusBitmap) Width() uint {
	return uint(this.Bitmap.GetWidth())
}

func (this *GdiplusBitmap) Height() uint {
	return uint(this.Bitmap.GetHeight())
}

func (this *GdiplusBitmap) RotateFlipX() {
	this.Bitmap.RotateFlip(gdiplus.RotateNoneFlipX)
}

func (this *GdiplusBitmap) RotateFlipY() {
	this.Bitmap.RotateFlip(gdiplus.RotateNoneFlipY)
}

func (this *GdiplusBitmap) Release() {
	this.Bitmap.Release()
}

type GdiplusRegion struct {
	*gdiplus.Region
	rects         []*image.Rectangle
	gdiplusEngine *GdiplusGraphicsEngine
}

func newGdiplusRegion(gdiplusEngine *GdiplusGraphicsEngine) *GdiplusRegion {
	rgn := &GdiplusRegion{rects: make([]*image.Rectangle, 0), gdiplusEngine: gdiplusEngine}
	rgn.Region, _ = gdiplus.NewRegionRectI(&gdiplus.Rect{})
	if rgn.Region.LastResult != gdiplus.Ok {
		panic(rgn.Region.LastError)
	}
	return rgn
}

func (this *GdiplusRegion) GetRects() []*image.Rectangle {
	return this.rects
}

//如果与已有区域相交，就合并之
func (this *GdiplusRegion) UnionWithIntersect(rect image.Rectangle) {
	rect = *FormatRect(&rect) //rect_.Min.X, rect_.Min.Y, rect_.Max.X, rect_.Max.Y)
	if len(this.rects) == 0 {
		this.rects = append(this.rects, &rect)
		return
	}
	rects := make([]*image.Rectangle, 0)
	for _, r := range this.rects {
		if IsIntersect(r, &rect) {
			rect = r.Union(rect)
		} else {
			rects = append(rects, r)
		}
	}
	rects = append(rects, &rect)
	this.rects = rects
}

func (this *GdiplusRegion) Clear() {
	this.rects = this.rects[0:0]
	this.Region.MakeEmpty()
}

func (this *GdiplusRegion) Release() {
	if len(this.rects) > 0 {
		this.Clear()
	}
}

type GdiplusStringFormat struct {
	*gdiplus.StringFormat
	formatFlags gdiplus.StringFormatFlags
}

func newGdiplusStringFormat(format *gdiplus.StringFormat) *GdiplusStringFormat {
	sfmt := &GdiplusStringFormat{StringFormat: format}
	sfmt.formatFlags = gdiplus.StringFormatFlagsMeasureTrailingSpaces
	sfmt.StringFormat.SetFormatFlags(INT(sfmt.formatFlags))
	return sfmt
}

func (this *GdiplusStringFormat) Release() {
	this.StringFormat.Release()
}

func (this *GdiplusStringFormat) SetHorizontalAlignment(align AlignmentType) {
	this.StringFormat.SetAlignment(gdiplus.StringAlignment(align))
}

func (this *GdiplusStringFormat) SetVerticalAlignment(align AlignmentType) {
	this.StringFormat.SetLineAlignment(gdiplus.StringAlignment(align))
}

//@see gdiplus.StringFormatFlags
//func (this *GdiplusStringFormat) SetFormatFlags(flags int) {
//	this.StringFormat.SetFormatFlags(INT(flags))
//}

func (this *GdiplusStringFormat) SetMultiline(b bool) {
	if !b {
		this.StringFormat.SetFormatFlags(INT(gdiplus.StringFormatFlagsNoWrap | this.formatFlags))
	} else {
		this.StringFormat.SetFormatFlags(INT(this.formatFlags))
	}
}
