package tui

type View interface {
	Title() string
	Body(hasFocus bool, state *GlobalState) []Text
	HandleEvent(event interface{}, state *GlobalState)
	BorderStyle() *CellStyle
	SubViews() []View
	Width() (int, int)
}

type ViewStyle struct {
	borderStyle CellStyle
	width       Fraction
}

type Fraction struct {
	numerator   int
	denominator int
}

type FooterView interface {
	Text() []Text
	Style() CellStyle
	HandleEvent(event interface{}, state *GlobalState)
}

type DefaultView struct {
}

func (*DefaultView) Title() string {
	return ""
}

func (*DefaultView) Body(bool, *GlobalState) []Text {
	return nil
}

func (*DefaultView) HandleEvent(interface{}, *GlobalState) {
}

func (*DefaultView) BorderStyle() *CellStyle {
	return nil
}

func (*DefaultView) SubViews() []View {
	return nil
}

// Width returns the numerator and the denominator of the width of a view
// math.big.Rat is not suitable, therefore this returns two int values.
func (*DefaultView) Width() (int, int) {
	return 2, 3
}
