package tui

import (
	"runtime"
)

var stateContainer = map[uintptr]any{}

func UseState[T any](initialState T) (T, func(T)) {
	return useState(initialState, 2)
}

func useState[T any](initialState T, skip int) (T, func(T)) {
	pc, _, _, _ := runtime.Caller(skip)
	if _, ok := stateContainer[pc]; !ok {
		stateContainer[pc] = initialState
	}
	return stateContainer[pc].(T), func(newState T) {
		stateContainer[pc] = newState
	}
}

func UseRef[T any](initialState T) *T {
	return useRef(initialState, 2)
}

func useRef[T any](initialState T, skip int) *T {
	pc, _, _, _ := runtime.Caller(skip)
	if _, ok := stateContainer[pc]; !ok {
		stateContainer[pc] = &initialState
	}
	return stateContainer[pc].(*T)
}
