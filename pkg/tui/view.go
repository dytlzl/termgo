package tui

type View struct {
	title           string
	x               int
	y               int
	reverseH        bool
	reverseV        bool
	relativeWidth   int
	relativeHeight  int
	absoluteWidth   int
	absoluteHeight  int
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
	style           *Style
	border          *Style
	children        []*View
	renderer        func(Size) []Text
	eventHandler    func(event any) any
}

func ViewWithRenderer(renderer func(Size) []Text) *View {
	return &View{renderer: renderer, relativeWidth: 12, relativeHeight: 12}
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

func (v *View) Position(x, y int) *View {
	v.x = x
	v.y = y
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

func TextView(body string) *View {
	v := &View{}
	v.renderer = func(s Size) []Text { return []Text{{Str: body, Style: *v.style}} }
	return v
}

func HStack(views ...*View) *View {
	x := 0
	for i, v := range views {
		if views[i] == nil {
			continue
		}
		views[i].x = x
		x += v.relativeWidth
	}
	return &View{children: views}
}

func VStack(views ...*View) *View {
	y := 0
	for i, v := range views {
		if views[i] == nil {
			continue
		}
		views[i].y = y
		y += v.relativeHeight
	}
	return &View{children: views}
}

func ReversedHStack(views ...*View) *View {
	v := HStack(views...)
	v.reverseH = true
	return v
}

func ReversedVStack(views ...*View) *View {
	v := VStack(views...)
	v.reverseV = true
	return v
}

func ZStack(views ...*View) *View {
	for i := range views {
		if views[i] == nil {
			continue
		}
		if views[i].relativeWidth == 0 {
			views[i].relativeWidth = 12
		}
		if views[i].relativeHeight == 0 {
			views[i].relativeHeight = 12
		}
		if views[i].x == 0 && views[i].y == 0 {
			views[i].x = (12 - views[i].relativeWidth) / 2
			views[i].y = (12 - views[i].relativeHeight) / 2
		}
	}
	return &View{children: views}
}
