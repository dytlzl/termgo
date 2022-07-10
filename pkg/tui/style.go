package tui

type style struct {
	f256          uint8
	b256          uint8
	bold          bool
	italic        bool
	strikethrough bool
	underline     bool
	reverse       bool
	hasCursor     bool
}

func (s *style) merge(defaultStyle style) {
	if s == nil {
		s = &defaultStyle
	} else {
		if s.f256 == 0 {
			s.f256 = defaultStyle.f256
		}
		if s.b256 == 0 {
			s.b256 = defaultStyle.b256
		}
		if !s.reverse {
			s.reverse = defaultStyle.reverse
		}
	}
}
