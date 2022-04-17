package tui

type View interface {
	Body(hasFocus bool, size Size) []Text
	HandleEvent(event interface{}) interface{}
	Options() ViewOptions
}

type ViewOptions struct {
	Title       string
	SubViews    []View
	BorderStyle *CellStyle
	Width       *Fraction
}

type FooterView interface {
	Text() []Text
	Style() CellStyle
	HandleEvent(event interface{}) interface{}
}

type DefaultView struct {
}

func (*DefaultView) Body(bool, Size) []Text {
	return nil
}

func (*DefaultView) HandleEvent(interface{}) interface{} {
	return nil
}

func (*DefaultView) Options() ViewOptions {
	return ViewOptions{
		Title:       "",
		SubViews:    nil,
		BorderStyle: nil,
		Width:       NewFraction(2, 3),
	}
}

// Check if *DefaultView implements View
var _ View = new(DefaultView)
