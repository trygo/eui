package eui

import (
	"github.com/tryor/commons/event"
)

/**
 * 事件类型, 100000以内的事件类型被系统保留
 */

/* 基本事件类型定义，类型编号从1-1000 */

/** XXXXXXX 中事件类型定义，类型编号从1000-1199 */

/** Easy Draw 中事件类型定义，类型编号从1200-1299 */
/*鼠标事件 1200-1219 */
//鼠标事件类型最小编号
const MOUSE_EVENT_MIN_TYPE event.Type = 1200

//鼠标事件类型最大编号
const MOUSE_EVENT_MAX_TYPE event.Type = 1220

//鼠标在对象上移动
const MOUSE_MOVE_EVENT_TYPE event.Type = 1201

//鼠标从对象上移出
const MOUSE_LEAVE_EVENT_TYPE event.Type = 1202

//按下鼠标按键
const MOUSE_PRESS_EVENT_TYPE event.Type = 1203

//释放鼠标按键
const MOUSE_RELEASE_EVENT_TYPE event.Type = 1204

//鼠标双击
const MOUSE_DOUBLE_CLICK_EVENT_TYPE event.Type = 1205

//鼠标中键滚动
const MOUSE_WHEEL_EVENT_TYPE event.Type = 1206

//鼠标移入到对象上
const MOUSE_ENTER_EVENT_TYPE event.Type = 1207

/*键盘事件 1220-1230 */
//按下
const KEY_PRESS_EVENT_TYPE event.Type = 1221

//释放
const KEY_RELEASE_EVENT_TYPE event.Type = 1222

//字符输入
const KEY_CHAR_EVENT_TYPE event.Type = 1223

//重复
const KEY_REPEAT_EVENT_TYPE event.Type = 1224

/* 元素选择/取消选择事件 */
const SELECT_EVENT_TYPE event.Type = 1231

//被修改状态
const MODIFIED_EVENT_TYPE event.Type = 1232

//改变可见状态
const VISIBLE_EVENT_TYPE event.Type = 1233

//  //绘制图层事件
//  const DRAW_LAYER_EVENT_TYPE event.Type = 1234;
//  //绘制元素事件
//  const DRAW_ELEMENT_EVENT_TYPE event.Type = 1235;
//绘制元素或层事件
const PAINT_EVENT_TYPE event.Type = 1236

//元素或层事件被移动
const MOVED_EVENT_TYPE event.Type = 1237

//元素或层设置或失去焦点时事件 Focus
const FOCUS_EVENT_TYPE event.Type = 1238

//对象销毁事件
const DESTROY_EVENT_TYPE event.Type = 1239
