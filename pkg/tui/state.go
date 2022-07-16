package tui

import (
	"runtime"
)

var stateMap = map[uintptr]any{}

func UseState[T any](initialState T) (T, func(T)) {
	pc, _, _, _ := runtime.Caller(1)
	if _, ok := stateMap[pc]; !ok {
		stateMap[pc] = initialState
	}
	return stateMap[pc].(T), func (newState T)  {
		stateMap[pc] = newState
	}
}
