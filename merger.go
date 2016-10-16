package eui

import (
	//	"fmt"
	"image/draw"
	//	"log"
)

//层合并器
type IMerger interface {
	InitLayerGraphicsEngine(l ILayer)
	Merge(page IDrawPage, drawFinish func()) bool
}

type layerCanvas struct {
	Donec          chan bool
	Layer          ILayer
	Canvas         IImage
	GraphicsEngine IGraphicsEngine
	X, Y           int
}

type GraphicsEngineCreaterType func(l ILayer) (IGraphicsEngine, IImage)

type DefaultLayerMerger struct {
	//Page                  IDrawPage
	LayerCanvasMap        map[ILayer]*layerCanvas
	graphicsEngineCreater GraphicsEngineCreaterType
}

//graphicsEngineCreater为IGraphicsEngine创建函数, 用于为层创建IGraphicsEngine
//canvasCreater目标设备的创建函数, 将在此设备上绘制图形
func NewDefaultLayerMerger(graphicsEngineCreater GraphicsEngineCreaterType) *DefaultLayerMerger {
	return &DefaultLayerMerger{
		graphicsEngineCreater: graphicsEngineCreater,
		LayerCanvasMap:        make(map[ILayer]*layerCanvas, 0)}
}

func (this *DefaultLayerMerger) InitLayerGraphicsEngine(l ILayer) {
	this.getLayerCanvas(l)
}

func (this *DefaultLayerMerger) DeleteLayerCanvas(l ILayer) {
	lc := this.LayerCanvasMap[l]
	if lc != nil {
		delete(this.LayerCanvasMap, l)
		lc.GraphicsEngine.Release()
		close(lc.Donec)
	}
}

func (this *DefaultLayerMerger) getLayerCanvas(l ILayer) *layerCanvas {
	//	pageGraphics := this.Page.GetGraphicsEngine()
	lc := this.LayerCanvasMap[l]
	if lc == nil {
		lc = &layerCanvas{Donec: make(chan bool), Layer: l}
		//lc.Canvas = this.canvasCreater(l)
		if this.graphicsEngineCreater != nil {
			lc.GraphicsEngine, lc.Canvas = this.graphicsEngineCreater(l)
			l.SetGraphicsEngine(lc.GraphicsEngine)

			//			ge, canvas := this.graphicsEngineCreater(l)
			//			if ge != nil {
			//				lc.GraphicsEngine, lc.Canvas = ge, canvas
			//				l.SetGraphicsEngine(lc.GraphicsEngine, true)
			//			} else {
			//				l.SetGraphicsEngine(pageGraphics, false)
			//			}
		}
		if l.GetBackground() != nil {
			if lc.GraphicsEngine != nil {
				lc.GraphicsEngine.SetFillColor(l.GetBackground())
				lc.GraphicsEngine.Clear()
			}
		}

		//lc.GraphicsEngine.SetLineWidth(1)
		this.LayerCanvasMap[l] = lc
	}

	return lc
}

//func (this *DefaultLayerMerger) MergeOld() bool {
//	pageGraphics := this.Page.GetGraphicsEngine()
//	lcs := make([]*layerCanvas, 0)
//	for _, layer := range this.Page.GetLayers() {
//		if layer.IsVisible() {
//			lc := this.getLayerCanvas(layer)
//			lcs = append(lcs, lc)
//			go func(lc *layerCanvas) {
//				lc.Layer.Draw(lc.GraphicsEngine)
//				lc.Donec <- true
//			}(lc)

//		}
//	}
//	for _, lc := range lcs {
//		<-lc.Donec
//		x, y := lc.Layer.GetCoordinate()
//		pageGraphics.DrawImage(lc.Canvas, x, y, 0, 0, 0, 0, draw.Over)
//	}

//	return true
//}

//func (this *DefaultLayerMerger) Merge() bool {
//	pageGraphics := this.Page.GetGraphicsEngine()
//	this.Page.SortLayers()
//	for _, layer := range this.Page.GetLayers() {
//		if layer.IsVisible() {
//			layerCanvas := this.getLayerCanvas(layer)
//			if layerCanvas.GraphicsEngine != nil {
//				layer.Draw(layerCanvas.GraphicsEngine)
//				x, y := layer.GetWorldCoordinate()
//				pageGraphics.DrawImage(layerCanvas.Canvas, x, y, 0, 0, 0, 0, draw.Over)
//			} else {
//				//没有为层准备绘图引擎，就直接使用Page.GetGraphicsEngine()引擎
//				layer.Draw(pageGraphics)
//			}
//		}
//		//println(fmt.Sprintf("DefaultLayerMerger.Merge(), Layer:%p, AllElements:%v, DrawElementCount:%v, DrawSpendTime(MS):%v", layer, layer.GetElementsCount(), layer.GetDrawElementsCount(), float64(layer.GetDrawSpendTime())/1000000))
//	}
//	return true
//}

func (this *DefaultLayerMerger) Merge_old(page IDrawPage, drawFinish func()) bool {
	pageGraphics := page.GetGraphicsEngine()
	page.SortLayers()
	//page.MakeVisibleRegionSnapshot()
	lcs := make([]*layerCanvas, 0)
	for _, layer := range page.GetLayers() {
		if layer.IsVisible() {
			lc := this.getLayerCanvas(layer)
			lc.X, lc.Y = lc.Layer.GetWorldCoordinate()
			lcs = append(lcs, lc)
			//if layerCanvas.GraphicsEngine != nil {
			layer.Draw(lc.GraphicsEngine)
			//log.Println("Layer:", layer.GetId(), float64(layer.GetDrawSpendTime())/1000000.0)

			//			x, y := layer.GetWorldCoordinate()
			//pageGraphics.DrawImage(layerCanvas.Canvas, x, y, 0, 0, 0, 0, draw.Over)
			//			pageGraphics.DrawImage(lc.GraphicsEngine.GetBuffer(), x, y, 0, 0, 0, 0, draw.Over)
			//			} else {
			//				//没有为层准备绘图引擎，就直接使用Page.GetGraphicsEngine()引擎
			//				layer.Draw(pageGraphics)
			//			}
		}
		//println(fmt.Sprintf("DefaultLayerMerger.Merge(), Layer:%p, AllElements:%v, DrawElementCount:%v, DrawSpendTime(MS):%v", layer, layer.GetElementsCount(), layer.GetDrawElementsCount(), float64(layer.GetDrawSpendTime())/1000000))
	}
	drawFinish()
	for _, lc := range lcs {
		//		x, y := lc.Layer.GetWorldCoordinate()
		//		if lc.Layer.GetId() == "activeLayer" || lc.Layer.GetId() == "earth" {
		//			fmt.Println(lc.Layer.GetId(), x, y)
		//		}
		pageGraphics.DrawImage(lc.GraphicsEngine.GetBuffer(), lc.X, lc.Y, 0, 0, 0, 0, draw.Over)
	}

	//log.Println("Page:", float64(page.GetDrawSpendTime())/1000000.0)
	return true
}

func (this *DefaultLayerMerger) Merge(page IDrawPage, drawFinish func()) bool {
	//timespender := NewTimespender("DefaultLayerMerger.Merge")
	pageGraphics := page.GetGraphicsEngine()
	page.SortLayers()
	layers := page.GetLayers()
	lcs := make([]*layerCanvas, 0, len(layers))
	//timespender.Print()

	for _, layer := range layers {
		if layer.IsVisible() {
			lc := this.getLayerCanvas(layer)
			lc.X, lc.Y = lc.Layer.GetWorldCoordinate()
			lcs = append(lcs, lc)
			go func(lc *layerCanvas) {
				lc.Layer.Draw(lc.GraphicsEngine)
				lc.Donec <- true
			}(lc)
		}
	}

	for _, lc := range lcs {
		<-lc.Donec
	}
	drawFinish()

	for _, lc := range lcs {
		pageGraphics.DrawImage(lc.GraphicsEngine.GetBuffer(), lc.X, lc.Y, 0, 0, 0, 0, draw.Over)
	}
	//timespender.Print()
	return true
}
