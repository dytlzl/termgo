package tui

import (
	"testing"
	"unsafe"
)

func Test_sizeOfStyle(t *testing.T) {
	t.Run("the size of Style{} must be 8 bytes", func(t *testing.T) {
		want := uintptr(8)
		size := unsafe.Sizeof(Style{})
		if size != want {
			t.Errorf("The size of Style = %v, want %v", size, want)
		}
	})
}
