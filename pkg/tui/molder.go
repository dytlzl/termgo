package tui

import (
	"errors"
	"fmt"

	"github.com/mattn/go-runewidth"
)

func moldView(r cellWriter, v *View, cfg *config, frame rect, parentFrame rect, defaultStyle style, allowOverflow bool) error {
	vr, err := newMolder(
		r,
		frame,
		parentFrame,
		int(v.paddingTop), int(v.paddingLeading), int(v.paddingBottom), int(v.paddingTrailing),
		allowOverflow,
	)
	if err != nil {
		return fmt.Errorf("failed to create viewRenderer: %w", err)
	}
	if v.style == nil {
		v.style = new(style)
	}
	v.style.merge(defaultStyle)
	if v.border != nil || v.title != "" || v.content != nil || v.style.b256 != 0 {
		vr.fill(cell{' ', 1, *v.style})
	}
	if v.border != nil {
		v.border.merge(*v.style)
		vr.putBorder(*v.border)
	}
	if v.title != "" {
		vr.putTitle([]text{{Str: " " + v.title + " ", Style: *v.style}})
	}
	if v.content != nil {
		vr.moldBody(v.content(), *v.style)
	}
	if v.keyHandler != nil {
		cfg.viewPQ.PushView(v)
	}

	if v.children == nil {
		return nil
	}

	availableWidth := frame.width - int(v.paddingLeading) - int(v.paddingTrailing)
	availableHeight := frame.height - int(v.paddingTop) - int(v.paddingBottom)

	remainedWidth := availableWidth
	remainedHeight := availableHeight
	numberOfAutoWidth := 0
	numberOfAutoHeight := 0

	children := v.children()

	for idx := range children {
		if children[idx] == nil {
			continue
		}

		if children[idx].style == nil {
			children[idx].style = new(style)
		}
		// calculate absolute size from relative size
		if children[idx].absoluteWidth == 0 {
			children[idx].absoluteWidth = availableWidth * int(children[idx].relativeWidth) / 12
		}
		if children[idx].absoluteHeight == 0 && v.dir == vertical && children[idx].content != nil {
			if children[idx].absoluteWidth == 0 {
				children[idx].absoluteWidth = availableWidth
			}
			children[idx].absoluteHeight = heightFromWidth(
				children[idx].content(),
				children[idx].absoluteWidth-int(children[idx].paddingLeading)-int(children[idx].paddingTrailing),
			) +
				int(children[idx].paddingTop) +
				int(children[idx].paddingBottom)
		}

		if children[idx].absoluteHeight == 0 {
			children[idx].absoluteHeight = availableHeight * int(children[idx].relativeHeight) / 12
		}

		remainedWidth -= children[idx].absoluteWidth
		remainedHeight -= children[idx].absoluteHeight

		// count auto-sizing view
		if children[idx].absoluteWidth == 0 {
			numberOfAutoWidth++
		}
		if children[idx].absoluteHeight == 0 {
			numberOfAutoHeight++
		}
	}

	accumulatedX := frame.x + int(v.paddingLeading)
	accumulatedY := frame.y + int(v.paddingTop) + v.offsetY

	for _, child := range children {
		if child == nil {
			continue
		}

		// calculate size of auto-sizing view
		if child.absoluteWidth == 0 {
			if v.dir == horizontal {
				child.absoluteWidth = remainedWidth / numberOfAutoWidth
				numberOfAutoWidth--
				remainedWidth -= child.absoluteWidth
			} else {
				child.absoluteWidth = availableWidth
			}
		}
		if child.absoluteHeight == 0 {
			if v.dir == vertical {
				child.absoluteHeight = remainedHeight / numberOfAutoHeight
				numberOfAutoHeight--
				remainedHeight -= child.absoluteHeight
			} else {
				child.absoluteHeight = availableHeight
			}
		}

		x := frame.x + int(v.paddingLeading) + (availableWidth-child.absoluteWidth)/2
		if v.dir == horizontal {
			x = accumulatedX
		}
		y := frame.y + int(v.paddingTop) + (availableHeight-child.absoluteHeight)/2
		if v.dir == vertical {
			y = accumulatedY
		}

		if y+int(v.paddingBottom) >= frame.y+frame.height {
			break
		}

		err = moldView(r, child, cfg,
			rect{
				x,
				y,
				child.absoluteWidth,
				child.absoluteHeight,
			},
			rect{
				frame.x + int(v.paddingLeading),
				frame.y + int(v.paddingTop),
				v.absoluteWidth - int(v.paddingLeading) - int(v.paddingTrailing),
				v.absoluteHeight - int(v.paddingTop) - int(v.paddingBottom),
			},
			*v.style,
			allowOverflow || v.allowOverflow)
		if err != nil {
			return err
		}
		if v.dir == horizontal {
			accumulatedX += child.absoluteWidth
		}
		if v.dir == vertical {
			accumulatedY += child.absoluteHeight
		}
	}
	return nil
}

type molder struct {
	renderer        cellWriter
	frame           rect
	parentFrame     rect
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
}

func newMolder(r cellWriter, frame rect, parentFrame rect, paddingTop, paddingLeading, paddingBottom, paddingTrailing int, allowOverflow bool) (*molder, error) {
	width, height := r.size()
	if !allowOverflow && (frame.x+frame.width > width || frame.y+frame.height > height) {
		return nil, errors.New("terminal size is too small")
	}
	return &molder{r, frame, parentFrame, paddingTop, paddingLeading, paddingBottom, paddingTrailing}, nil
}

func (m *molder) moldBody(slice []text, defaultStyle style) {
	x, y := 0, 0
	for _, as := range slice {
		as.Style.merge(defaultStyle)
		for _, r := range as.Str {
			if r == 13 { // CR
				continue
			}
			if r == 10 { // NL
				y++
				x = 0
				continue
			}
			width := runewidth.RuneWidth(r)
			if x+width > m.frame.width-m.paddingLeading-m.paddingTrailing {
				y++
				x = 0
			}
			if m.paddingTop+y+m.paddingBottom >= m.frame.height {
				return
			}
			if m.frame.y+m.paddingTop+y-1+m.paddingBottom >= m.parentFrame.y+m.parentFrame.height {
				return
			}
			if m.frame.y+m.paddingTop+y < m.parentFrame.y {
				x += width
				continue
			}
			m.put(cell{Char: r, Width: width, Style: as.Style}, x, y)
			if width == 2 {
				if as.Style.hasCursor {
					s := as.Style
					s.hasCursor = false
					m.put(cell{Char: ' ', Width: 0, Style: s}, x+1, y)
				} else {
					m.put(cell{Char: ' ', Width: 0, Style: as.Style}, x+1, y)
				}
			}
			x += width
		}
	}
}

func (m *molder) putBorder(s style) {
	if m.frame.y >= m.parentFrame.y {
		m.renderer.matrix()[m.frame.y][m.frame.x] = cell{Char: '╭', Width: 1, Style: s}
		for x := 1; x < m.frame.width-1; x++ {
			c := cell{Char: '─', Width: 1, Style: s}
			m.renderer.matrix()[m.frame.y][m.frame.x+x] = c
		}
		m.renderer.matrix()[m.frame.y][m.frame.x+m.frame.width-1] = cell{Char: '╮', Width: 1, Style: s}
	}
	for y := 1; y < m.frame.height-1; y++ {
		if m.frame.y+y < m.parentFrame.y {
			continue
		}
		if m.frame.y+y >= m.parentFrame.y+m.parentFrame.height {
			return
		}
		c := cell{Char: '│', Width: 1, Style: s}
		m.renderer.matrix()[m.frame.y+y][m.frame.x] = c
		m.renderer.matrix()[m.frame.y+y][m.frame.x+m.frame.width-1] = c
	}
	if m.frame.y+m.frame.height-1 >= m.parentFrame.y+m.parentFrame.height {
		return
	}
	if m.frame.y+m.frame.height-1 < m.parentFrame.y {
		return
	}
	m.renderer.matrix()[m.frame.y+m.frame.height-1][m.frame.x] = cell{Char: '╰', Width: 1, Style: s}
	for x := 1; x < m.frame.width-1; x++ {
		c := cell{Char: '─', Width: 1, Style: s}
		m.renderer.matrix()[m.frame.y+m.frame.height-1][m.frame.x+x] = c
	}
	m.renderer.matrix()[m.frame.y+m.frame.height-1][m.frame.x+m.frame.width-1] = cell{Char: '╯', Width: 1, Style: s}
}

func (m *molder) putTitle(slice []text) {
	if m.frame.y < m.parentFrame.y {
		return
	}
	x := 2 - m.paddingTop
	for _, as := range slice {
		for _, r := range as.Str {
			if r == '\n' {
				return
			}
			width := RuneWidth(r)
			if x+width > m.frame.width-m.paddingLeading-m.paddingTrailing {
				return
			}
			m.put(cell{Char: r, Width: width, Style: as.Style}, x, -m.paddingTop)
			if width == 2 {
				m.put(cell{Char: ' ', Width: 0}, x+1, -m.paddingTop)
			}
			x += width
		}
	}
}

func (m *molder) fill(c cell) {
	for y := 0; y < m.frame.height; y++ {
		if m.frame.y+y <= m.parentFrame.y {
			continue
		}
		if m.frame.y+y >= m.parentFrame.y+m.parentFrame.height {
			return
		}
		for x := 0; x < m.frame.width; x++ {
			if m.frame.x+x > 0 && m.renderer.matrix()[m.frame.y+y][m.frame.x+x-1].Width == 2 {
				m.renderer.matrix()[m.frame.y+y][m.frame.x+x-1] =
					cell{' ', 1, m.renderer.matrix()[m.frame.y+y][m.frame.x+x-1].Style}
			}
			m.renderer.matrix()[m.frame.y+y][m.frame.x+x] = c
		}
	}
}

func (m *molder) put(c cell, x, y int) {
	m.renderer.put(c, m.frame.x+x+m.paddingLeading, m.frame.y+y+m.paddingTop)
}

func RuneWidth(r rune) int {
	switch r {
	case '※', '│':
		return 1
	}
	return runewidth.RuneWidth(r)
}

func heightFromWidth(slice []text, width int) int {
	x, y := 0, 0
	for _, as := range slice {
		for _, r := range as.Str {
			if r == 13 { // CR
				continue
			}
			if r == 10 { // NL
				y++
				x = 0
				continue
			}
			w := runewidth.RuneWidth(r)
			if x+w > width {
				y++
				x = 0
			}
			x += w
		}
	}
	return y + 1
}
