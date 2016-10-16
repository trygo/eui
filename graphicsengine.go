package eui

import (
	. "github.com/trygo/winapi"
	"image"
	"image/color"
	"image/draw"
)

type RenderMode byte //FillStroke
const (
	OnlyFill   RenderMode = 1
	OnlyStroke RenderMode = 2
	FillStroke RenderMode = OnlyFill | OnlyStroke
)

type IImage interface{}

type IBitmap interface {
	Width() uint
	Height() uint
	Release()
	GetFrameCount() uint
	GetFrameDelays() []int
	SelectActiveFrame(frameIndex uint)
	SelectNextFrame()
	GetFrameDelay() int
	GetFrameIndex() uint

	RotateFlipX()
	RotateFlipY()
}

type ICachedBitmap interface {
	Release()
}

type IVisibleRegion interface {
	GetVisibleRegionCoordinate() (x, y int)
	//GetVisibleRegionSnapshotCoordinate() (x, y int)
}

//绘图引擎
type IGraphicsEngine interface {
	Clear() //清除，填充色采用SetFillColor()方法设置的值
	GetLayerMerger() IMerger
	SetVisibleRegion(vr IVisibleRegion)

	SetRenderPen(pen IPen)       //每次使用后将被清空
	SetRenderBrush(brush IBrush) //每次使用后将被清空
	Render(mode ...RenderMode)   //将图形数据输出到目标设备上
	SwapBuffers()
	SwapThreeBuffers()

	GetBuffer() IImage //返回最新缓冲
	CopyBuffer(x, y INT, srcx, srcy, srcw, srch INT, clear ...bool)

	NewRegion() IRegion
	NewFontFamily(name string) (IFontFamily, error)
	NewStringFormat() (IStringFormat, error)
	NewFont(family IFontFamily, emSize float32, style IFontStyle, unit IUnit) (IFont, error)
	NewBrush(c color.Color) (IBrush, error)
	NewPen(c color.Color, width ...float32) (IPen, error)
	NewPen2(brush IBrush, width ...float32) (IPen, error)

	//adjustPos 指示是否格式化坐标为世界坐标， 默认会调整坐标
	SetClipRect(rect *image.Rectangle, adjustPos ...bool)
	SetClip(cliped IRegion) //设置可绘制区域，clipedRegion指定区域范围
	ResetClip()             //恢复SetClip设置的区域
	AddPaths(path ...IPath)

	LoadImage(filename string) (IBitmap, error)
	CacheImage(bitmap IBitmap) (ICachedBitmap, error)
	CloneImage(bitmap IBitmap, x, y, width, height int) IBitmap

	SetStrokeColor(c color.Color)
	SetFillColor(c color.Color)
	SetLineWidth(lineWidth float32)

	//	DrawImage(p image.Point, src IImage, sp image.Point, ops ...draw.Op)
	//x, y, srcx, srcy, srcwidth, srcheight INT
	DrawImage(src IImage, x, y, srcx, srcy, srcwidth, srcheight int, ops ...draw.Op) //将src(srcx, srcy, srcwidth, srcheight)绘制到(x, y)坐标位置

	DrawRectangle(rect *image.Rectangle)
	//DrawString(p *image.Rectangle, text string, emSize float32)
	DrawString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat, brush IBrush)
	MeasureString(text string, font IFont, layoutRect *image.Rectangle, stringFormat IStringFormat) (boundingBox *image.Rectangle, codepointsFitted, linesFilled int)
	DrawLine(x1, y1, x2, y2 float32)

	LastError() error
	Release()
}

type IFontFamily interface {
	Release()
}
type IFontStyle int
type IStringFormat interface {
	Release()
	SetHorizontalAlignment(align AlignmentType)
	SetVerticalAlignment(align AlignmentType)
	SetMultiline(b bool)
}

type IFont interface {
	Release()
}
type IUnit int

type IBrush interface {
	Release()
}

type IPen interface {
	Release()
}

type IPath interface {
	LastPoint() (x, y float32)
	AddLine(x1, y1, x2, y2 float32)
	AddBezier(x1, y1, x2, y2, x3, y3, x4, y4 float32)
	AddArc(r *image.Rectangle, startAngle, sweepAngle float32)
	AddEllipse(r *image.Rectangle)
	AddRectangle(r *image.Rectangle)
	AddPolygon(points []image.Point)
	AddString(layoutRect *image.Rectangle, text string, fontSize float32, family IFontFamily, style IFontStyle, stringFormat IStringFormat)

	IsVisible(x, y float32, ge IGraphicsEngine) bool
	IsOutlineVisible(x, y float32, ge IGraphicsEngine) bool
	Reset()
	Release()
	LastError() error
	SetVisibleRegion(vr IVisibleRegion)
}

type IRegion interface {
	GetRects() []*image.Rectangle
	UnionWithIntersect(rect image.Rectangle)
	Clear()
	Release()
}
