package eui

import (
	//"bytes"
	"fmt"
	"testing"
	//"time"
	//"strings"
)

func Test(t *testing.T) {
	ks := NewKeySequence(Key_LControl, 'c', 'C', 'd')

	fmt.Println("NocaseTest Key_LControl:", ks.NocaseTest(Key_LControl))
	fmt.Println("NocaseTests Key_LControl, c:", ks.NocaseTests(ExactMatch, Key_LControl, 'c'))
	fmt.Println("NocaseTests Key_LControl, v:", ks.NocaseTests(ExactMatch, Key_LControl, 'v'))
	fmt.Println("NocaseTests Key_LControl, v:", ks.NocaseTests(PartialMatch, Key_LControl, 'v'))

	fmt.Println("CaseTest Key_LControl:", ks.Test(Key_LControl))
	fmt.Println("CaseTests Key_LControl, d:", ks.Tests(ExactMatch, Key_LControl, 'd'))
	fmt.Println("CaseTests Key_LControl, D:", ks.Tests(ExactMatch, Key_LControl, 'D'))

	ks.Add('V')
	fmt.Println("CaseTests Key_LControl, v:", ks.Tests(ExactMatch, Key_LControl, 'v'))
}
