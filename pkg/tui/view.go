package tui

import (
	"fmt"
)

type View struct {
	relativeWidth   int
	relativeHeight  int
	absoluteWidth   int
	absoluteHeight  int
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
	title           string
	dir             direction
	style           *Style
	border          *Style
	children        []*View
	renderer        func() []text
}

type direction int

const (
	horizontal = iota + 1
	vertical
)

func (v *View) RelativeSize(width, height int) *View {
	if v == nil {
		return nil
	}
	v.relativeWidth = width
	v.relativeHeight = height
	return v
}

func (v *View) AbsoluteSize(width, height int) *View {
	if v == nil {
		return nil
	}
	v.absoluteWidth = width
	v.absoluteHeight = height
	return v
}

func (v *View) Padding(top, leading, bottom, trailing int) *View {
	if v == nil {
		return nil
	}
	v.paddingTop = top
	v.paddingLeading = leading
	v.paddingBottom = bottom
	v.paddingTrailing = trailing
	return v
}

func (v *View) Title(title string) *View {
	if v == nil {
		return nil
	}
	v.title = title
	v.paddingTop = 2
	return v
}

func (v *View) Style(style Style) *View {
	if v == nil {
		return nil
	}
	v.style = &style
	return v
}

func (v *View) FGColor(color int) *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.F256 = color
	return v
}

func (v *View) BGColor(color int) *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.B256 = color
	return v
}

func BorderOptionFGColor(color int) func(*View) {
	return func(v *View) {
		v.border.F256 = color
	}
}

func BorderOptionBGColor(color int) func(*View) {
	return func(v *View) {
		v.border.B256 = color
	}
}

type borderOption = func(*View) error

func (v *View) Border(options ...borderOption) *View {
	if v == nil {
		return nil
	}
	v.paddingTop = 2
	v.paddingLeading = 2
	v.paddingBottom = 2
	v.paddingTrailing = 2
	v.border = new(Style)
	for _, option := range options {
		option(v)
	}
	return v
}

func (v *View) Hidden(isHidden bool) *View {
	if isHidden {
		return nil
	} else {
		return v
	}
}

func TextView(body string) *View {
	v := &View{}
	v.renderer = func() []text { return []text{{Str: body, Style: *v.style}} }
	return v
}

func Spacer(views ...*View) *View {
	return &View{}
}

func HStack(views ...*View) *View {
	return &View{children: views, dir: horizontal}
}

func VStack(views ...*View) *View {
	return &View{children: views, dir: vertical}
}

func ZStack(views ...*View) *View {
	return &View{children: views}
}

func String(s string) *View {
	view := &View{}
	view.renderer = func() []text {
		return []text{{Str: s, Style: *view.style}}
	}
	return view
}

func Cursor(s string) *View {
	v := String(s)
	v.style = new(Style)
	v.style.HasCursor = true
	return v
}

func Fmt(format string, a ...any) *View {
	return String(fmt.Sprintf(format, a...))
}

func Break() *View {
	return String("\n")
}

func InlineStack(views ...*View) *View {
	view := &View{}
	view.renderer = func() []text {
		slice := make([]text, 0)
		for _, child := range views {
			if child.style == nil {
				child.style = new(Style)
			}
			mergeDefaultStyle(child.style, *view.style)
			slice = append(slice, child.renderer()...)
		}
		return slice
	}
	return view
}

func Map[T1 any, T2 any](slice []T1, fn func(T1) T2) []T2 {
	slice2 := make([]T2, len(slice))
	for idx, element := range slice {
		slice2[idx] = fn(element)
	}
	return slice2
}

func MapN[T any](number int, fn func(int) T) []T {
	slice := make([]T, number)
	for i := 0; i < number; i++ {
		slice[i] = fn(i)
	}
	return slice
}

func HMap[T any](slice []T, fn func(T) *View) *View {
	return HStack(Map(slice, fn)...)
}

func VMap[T any](slice []T, fn func(T) *View) *View {
	return VStack(Map(slice, fn)...)
}

func ZMap[T any](slice []T, fn func(T) *View) *View {
	return VStack(Map(slice, fn)...)
}

func InlineMap[T any](slice []T, fn func(T) *View) *View {
	return InlineStack(Map(slice, fn)...)
}

func InlineMapN(number int, fn func(int) *View) *View {
	return InlineStack(MapN(number, fn)...)
}

func If[T any](condition bool, t T, f T) T {
	if condition {
		return t
	} else {
		return f
	}
}
