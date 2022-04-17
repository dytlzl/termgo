package tui

import "reflect"

type GlobalState struct {
	Width               int
	Height              int
	FocusedModelType    reflect.Type
	ShouldTerminate     bool
	ShouldSkipRendering bool
}
