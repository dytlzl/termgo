package tui

type Style struct {
	F256          uint8
	B256          uint8
	invert        bool
	bold          bool
	italic        bool
	strikethrough bool
	underline     bool
	HasCursor     bool
}

func (style *Style) merge(defaultStyle Style) {
	if style == nil {
		style = &defaultStyle
	} else {
		if style.F256 == 0 {
			style.F256 = defaultStyle.F256
		}
		if style.B256 == 0 {
			style.B256 = defaultStyle.B256
		}
		if !style.strikethrough {
			style.strikethrough = defaultStyle.strikethrough
		}
		if style.invert {
			style.B256 = defaultStyle.F256
			style.F256 = defaultStyle.B256
			if defaultStyle.B256 == 0 {
				style.F256 = 15
			}
			if defaultStyle.F256 == 0 {
				style.B256 = 15
			}
		}
	}
}
