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

// RelativeSize specifies relative width and height of units that the view used.
// The maximum number allowed is 12(100% of the parent view),
// and 0 means auto-resizing.
func (v *View) RelativeSize(width, height int) *View {
	if v == nil {
		return nil
	}
	v.relativeWidth = width
	v.relativeHeight = height
	return v
}

// AbsoluteSize specifies absolute width and height of units that the view used.
// 0 means auto-resizing.
func (v *View) AbsoluteSize(width, height int) *View {
	if v == nil {
		return nil
	}
	v.absoluteWidth = width
	v.absoluteHeight = height
	return v
}

// Padding sets padding values to the view.
// When one value is specified, it applies the same padding to all four sides.
// When two values are specified, the first padding applies to the top and bottom, the second to the left and right.
// When three values are specified, the first padding applies to the top, the second to the right and left, the third to the bottom.
// When four values are specified, the paddings apply to the top, right, bottom, and left in that order (clockwise).
func (v *View) Padding(values ...int) *View {
	if v == nil {
		return nil
	}
	top, leading, bottom, trailing := 0, 0, 0, 0
	switch len(values) {
	case 1:
		top, leading, bottom, trailing = values[0], values[0], values[0], values[0]
	case 2:
		top, leading, bottom, trailing = values[0], values[1], values[0], values[1]
	case 3:
		top, leading, bottom, trailing = values[0], values[1], values[2], values[1]
	case 4:
		top, leading, bottom, trailing = values[0], values[1], values[2], values[3]
	}
	v.paddingTop = top
	v.paddingLeading = leading
	v.paddingBottom = bottom
	v.paddingTrailing = trailing
	return v
}

// Title sets title to the view.
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

// FGColor sets a foreground color to the view.
func (v *View) FGColor(color uint8) *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.F256 = color
	return v
}

// BGColor sets a background color to the view.
func (v *View) BGColor(color uint8) *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.B256 = color
	return v
}

// Invert inverts foreground color and background color.
func (v *View) Invert(b bool) *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.invert = b
	return v
}

// Bold sets bold style to the view.
func (v *View) Bold() *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.bold = true
	return v
}

// Italic sets italic style to the view.
func (v *View) Italic() *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.italic = true
	return v
}

// Underline sets underline style to the view.
func (v *View) Underline() *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.underline = true
	return v
}

// Strikethrough sets strikethourgh style to the view.
func (v *View) Strikethrough() *View {
	if v == nil {
		return nil
	}
	if v.style == nil {
		v.style = new(Style)
	}
	v.style.strikethrough = true
	return v
}

func BorderOptionFGColor(color uint8) func(*View) {
	return func(v *View) {
		v.border.F256 = color
	}
}

func BorderOptionBGColor(color uint8) func(*View) {
	return func(v *View) {
		v.border.B256 = color
	}
}

type borderOption = func(*View)

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

func Spacer() *View {
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
			if child == nil {
				continue
			}
			if child.style == nil {
				child.style = new(Style)
			}
			child.style.merge(*view.style)
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
