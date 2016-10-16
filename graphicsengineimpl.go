package eui

import (
	"errors"
	//	"fmt"
	. "github.com/trygo/winapi"
	"image"
	"image/color"
	"image/draw"
)

/*
  IGraphicsEngine的默认实现
*/

//标记是否绘制区域边线，调试用
var drawRegionStrokeLine bool
var regionStrokeLineColor color.Color = color.NRGBA{255, 255, 255, 255}

func DrawRegionStrokeLine(b bool, c ...color.Color) {
	drawRegionStrokeLine = b
	if len(c) > 0 {
		regionStrokeLineColor = c[0]
	}
}
func DrawRegionStrokeLineEnable() (bool, color.Color) {
	return drawRegionStrokeLine, regionStrokeLineColor
}

type pathtype int

const (
	ellipse pathtype = iota
	rectangle
)

type pathitem struct {
	typ pathtype
	val interface{}
}

type Path struct {
	items []*pathitem
}

func NewDefaultPath() *Path {
	return &Path{items: make([]*pathitem, 0)}
}

func (this *Path) AddLine(x1, y1, x2, y2 float32) {
	panic(errors.New("Not implemented"))
}

func (this *Path) AddBezier(x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	panic(errors.New("Not implemented"))
}

func (this *Path) AddArc(r *image.Rectangle, startAngle, sweepAngle float32) {
	panic(errors.New("Not implemented"))
}

func (this *Path) AddEllipse(r *image.Rectangle) {
	this.items = append(this.items, &pathitem{typ: ellipse, val: r})
}

func (this *Path) AddRectangle(r *image.Rectangle) {
	this.items = append(this.items, &pathitem{typ: rectangle, val: r})
}

func (this *Path) AddPolygon(points []image.Point) {
	panic(errors.New("Not implemented"))
}

func (this *Path) AddString(layoutRect *image.Rectangle, text string, fontSize float32, family IFontFamily, style IFontStyle, stringFormat IStringFormat) {
	panic(errors.New("Not implemented"))
}

func (this *Path) LastPoint() (x, y float32) {
	panic(errors.New("Not implemented"))
}

func (this *Path) LastError() error {
	return nil
}

func (this *Path) SetVisibleRegion(vr IVisibleRegion) {

}

func (this *Path) intersects(r *image.Rectangle, x, y int) bool {
	return r.Min.X <= x && x < r.Max.X &&
		r.Min.Y <= y && y < r.Max.Y
}

func (this *Path) IsVisible(x, y float32, ge IGraphicsEngine) bool {
	for _, item := range this.items {
		switch item.typ {
		case ellipse:
			return this.intersects(item.val.(*image.Rectangle), int(x), int(y))
		case rectangle:
			return this.intersects(item.val.(*image.Rectangle), int(x), int(y))
		}
	}
	return false
}

func (this *Path) IsOutlineVisible(x, y float32, ge IGraphicsEngine) bool {
	//	panic(errors.New("Not implemented"))
	return false
}

func (this *Path) Reset() {
	this.items = this.items[0:0]
}

func (this *Path) Release() {
	this.items = this.items[0:0]
}

type Region struct {
	rects []*image.Rectangle
}

func newRegion() *Region {
	return &Region{rects: make([]*image.Rectangle, 0)}
}

func (this *Region) GetRects() []*image.Rectangle {
	return this.rects
}

//如果与已有区域相交，就合并之
func (this *Region) UnionWithIntersect(rect image.Rectangle) {
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

func (this *Region) Clear() {
	this.rects = this.rects[0:0]
}

func (this *Region) Release() {

}

//默认绘图引擎，此实现只将图形绘制到image.Image上
type DefaultGraphicsEngine struct {
	//Page         IDrawPage
	Canvas       draw.Image
	bakCanvas    draw.Image
	ClipedRegion IRegion
	LayerMerger  IMerger
	//LayerCanvasCreater         func(l ILayer) draw.Image
	LayerGraphicsEngineCreater GraphicsEngineCreaterType //func(l ILayer) (IGraphicsEngine, draw.Image)
	Op                         draw.Op

	Paths []IPath

	LineWidth  float32
	StrokColor color.Color
	FillColor  color.Color
}

//canvas页输出设备
//layerCanvasCreater函数中创建层输出设备
//func NewDefaultGraphicsEngine(page IDrawPage, canvas draw.Image, layerCanvasCreater func(l ILayer) draw.Image) *DefaultGraphicsEngine {
//canvas页输出设备
//layerGraphicsEngineCreater函数用于创建层IGraphicsEngine
//func NewDefaultGraphicsEngine(page IDrawPage, icanvas draw.Image, layerGraphicsEngineCreater GraphicsEngineCreaterType) *DefaultGraphicsEngine {
func NewDefaultGraphicsEngine(icanvas draw.Image, layerGraphicsEngineCreater GraphicsEngineCreaterType) *DefaultGraphicsEngine {
	ge := &DefaultGraphicsEngine{Canvas: icanvas, LayerGraphicsEngineCreater: layerGraphicsEngineCreater, Op: draw.Src}
	ge.Paths = make([]IPath, 0)
	ge.LineWidth = 1.0
	ge.StrokColor = color.NRGBA{0, 0, 0, 255}
	ge.FillColor = color.NRGBA{255, 255, 255, 255}
	return ge
}

func (this *DefaultGraphicsEngine) SetVisibleRegion(vr IVisibleRegion) {

}

func (this *DefaultGraphicsEngine) Release() {

}

func (this *DefaultGraphicsEngine) LastError() error {
	return nil
}

func (this *DefaultGraphicsEngine) NewRegion() IRegion {
	return newRegion()
}

func (this *DefaultGraphicsEngine) NewFontFamily(name string) (IFontFamily, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) NewStringFormat() (IStringFormat, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) NewFont(family IFontFamily, emSize float32, style IFontStyle, unit IUnit) (IFont, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) NewBrush(c color.Color) (IBrush, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) NewPen(c color.Color, width ...float32) (IPen, error) {
	panic(errors.New("Not implemented"))
}
func (this *DefaultGraphicsEngine) NewPen2(brush IBrush, width ...float32) (IPen, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) Clear() {
	var p1, p2 image.Point
	rects := this.getDrawRectangles()
	for _, r := range rects {
		p1, p2 = r.Min, r.Max
		for x := p1.X; x < p2.X; x++ {
			for y := p1.Y; y < p2.Y; y++ {
				this.Canvas.Set(x, y, this.FillColor)
			}
		}
	}

	//p1, p2 = this.Canvas.Bounds().Min, this.Canvas.Bounds().Max
	//for x := p1.X; x < p2.X; x++ {
	//	for y := p1.Y; y < p2.Y; y++ {
	//		this.Canvas.Set(x, y, this.FillColor)
	//	}
	//}
}

func (this *DefaultGraphicsEngine) GetLayerMerger() IMerger {
	//	fmt.Println("GetLayerMerger")
	if this.LayerMerger == nil {
		//this.LayerMerger = NewDefaultLayerMerger(this.Page, func(l ILayer) (IGraphicsEngine, IImage) {
		this.LayerMerger = NewDefaultLayerMerger(func(l ILayer) (IGraphicsEngine, IImage) {
			return this.LayerGraphicsEngineCreater(l)
		})
	}
	return this.LayerMerger
}

func (this *DefaultGraphicsEngine) SwapBuffers() {
}

func (this *DefaultGraphicsEngine) SwapThreeBuffers() {
}

func (this *DefaultGraphicsEngine) GetBuffer() IImage {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) CopyBuffer(x, y INT, srcx, srcy, srcw, srch INT, clear ...bool) {

}

func (this *DefaultGraphicsEngine) AddPaths(path ...IPath) {
	this.Paths = append(this.Paths, path...)
}

func (this *DefaultGraphicsEngine) SetRenderPen(pen IPen) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) SetRenderBrush(brush IBrush) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) Render(mode ...RenderMode) {
	//println("DefaultGraphicsEngine.Render", len(this.Paths))
	if len(this.Paths) == 0 {
		return
	}
	for _, p := range this.Paths {
		for _, item := range (p.(*Path)).items {
			switch item.typ {
			case ellipse:
				this.DrawRectangle(item.val.(*image.Rectangle)) //还没实现画圆，先用矩形代替
			case rectangle:
				this.DrawRectangle(item.val.(*image.Rectangle))
			}
		}
	}
	this.Paths = this.Paths[0:0]
}

func (this *DefaultGraphicsEngine) SetClipRect(rect *image.Rectangle, adjustPos ...bool) {
	rg := this.NewRegion()
	rg.UnionWithIntersect(*rect)
	this.SetClip(rg)
}

func (this *DefaultGraphicsEngine) SetClip(cliped IRegion) {
	this.ClipedRegion = cliped
	//println("DefaultGraphicsEngine.SetClip.clipedRegion", len(this.ClipedRegion.GetRects()))
	this.bakCanvas = this.Canvas
	this.Canvas = image.NewNRGBA(this.Canvas.Bounds())
}

func (this *DefaultGraphicsEngine) ResetClip() {
	clipCanvas := this.Canvas
	this.Canvas = this.bakCanvas
	rects := this.getDrawRectangles()
	for _, r := range rects {
		draw.Draw(this.Canvas, *r, clipCanvas, r.Min, this.Op)
	}

	if drawRegionStrokeLine {
		strokColor := this.StrokColor
		op := this.Op
		this.Op = draw.Over
		this.SetStrokeColor(regionStrokeLineColor)
		for _, r := range rects {
			this.strokeRectangle(image.ZP, r)
		}
		this.Op = op
		this.StrokColor = strokColor
	}

	this.bakCanvas = nil
	this.ClipedRegion = nil
}

func (this *DefaultGraphicsEngine) SetStrokeColor(c color.Color) {
	this.StrokColor = c
}
func (this *DefaultGraphicsEngine) SetFillColor(c color.Color) {
	this.FillColor = c
}
func (this *DefaultGraphicsEngine) SetLineWidth(lineWidth float32) {
	this.LineWidth = lineWidth
}

//在p坐标位置绘制src(sp)
//func (this *DefaultGraphicsEngine) DrawImage(p image.Point, src_ IImage, sp image.Point, ops ...draw.Op) {
func (this *DefaultGraphicsEngine) DrawImage(src_ IImage, x, y, srcx, srcy, srcwidth, srcheight int, ops ...draw.Op) {
	var src image.Image
	switch src_.(type) {
	case image.Image:
		src = src_.(image.Image) //gdiplus.IImage
	case IBitmap:
		src = src_.(*DefaultBitmap).img
	default:
		return
	}

	op := this.Op
	if len(ops) > 0 {
		op = ops[0]
	}

	p := image.Point{x, y}
	sp := image.Point{srcx, srcy}

	srcRect := src.Bounds()
	srcMin := srcRect.Min.Add(sp.Add(p))
	srcRect2 := image.Rectangle{Min: srcMin, Max: image.Pt(srcMin.X+srcRect.Dx(), srcMin.Y+srcRect.Dy())}
	for _, r := range this.getDrawRectangles() {
		isIntersect := IsIntersect(r, &srcRect2)
		//println("DefaultGraphicsEngine.DrawImage", len(this.getDrawRectangles()), r.String(), "p:", p.String(), "sp", sp.String(), "srcRect:", srcRect.String(), "srcRect2:", srcRect2.String(), r.Intersect(srcRect2).String(), isIntersect)
		if isIntersect {
			region := image.Rectangle{Min: p, Max: r.Max}
			//println("DefaultGraphicsEngine.DrawImage.region", region.String(), "src.Bounds():", src.Bounds().String(), "sp:", sp.String())
			draw.Draw(this.Canvas, region, src, sp, op)
		}
	}
}

func (this *DefaultGraphicsEngine) getDrawRectangles() []*image.Rectangle {
	if this.ClipedRegion == nil {
		r := this.Canvas.Bounds()
		return []*image.Rectangle{&r}
	} else {
		return this.ClipedRegion.GetRects()
	}
}

func (this *DefaultGraphicsEngine) fillColors(img draw.Image, c color.Color, left, top, right, bottom int) {
	for x := left; x < right; x++ {
		for y := top; y < bottom; y++ {
			img.Set(x, y, c)
		}
	}
}

func (this *DefaultGraphicsEngine) strokeRectangle(p image.Point, rect *image.Rectangle) {
	srcRect := image.Rectangle{image.ZP, image.Pt(rect.Dx(), rect.Dy())}
	src := image.NewNRGBA(srcRect)
	lw := int(this.LineWidth)

	//画线, 上边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Min.Y, srcRect.Max.X, srcRect.Min.Y+lw)
	//画下边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Max.Y-lw, srcRect.Max.X, srcRect.Max.Y)
	//画左边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Min.Y, srcRect.Min.X+lw, srcRect.Max.Y)
	//画右边
	this.fillColors(src, this.StrokColor, srcRect.Max.X-lw, srcRect.Min.Y, srcRect.Max.X, srcRect.Max.Y)

	//this.DrawImage(rect.Min, src, image.ZP)
	this.DrawImage(src, rect.Min.X, rect.Min.Y, 0, 0, srcRect.Dx(), srcRect.Dy())
}

func (this *DefaultGraphicsEngine) DrawRectangle(rect *image.Rectangle) {
	srcRect := image.Rectangle{image.ZP, image.Pt(rect.Dx(), rect.Dy())}
	src := image.NewNRGBA(srcRect)
	lw := int(this.LineWidth)

	//填充
	left := srcRect.Min.X + lw
	top := srcRect.Min.Y + lw
	right := srcRect.Max.X - lw
	bottom := srcRect.Max.Y - lw
	this.fillColors(src, this.FillColor, left, top, right, bottom)

	//画线, 上边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Min.Y, srcRect.Max.X, srcRect.Min.Y+lw)
	//画下边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Max.Y-lw, srcRect.Max.X, srcRect.Max.Y)
	//画左边
	this.fillColors(src, this.StrokColor, srcRect.Min.X, srcRect.Min.Y, srcRect.Min.X+lw, srcRect.Max.Y)
	//画右边
	this.fillColors(src, this.StrokColor, srcRect.Max.X-lw, srcRect.Min.Y, srcRect.Max.X, srcRect.Max.Y)
	//println("DefaultGraphicsEngine.DrawRectangle, rect:", rect.String(), srcRect.String())
	//SaveToPngFile(src, fmt.Sprintf("%v-%v-%v.png", src.Bounds().String(), src.Bounds().Dx(), src.Bounds().Dy()))
	//this.DrawImage(rect.Min, src, image.ZP)
	this.DrawImage(src, rect.Min.X, rect.Min.Y, 0, 0, srcRect.Dx(), srcRect.Dy())
}

//func (this *DefaultGraphicsEngine) DrawString(p *image.Rectangle, text string, emSize float32) {
//	panic(errors.New("Not implemented"))
//}

func (this *DefaultGraphicsEngine) MeasureString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat) (boundingBox *image.Rectangle, codepointsFitted, linesFilled int) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) DrawString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat, brush IBrush) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) DrawLine(x1, y1, x2, y2 float32) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) CacheImage(bitmap IBitmap) (ICachedBitmap, error) {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) CloneImage(bitmap IBitmap, x, y, width, height int) IBitmap {
	panic(errors.New("Not implemented"))
}

func (this *DefaultGraphicsEngine) LoadImage(filename string) (IBitmap, error) {
	img, err := LoadImage(filename)
	if err != nil {
		return nil, err
	}
	rect := img.Bounds()
	return &DefaultBitmap{img: img, W: uint(rect.Dx()), H: uint(rect.Dy())}, err
}

type DefaultBitmap struct {
	img  image.Image
	W, H uint
}

func (this *DefaultBitmap) Width() uint {
	return this.W
}

func (this *DefaultBitmap) Height() uint {
	return this.H
}

func (this *DefaultBitmap) Release() {

}

func (this *DefaultBitmap) GetFrameCount() uint {
	return 1
}

func (this *DefaultBitmap) GetFrameDelays() []int {
	return nil
}

func (this *DefaultBitmap) SelectActiveFrame(frameIndex uint) {

}

func (this *DefaultBitmap) SelectNextFrame() {

}

func (this *DefaultBitmap) GetFrameDelay() int {
	return 0
}
func (this *DefaultBitmap) GetFrameIndex() uint {
	return 0
}

func (this *DefaultBitmap) RotateFlipX() {

}

func (this *DefaultBitmap) RotateFlipY() {

}
