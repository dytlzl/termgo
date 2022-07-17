package tui

import (
	"testing"
)

func Test_UseRef(t *testing.T) {
	funcWithUseRef := func() int {
		p := UseRef(3)
		*p++
		return *p
	}
	t.Run("retain the state", func(t *testing.T) {
		got1 := funcWithUseRef()
		want1 := 4
		if got1 != want1 {
			t.Errorf("useRef() got key = %v, want %v", got1, want1)
		}
		got2 := funcWithUseRef()
		want2 := 5
		if got2 != want2 {
			t.Errorf("useRef() got key = %v, want %v", got2, want2)
		}
	})
}
