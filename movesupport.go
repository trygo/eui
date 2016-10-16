package eui

import (
	"sync/atomic"
)

type IMoveSupport interface {
	PrepareTransform(x, y int)
	ReferencePointX() int
	ReferencePointY() int
	IsPreparingTransform() bool
	IsMoving() bool //正在移动
	SetMoving(b bool)
	EndTransform()
	MoveTo(x, y int, angle float32)   //移动到目标坐标
	MoveBy(dx, dy int, angle float32) //增量移动到目标位置
}

type MoveSupport struct {
	refPointX, refPointY int32 //移动参照点
	preparingTransform   bool  //准备移动标记
	moving               bool
}

/**
 * 准备变换参照点，x,y 为参照点坐标，一般是在keyPressEvent事件中调用此方法, 参数是当前鼠标x,y坐标
 */
func (this *MoveSupport) PrepareTransform(x, y int) {
	atomic.StoreInt32(&this.refPointX, int32(x))
	atomic.StoreInt32(&this.refPointY, int32(y))
	this.preparingTransform = true
	//	atomic.StoreInt32(&this.preparingTransform, 1)
}

/**
 * 返回参照点X坐标
 */
func (this *MoveSupport) ReferencePointX() int {
	return int(atomic.LoadInt32(&this.refPointX))
}

/**
 * 返回参照点Y坐标
 */
func (this *MoveSupport) ReferencePointY() int {
	return int(atomic.LoadInt32(&this.refPointY))
}

/**
 * 检查是否已经准备好了参照点
 */
func (this *MoveSupport) IsPreparingTransform() bool {
	//return atomic.LoadInt32(&this.preparingTransform) == 1
	return this.preparingTransform
}

/**
 * 结束变换
 */
func (this *MoveSupport) EndTransform() {
	this.preparingTransform = false
	//atomic.StoreInt32(&this.preparingTransform, 0)
}

func (this *MoveSupport) IsMoving() bool {
	//	return atomic.LoadInt32(&this.moving) == 1
	return this.moving
}

func (this *MoveSupport) SetMoving(b bool) {
	//	if b {
	//		atomic.StoreInt32(&this.moving, 1)
	//	} else {
	//		atomic.StoreInt32(&this.moving, 0)
	//	}
	this.moving = b
}
