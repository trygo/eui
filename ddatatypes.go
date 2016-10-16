package eui

//type Real int
type REAL float32

type AlignmentType int8

const (
	AlignmentNear   AlignmentType = 0
	AlignmentCenter AlignmentType = 1
	AlignmentFar    AlignmentType = 2
	AlignmentClean  AlignmentType = -1
)

/**
 * 图层或元素的鼠标事件响应模式
 */
type MouseEventMode int8

const (
	MEventMode_Hovering MouseEventMode = iota //鼠标正悬停的层或元素
	//MEventMode_All, //所有图层或元素
	MEventMode_Top    //仅顶层或元素
	MEventMode_Manual //指定层或元素
)

/**
 * 层中元素绘制模式
 */
type DrawMode int8

const (
	DrawMode_Auto   DrawMode = iota //自动选择
	DrawMode_Region                 //区域绘制，仅绘制被修改过的可视元素, (默认)
	DrawMode_All                    //绘制所有可视元素
)

/**
 * 层卷动模式,当层高或宽度大于视口高或宽度时，将采用指定模式卷动
 */
type ScrollMode int8

const (
	ScrollMode_Auto  ScrollMode = 0 //自动选择
	ScrollMode_CAMAC ScrollMode = 1 //卡马克
)

/**
 * 鼠标按钮状态定义
 */
type MButton uint16

const (
	MButton_No    MButton = 0x0000
	MButton_Left  MButton = 0x0001
	MButton_Right MButton = 0x0002
	MButton_Mid   MButton = 0x0010
	MButton_X1    MButton = 0x0020
	MButton_X2    MButton = 0x0040
	MButton_Mask  MButton = 0x00ff
)

/**
 * 方向, Orientation
 */

type Orientation uint8

const (
	OHorizontal            Orientation = 0x1                     //水平，横向
	OVertical              Orientation = 0x2                     //垂直，纵向
	OHorizontalAndVertical Orientation = OHorizontal | OVertical //双向
)

type KeyboardModifier uint8

const (
	ModNone    KeyboardModifier = 0
	ModShift   KeyboardModifier = 1
	ModControl KeyboardModifier = 2
	ModAlt     KeyboardModifier = 4
	ModSuper   KeyboardModifier = 8
)

func (m KeyboardModifier) Shift() bool {
	return m&ModShift != 0
}

func (m KeyboardModifier) Control() bool {
	return m&ModControl != 0
}

func (m KeyboardModifier) Alt() bool {
	return m&ModAlt != 0
}

func (m KeyboardModifier) Super() bool {
	return m&ModSuper != 0
}

/**
 * 键定义
 * @SEE gxui.KeyboardKey
 */

//const (
//	Key_Back       = 8
//	Key_Tab        = 9
//	Key_LF         = 10
//	Key_Clear      = 12
//	Key_Return     = 13
//	Key_Shift      = 16
//	Key_Control    = 17
//	Key_Menu       = 18 //Key_Control + Key_Alt
//	Key_Pause      = 19
//	Key_Capital    = 20
//	Key_Kana       = 0x15
//	Key_Hangeul    = 0x15
//	Key_Hangul     = 0x15
//	Key_Junja      = 0x17
//	Key_Final      = 0x18
//	Key_Hanja      = 0x19
//	Key_Kanji      = 0x19
//	Key_Escape     = 0x1b
//	Key_Convert    = 0x1c
//	Key_Nonconvert = 0x1d
//	Key_Accept     = 0x1e
//	Key_Modechange = 0x1f
//	Key_Space      = 32
//	Key_Prior      = 33
//	Key_Next       = 34
//	Key_End        = 35
//	Key_Home       = 36
//	Key_Left       = 37
//	Key_Up         = 38
//	Key_Right      = 39
//	Key_Down       = 40
//	Key_Select     = 41
//	Key_Print      = 42
//	Key_Execute    = 43
//	Key_Snapshot   = 44
//	Key_Insert     = 45
//	Key_Delete     = 46
//	Key_Help       = 47
//	Key_Lwin       = 0x5b
//	Key_Rwin       = 0x5c
//	Key_Apps       = 0x5d
//	Key_Sleep      = 0x5f
//	Key_Numpad0    = 0x60
//	Key_Numpad1    = 0x61
//	Key_Numpad2    = 0x62
//	Key_Numpad3    = 0x63
//	Key_Numpad4    = 0x64
//	Key_Numpad5    = 0x65
//	Key_Numpad6    = 0x66
//	Key_Numpad7    = 0x67
//	Key_Numpad8    = 0x68
//	Key_Numpad9    = 0x69
//	Key_Multiply   = 0x6a
//	Key_Add        = 0x6b
//	Key_Separator  = 0x6c
//	Key_Subtract   = 0x6d
//	Key_Decimal    = 0x6e
//	Key_Divide     = 0x6f
//	Key_F1         = 0x70
//	Key_F2         = 0x71
//	Key_F3         = 0x72
//	Key_F4         = 0x73
//	Key_F5         = 0x74
//	Key_F6         = 0x75
//	Key_F7         = 0x76
//	Key_F8         = 0x77
//	Key_F9         = 0x78
//	Key_F10        = 0x79
//	Key_F11        = 0x7a
//	Key_F12        = 0x7b
//	Key_F13        = 0x7c
//	Key_F14        = 0x7d
//	Key_F15        = 0x7e
//	Key_F16        = 0x7f
//	Key_F17        = 0x80
//	Key_F18        = 0x81
//	Key_F19        = 0x82
//	Key_F20        = 0x83
//	Key_F21        = 0x84
//	Key_F22        = 0x85
//	Key_F23        = 0x86
//	Key_F24        = 0x87
//	Key_Numlock    = 0x90
//	Key_Scroll     = 0x91
//	Key_LShift     = 0xa0
//	Key_RShift     = 0xa1
//	Key_LControl   = 0xa2
//	Key_RControl   = 0xa3
//	Key_LMenu      = 0xa4
//	Key_RMenu      = 0xa5
//)

//layer_type_def":"1:背景层(所有元素都为静态非动画元素), 2:活动层,即操作层, 99:其它",
type LayerType int

const (
	LayerTypeBackground LayerType = 1
	LayerTypeActive     LayerType = 2
	LayerTypeOther      LayerType = 99
)

type ElementType int

const (
	ElementTypeImage     ElementType = 1  //图片
	ElementTypeAnimation ElementType = 2  //动画
	ElementTypeSpirit    ElementType = 3  //精灵
	ElementTypeLine      ElementType = 4  //直线
	ElementTypeBezier    ElementType = 5  //曲线,Bezier
	ElementTypeCurve     ElementType = 6  //曲线,多顶点曲线 Curve
	ElementTypeRectangle ElementType = 7  //矩形
	ElementTypeEllipse   ElementType = 8  //椭圆
	ElementTypePolygon   ElementType = 9  //多边型
	ElementTypeText      ElementType = 10 //文本
	ElementTypeEdit      ElementType = 11 //文本输入
	ElementTypeButton    ElementType = 12 //按钮
	ElementTypeElement   ElementType = 99 //1-n中任何一种
)

const (
	ObstacleNo    = 0 //非障碍物
	ObstacleYes   = 1 //障碍物
	ObstacleUser2 = 2 //自定义2
	ObstacleUser3 = 3 //自定义3
	//自定义N
)

type Alignment struct {
	Horizontal AlignmentType `json:"h" xml:"h"`
	Vertical   AlignmentType `json:"v" xml:"v"`
}

var ZRGBA RGBA

type RGBA struct {
	R uint8 `json:"r" xml:"r"`
	G uint8 `json:"g" xml:"g"`
	B uint8 `json:"b" xml:"b"`
	A uint8 `json:"a" xml:"a"`
}

func (c RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

var ZS Size

type Size struct {
	W int `json:"w" xml:"w"`
	H int `json:"h" xml:"h"`
}

var ZP Point

type Point struct {
	X int `json:"x" xml:"x"`
	Y int `json:"y" xml:"y"`
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

func (this *Point) AddPoint(point *Point) *Point {
	return &Point{X: this.X + point.X, Y: this.Y + point.Y}
}

func (this *Point) SubPoint(point *Point) *Point {
	return &Point{X: this.X - point.X, Y: this.Y - point.Y}
}

func (this *Point) Equals(point *Point) bool {
	return (this.X == point.X) && (this.Y == point.Y)
}

type PointF struct {
	X REAL `json:"x" xml:"x"`
	Y REAL `json:"y" xml:"y"`
}

func NewPointF(x, y REAL) *PointF {
	return &PointF{X: x, Y: y}
}

func (this *PointF) AddPoint(point *PointF) *PointF {
	return &PointF{X: this.X + point.X, Y: this.Y + point.Y}
}

func (this *PointF) SubPoint(point *PointF) *PointF {
	return &PointF{X: this.X - point.X, Y: this.Y - point.Y}
}

func (this *PointF) Equals(point *PointF) bool {
	return (this.X == point.X) && (this.Y == point.Y)
}

var ZR Rect

type Rect struct {
	X int `json:"x" xml:"x"`
	Y int `json:"y" xml:"y"`
	W int `json:"w" xml:"w"`
	H int `json:"h" xml:"h"`
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{X: x, Y: y, W: w, H: h}
}

func (this *Rect) Clone() *Rect {
	return NewRect(this.X, this.Y, this.W, this.H)
}

func (this *Rect) GetLeft() int {
	return this.X
}

func (this *Rect) GetTop() int {
	return this.Y
}

func (this *Rect) GetRight() int {
	return this.X + this.W
}

func (this *Rect) GetBottom() int {
	return this.Y + this.H
}

func (this *Rect) IsEmptyArea() bool {
	return (this.W <= 0) || (this.H <= 0)
}

func (this *Rect) Equals(rect *Rect) bool {
	return this.X == rect.X &&
		this.Y == rect.Y &&
		this.W == rect.W &&
		this.H == rect.H
}

func (this *Rect) Contains(x, y int) bool {
	return x >= this.X && x < this.X+this.W &&
		y >= this.Y && y < this.Y+this.H
}

func (this *Rect) Contains2(pt *Point) bool {
	return this.Contains(pt.X, pt.Y)
}

func (this *Rect) Contains3(rect *Rect) bool {
	return (this.X <= rect.X) && (rect.GetRight() <= this.GetRight()) &&
		(this.Y <= rect.Y) && (rect.GetBottom() <= this.GetBottom())
}

func (this *Rect) Inflate(dx, dy int) {
	this.X -= dx
	this.Y -= dy
	this.W += 2 * dx
	this.H += 2 * dy
}

func (this *Rect) Inflate2(point *Point) {
	this.Inflate(point.X, point.Y)
}

func (this *Rect) Intersect(rect *Rect) bool {
	r, ok := this.intersect(this, rect)
	this.X, this.Y, this.W, this.H = r.X, r.Y, r.W, r.H
	return ok
}

func (this *Rect) intersect(in_a, in_b *Rect) (*Rect, bool) {
	right := Min((in_a.GetRight()), (in_b.GetRight()))
	bottom := Min((in_a.GetBottom()), (in_b.GetBottom()))
	left := Max((in_a.GetLeft()), (in_b.GetLeft()))
	top := Max((in_a.GetTop()), (in_b.GetTop()))
	c := &Rect{X: left, Y: top, W: (right - left), H: (bottom - top)}
	return c, !c.IsEmptyArea()
}

func (this *Rect) IntersectsWith(rect *Rect) bool {
	return (this.GetLeft() < rect.GetRight() &&
		this.GetTop() < rect.GetTop() &&
		this.GetRight() > rect.GetLeft() &&
		this.GetBottom() > rect.GetTop())
}

func UnionI(a, b *Rect) (*Rect, bool) {
	right := Max((a.GetRight()), (b.GetRight()))
	bottom := Max((a.GetBottom()), (b.GetBottom()))
	left := Min((a.GetLeft()), (b.GetLeft()))
	top := Min((a.GetTop()), (b.GetTop()))
	c := &Rect{X: left, Y: top, W: (right - left), H: (bottom - top)}
	return c, !c.IsEmptyArea()
}

func (this *Rect) Offset(dx, dy int) {
	this.X += dx
	this.Y += dy
}

func (this *Rect) Offset2(point *Point) {
	this.X += point.X
	this.Y += point.Y
}

var ZRF RectF

type RectF struct {
	X REAL `json:"x" xml:"x"`
	Y REAL `json:"y" xml:"y"`
	W REAL `json:"w" xml:"w"`
	H REAL `json:"h" xml:"h"`
}

func NewRectF(x, y, w, h REAL) *RectF {
	return &RectF{X: x, Y: y, W: w, H: h}
}

func (this *RectF) Clone() *RectF {
	return NewRectF(this.X, this.Y, this.W, this.H)
}

func (this *RectF) GetLeft() REAL {
	return this.X
}

func (this *RectF) GetTop() REAL {
	return this.Y
}

func (this *RectF) GetRight() REAL {
	return this.X + this.W
}

func (this *RectF) GetBottom() REAL {
	return this.Y + this.H
}

func (this *RectF) IsEmptyArea() bool {
	return (this.W <= 0) || (this.H <= 0)
}

func (this *RectF) Equals(rect *RectF) bool {
	return this.X == rect.X &&
		this.Y == rect.Y &&
		this.W == rect.W &&
		this.H == rect.H
}

func (this *RectF) Contains(x, y REAL) bool {
	return x >= this.X && x < this.X+this.W &&
		y >= this.Y && y < this.Y+this.H
}

func (this *RectF) Contains2(pt *PointF) bool {
	return this.Contains(pt.X, pt.Y)
}

func (this *RectF) Contains3(rect *RectF) bool {
	return (this.X <= rect.X) && (rect.GetRight() <= this.GetRight()) &&
		(this.Y <= rect.Y) && (rect.GetBottom() <= this.GetBottom())
}

func (this *RectF) Inflate(dx, dy REAL) {
	this.X -= dx
	this.Y -= dy
	this.W += 2 * dx
	this.H += 2 * dy
}

func (this *RectF) Inflate2(point *PointF) {
	this.Inflate(point.X, point.Y)
}

func (this *RectF) Intersect(rect *RectF) bool {
	r, ok := this.intersect(this, rect)
	this.X, this.Y, this.W, this.H = r.X, r.Y, r.W, r.H
	return ok
}

func (this *RectF) intersect(in_a, in_b *RectF) (*RectF, bool) {
	right := Minf((in_a.GetRight()), (in_b.GetRight()))
	bottom := Minf((in_a.GetBottom()), (in_b.GetBottom()))
	left := Maxf((in_a.GetLeft()), (in_b.GetLeft()))
	top := Maxf((in_a.GetTop()), (in_b.GetTop()))
	c := &RectF{X: left, Y: top, W: (right - left), H: (bottom - top)}
	return c, !c.IsEmptyArea()
}

func (this *RectF) IntersectsWith(rect *RectF) bool {
	return (this.GetLeft() < rect.GetRight() &&
		this.GetTop() < rect.GetTop() &&
		this.GetRight() > rect.GetLeft() &&
		this.GetBottom() > rect.GetTop())
}

func UnionF(a, b *RectF) (*RectF, bool) {
	right := Maxf((a.GetRight()), (b.GetRight()))
	bottom := Maxf((a.GetBottom()), (b.GetBottom()))
	left := Minf((a.GetLeft()), (b.GetLeft()))
	top := Minf((a.GetTop()), (b.GetTop()))
	c := &RectF{X: left, Y: top, W: (right - left), H: (bottom - top)}
	return c, !c.IsEmptyArea()
}

func (this *RectF) Offset(dx, dy REAL) {
	this.X += dx
	this.Y += dy
}

func (this *RectF) Offset2(point *PointF) {
	this.X += point.X
	this.Y += point.Y
}
