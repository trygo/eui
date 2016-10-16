package eui

import (
	"github.com/google/gxui"
)

const (
	//NoMatch      = iota //不匹配
	PartialMatch = iota //部分匹配, 即包含
	ExactMatch          //精确匹配, 完全相等
)

type KeySequence struct {
	Keys []gxui.KeyboardKey
}

func NewKeySequence(keys ...gxui.KeyboardKey) *KeySequence {
	return &KeySequence{Keys: keys}
}

func (ks *KeySequence) Add(keys ...gxui.KeyboardKey) {
	if len(keys) > 0 {
		ks.Keys = append(ks.Keys, keys...)
	}
}

func (ks *KeySequence) At(index int) gxui.KeyboardKey {
	return ks.Keys[index]
}

func (ks *KeySequence) Size() int {
	return len(ks.Keys)
}

func (ks *KeySequence) Empty() bool {
	return len(ks.Keys) == 0
}

//不会忽略大小写
func (ks *KeySequence) Test(key gxui.KeyboardKey) bool {
	for _, v := range ks.Keys {
		if v == key {
			return true
		}
	}
	return false
}

////忽略大小写
//func (ks *KeySequence) Test(key gxui.KeyboardKey) bool {
//	key = toLower(key)
//	for _, v := range ks.Keys {
//		if toLower(v) == key {
//			return true
//		}
//	}
//	return false
//}

////不会忽略大小写
//func (ks *KeySequence) NocaseTests(matchMode int, keys ...byte) bool {
//	if matchMode == PartialMatch {
//		for _, key := range keys {
//			if !ks.NocaseTest(key) {
//				return false
//			}
//		}
//		return true
//	} else if matchMode == ExactMatch {
//		if len(ks.Keys) != len(keys) {
//			return false
//		}
//		for _, key := range keys {
//			if !ks.NocaseTest(key) {
//				return false
//			}
//		}
//		return true
//	} else {
//		return false
//	}
//}

func (ks *KeySequence) Tests(matchMode int, keys ...gxui.KeyboardKey) bool {
	if matchMode == PartialMatch {
		for _, key := range keys {
			if !ks.Test(key) {
				return false
			}
		}
		return true
	} else if matchMode == ExactMatch {
		if len(ks.Keys) != len(keys) {
			return false
		}
		for _, key := range keys {
			if !ks.Test(key) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

//func toLower(c byte) byte {
//	if 'A' <= c && c <= 'Z' {
//		c += ('a' - 'A')
//	}
//	return c
//}
