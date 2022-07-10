package tui

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_sizeOfStyle(t *testing.T) {
	want := uintptr(16)
	t.Run(fmt.Sprintf("the size of Style{} must be less than or equal %d bytes", want), func(t *testing.T) {
		size := unsafe.Sizeof(style{})
		if size >= want {
			t.Errorf("The size of Style = %v, bigger than %v", size, want)
		}
	})
}
