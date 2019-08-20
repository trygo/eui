package eui

type KeyboardKey int

const (
	KeyUnknown KeyboardKey = iota
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKp0
	KeyKp1
	KeyKp2
	KeyKp3
	KeyKp4
	KeyKp5
	KeyKp6
	KeyKp7
	KeyKp8
	KeyKp9
	KeyKpDecimal
	KeyKpDivide
	KeyKpMultiply
	KeyKpSubtract
	KeyKpAdd
	KeyKpEnter
	KeyKpEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
	KeyLast
)

type Key int

const (
	/* The unknown key */
	KEY_UNKNOWN = -1
	/* Printable keys */
	KEY_SPACE         = 32
	KEY_APOSTROPHE    = 39
	KEY_COMMA         = 44
	KEY_MINUS         = 45
	KEY_PERIOD        = 46
	KEY_SLASH         = 47
	KEY_0             = 48
	KEY_1             = 49
	KEY_2             = 50
	KEY_3             = 51
	KEY_4             = 52
	KEY_5             = 53
	KEY_6             = 54
	KEY_7             = 55
	KEY_8             = 56
	KEY_9             = 57
	KEY_SEMICOLON     = 59
	KEY_EQUAL         = 61
	KEY_A             = 65
	KEY_B             = 66
	KEY_C             = 67
	KEY_D             = 68
	KEY_E             = 69
	KEY_F             = 70
	KEY_G             = 71
	KEY_H             = 72
	KEY_I             = 73
	KEY_J             = 74
	KEY_K             = 75
	KEY_L             = 76
	KEY_M             = 77
	KEY_N             = 78
	KEY_O             = 79
	KEY_P             = 80
	KEY_Q             = 81
	KEY_R             = 82
	KEY_S             = 83
	KEY_T             = 84
	KEY_U             = 85
	KEY_V             = 86
	KEY_W             = 87
	KEY_X             = 88
	KEY_Y             = 89
	KEY_Z             = 90
	KEY_LEFT_BRACKET  = 91
	KEY_BACKSLASH     = 92
	KEY_RIGHT_BRACKET = 93
	KEY_GRAVE_ACCENT  = 96
	KEY_WORLD_1       = 161
	KEY_WORLD_2       = 162

	//Function       keys
	KEY_ESCAPE        = 256
	KEY_ENTER         = 257
	KEY_TAB           = 258
	KEY_BACKSPACE     = 259
	KEY_INSERT        = 260
	KEY_DELETE        = 261
	KEY_RIGHT         = 262
	KEY_LEFT          = 263
	KEY_DOWN          = 264
	KEY_UP            = 265
	KEY_PAGE_UP       = 266
	KEY_PAGE_DOWN     = 267
	KEY_HOME          = 268
	KEY_END           = 269
	KEY_CAPS_LOCK     = 280
	KEY_SCROLL_LOCK   = 281
	KEY_NUM_LOCK      = 282
	KEY_PRINT_SCREEN  = 283
	KEY_PAUSE         = 284
	KEY_F1            = 290
	KEY_F2            = 291
	KEY_F3            = 292
	KEY_F4            = 293
	KEY_F5            = 294
	KEY_F6            = 295
	KEY_F7            = 296
	KEY_F8            = 297
	KEY_F9            = 298
	KEY_F10           = 299
	KEY_F11           = 300
	KEY_F12           = 301
	KEY_F13           = 302
	KEY_F14           = 303
	KEY_F15           = 304
	KEY_F16           = 305
	KEY_F17           = 306
	KEY_F18           = 307
	KEY_F19           = 308
	KEY_F20           = 309
	KEY_F21           = 310
	KEY_F22           = 311
	KEY_F23           = 312
	KEY_F24           = 313
	KEY_F25           = 314
	KEY_KP_0          = 320
	KEY_KP_1          = 321
	KEY_KP_2          = 322
	KEY_KP_3          = 323
	KEY_KP_4          = 324
	KEY_KP_5          = 325
	KEY_KP_6          = 326
	KEY_KP_7          = 327
	KEY_KP_8          = 328
	KEY_KP_9          = 329
	KEY_KP_DECIMAL    = 330
	KEY_KP_DIVIDE     = 331
	KEY_KP_MULTIPLY   = 332
	KEY_KP_SUBTRACT   = 333
	KEY_KP_ADD        = 334
	KEY_KP_ENTER      = 335
	KEY_KP_EQUAL      = 336
	KEY_LEFT_SHIFT    = 340
	KEY_LEFT_CONTROL  = 341
	KEY_LEFT_ALT      = 342
	KEY_LEFT_SUPER    = 343
	KEY_RIGHT_SHIFT   = 344
	KEY_RIGHT_CONTROL = 345
	KEY_RIGHT_ALT     = 346
	KEY_RIGHT_SUPER   = 347
	KEY_MENU          = 348

	KEY_LAST = KEY_MENU
)

func TranslateKeyboardKey(in Key) KeyboardKey {
	switch in {
	case KEY_UNKNOWN:
		return KeyUnknown
	case KEY_SPACE:
		return KeySpace
	case KEY_APOSTROPHE:
		return KeyApostrophe
	case KEY_COMMA:
		return KeyComma
	case KEY_MINUS:
		return KeyMinus
	case KEY_PERIOD:
		return KeyPeriod
	case KEY_SLASH:
		return KeySlash
	case KEY_0:
		return Key0
	case KEY_1:
		return Key1
	case KEY_2:
		return Key2
	case KEY_3:
		return Key3
	case KEY_4:
		return Key4
	case KEY_5:
		return Key5
	case KEY_6:
		return Key6
	case KEY_7:
		return Key7
	case KEY_8:
		return Key8
	case KEY_9:
		return Key9
	case KEY_SEMICOLON:
		return KeySemicolon
	case KEY_EQUAL:
		return KeyEqual
	case KEY_A:
		return KeyA
	case KEY_B:
		return KeyB
	case KEY_C:
		return KeyC
	case KEY_D:
		return KeyD
	case KEY_E:
		return KeyE
	case KEY_F:
		return KeyF
	case KEY_G:
		return KeyG
	case KEY_H:
		return KeyH
	case KEY_I:
		return KeyI
	case KEY_J:
		return KeyJ
	case KEY_K:
		return KeyK
	case KEY_L:
		return KeyL
	case KEY_M:
		return KeyM
	case KEY_N:
		return KeyN
	case KEY_O:
		return KeyO
	case KEY_P:
		return KeyP
	case KEY_Q:
		return KeyQ
	case KEY_R:
		return KeyR
	case KEY_S:
		return KeyS
	case KEY_T:
		return KeyT
	case KEY_U:
		return KeyU
	case KEY_V:
		return KeyV
	case KEY_W:
		return KeyW
	case KEY_X:
		return KeyX
	case KEY_Y:
		return KeyY
	case KEY_Z:
		return KeyZ
	case KEY_LEFT_BRACKET:
		return KeyLeftBracket
	case KEY_BACKSLASH:
		return KeyBackslash
	case KEY_RIGHT_BRACKET:
		return KeyRightBracket
	case KEY_GRAVE_ACCENT:
		return KeyGraveAccent
	case KEY_WORLD_1:
		return KeyWorld1
	case KEY_WORLD_2:
		return KeyWorld2
	case KEY_ESCAPE:
		return KeyEscape
	case KEY_ENTER:
		return KeyEnter
	case KEY_TAB:
		return KeyTab
	case KEY_BACKSPACE:
		return KeyBackspace
	case KEY_INSERT:
		return KeyInsert
	case KEY_DELETE:
		return KeyDelete
	case KEY_RIGHT:
		return KeyRight
	case KEY_LEFT:
		return KeyLeft
	case KEY_DOWN:
		return KeyDown
	case KEY_UP:
		return KeyUp
	case KEY_PAGE_UP:
		return KeyPageUp
	case KEY_PAGE_DOWN:
		return KeyPageDown
	case KEY_HOME:
		return KeyHome
	case KEY_END:
		return KeyEnd
	case KEY_CAPS_LOCK:
		return KeyCapsLock
	case KEY_SCROLL_LOCK:
		return KeyScrollLock
	case KEY_NUM_LOCK:
		return KeyNumLock
	case KEY_PRINT_SCREEN:
		return KeyPrintScreen
	case KEY_PAUSE:
		return KeyPause
	case KEY_F1:
		return KeyF1
	case KEY_F2:
		return KeyF2
	case KEY_F3:
		return KeyF3
	case KEY_F4:
		return KeyF4
	case KEY_F5:
		return KeyF5
	case KEY_F6:
		return KeyF6
	case KEY_F7:
		return KeyF7
	case KEY_F8:
		return KeyF8
	case KEY_F9:
		return KeyF9
	case KEY_F10:
		return KeyF10
	case KEY_F11:
		return KeyF11
	case KEY_F12:
		return KeyF12
	case KEY_F13:
		return KeyF13
	case KEY_F14:
		return KeyF14
	case KEY_F15:
		return KeyF15
	case KEY_F16:
		return KeyF16
	case KEY_F17:
		return KeyF17
	case KEY_F18:
		return KeyF18
	case KEY_F19:
		return KeyF19
	case KEY_F20:
		return KeyF20
	case KEY_F21:
		return KeyF21
	case KEY_F22:
		return KeyF22
	case KEY_F23:
		return KeyF23
	case KEY_F24:
		return KeyF24
	case KEY_F25:
		return KeyF25
	case KEY_KP_0:
		return KeyKp0
	case KEY_KP_1:
		return KeyKp1
	case KEY_KP_2:
		return KeyKp2
	case KEY_KP_3:
		return KeyKp3
	case KEY_KP_4:
		return KeyKp4
	case KEY_KP_5:
		return KeyKp5
	case KEY_KP_6:
		return KeyKp6
	case KEY_KP_7:
		return KeyKp7
	case KEY_KP_8:
		return KeyKp8
	case KEY_KP_9:
		return KeyKp9
	case KEY_KP_DECIMAL:
		return KeyKpDecimal
	case KEY_KP_DIVIDE:
		return KeyKpDivide
	case KEY_KP_MULTIPLY:
		return KeyKpMultiply
	case KEY_KP_SUBTRACT:
		return KeyKpSubtract
	case KEY_KP_ADD:
		return KeyKpAdd
	case KEY_KP_ENTER:
		return KeyKpEnter
	case KEY_KP_EQUAL:
		return KeyKpEqual
	case KEY_LEFT_SHIFT:
		return KeyLeftShift
	case KEY_LEFT_CONTROL:
		return KeyLeftControl
	case KEY_LEFT_ALT:
		return KeyLeftAlt
	case KEY_LEFT_SUPER:
		return KeyLeftSuper
	case KEY_RIGHT_SHIFT:
		return KeyRightShift
	case KEY_RIGHT_CONTROL:
		return KeyRightControl
	case KEY_RIGHT_ALT:
		return KeyRightAlt
	case KEY_RIGHT_SUPER:
		return KeyRightSuper
	case KEY_MENU:
		return KeyMenu

	default:
		return KeyUnknown
	}
}

var Keys []Key = []Key{ //KEY_UNKNOWN,
	KEY_SPACE,
	KEY_APOSTROPHE,
	KEY_COMMA,
	KEY_MINUS,
	KEY_PERIOD,
	KEY_SLASH,
	KEY_0,
	KEY_1,
	KEY_2,
	KEY_3,
	KEY_4,
	KEY_5,
	KEY_6,
	KEY_7,
	KEY_8,
	KEY_9,
	KEY_SEMICOLON,
	KEY_EQUAL,
	KEY_A,
	KEY_B,
	KEY_C,
	KEY_D,
	KEY_E,
	KEY_F,
	KEY_G,
	KEY_H,
	KEY_I,
	KEY_J,
	KEY_K,
	KEY_L,
	KEY_M,
	KEY_N,
	KEY_O,
	KEY_P,
	KEY_Q,
	KEY_R,
	KEY_S,
	KEY_T,
	KEY_U,
	KEY_V,
	KEY_W,
	KEY_X,
	KEY_Y,
	KEY_Z,
	KEY_LEFT_BRACKET,
	KEY_BACKSLASH,
	KEY_RIGHT_BRACKET,
	KEY_GRAVE_ACCENT,
	KEY_WORLD_1,
	KEY_WORLD_2,
	KEY_ESCAPE,
	KEY_ENTER,
	KEY_TAB,
	KEY_BACKSPACE,
	KEY_INSERT,
	KEY_DELETE,
	KEY_RIGHT,
	KEY_LEFT,
	KEY_DOWN,
	KEY_UP,
	KEY_PAGE_UP,
	KEY_PAGE_DOWN,
	KEY_HOME,
	KEY_END,
	KEY_CAPS_LOCK,
	KEY_SCROLL_LOCK,
	KEY_NUM_LOCK,
	KEY_PRINT_SCREEN,
	KEY_PAUSE,
	KEY_F1,
	KEY_F2,
	KEY_F3,
	KEY_F4,
	KEY_F5,
	KEY_F6,
	KEY_F7,
	KEY_F8,
	KEY_F9,
	KEY_F10,
	KEY_F11,
	KEY_F12,
	KEY_F13,
	KEY_F14,
	KEY_F15,
	KEY_F16,
	KEY_F17,
	KEY_F18,
	KEY_F19,
	KEY_F20,
	KEY_F21,
	KEY_F22,
	KEY_F23,
	KEY_F24,
	KEY_F25,
	KEY_KP_0,
	KEY_KP_1,
	KEY_KP_2,
	KEY_KP_3,
	KEY_KP_4,
	KEY_KP_5,
	KEY_KP_6,
	KEY_KP_7,
	KEY_KP_8,
	KEY_KP_9,
	KEY_KP_DECIMAL,
	KEY_KP_DIVIDE,
	KEY_KP_MULTIPLY,
	KEY_KP_SUBTRACT,
	KEY_KP_ADD,
	KEY_KP_ENTER,
	KEY_KP_EQUAL,
	KEY_LEFT_SHIFT,
	KEY_LEFT_CONTROL,
	KEY_LEFT_ALT,
	KEY_LEFT_SUPER,
	KEY_RIGHT_SHIFT,
	KEY_RIGHT_CONTROL,
	KEY_RIGHT_ALT,
	KEY_RIGHT_SUPER,
	KEY_MENU,
	KEY_LAST,
}
