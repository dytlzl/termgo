package tui

type View struct {
	title           string
	relativeWidth   int
	relativeHeight  int
	absoluteWidth   int
	absoluteHeight  int
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
	dir             direction
	style           *Style
	border          *Style
	children        []*View
	renderer        func(Size) []Text
	eventHandler    func(event any) any
}

type direction int

const (
	horizontal = iota + 1
	vertical
)

func ViewWithRenderer(renderer func(Size) []Text) *View {
	return &View{renderer: renderer}
}

func (v *View) RelativeSize(width, height int) *View {
	v.relativeWidth = width
	v.relativeHeight = height
	return v
}

func (v *View) AbsoluteSize(width, height int) *View {
	v.absoluteWidth = width
	v.absoluteHeight = height
	return v
}

func (v *View) Padding(top, leading, bottom, trailing int) *View {
	v.paddingTop = top
	v.paddingLeading = leading
	v.paddingBottom = bottom
	v.paddingTrailing = trailing
	return v
}

func (v *View) Title(title string) *View {
	v.title = title
	v.paddingTop = 1
	return v
}

func (v *View) Style(style Style) *View {
	v.style = &style
	return v
}

func (v *View) Border(style Style) *View {
	v.paddingTop = 2
	v.paddingLeading = 2
	v.paddingBottom = 2
	v.paddingTrailing = 2
	v.border = &style
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
	v.renderer = func(s Size) []Text { return []Text{{Str: body, Style: *v.style}} }
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
