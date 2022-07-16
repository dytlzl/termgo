package tui

import (
	"errors"

	"github.com/mattn/go-runewidth"
)

type viewRenderer struct {
	renderer        *renderer
	frame           rect
	parentFrame     rect
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
}

func newViewRenderer(r *renderer, frame rect, parentFrame rect, paddingTop, paddingLeading, paddingBottom, paddingTrailing int, allowOverflow bool) (*viewRenderer, error) {
	if !allowOverflow && (frame.x+frame.width > r.width || frame.y+frame.height > r.height) {
		return nil, errors.New("terminal size is too small")
	}
	return &viewRenderer{r, frame, parentFrame, paddingTop, paddingLeading, paddingBottom, paddingTrailing}, nil
}

func (w *viewRenderer) putBody(slice []text, defaultStyle style) {
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
			if x+width > w.frame.width-w.paddingLeading-w.paddingTrailing {
				y++
				// debugf("%d %d %d %d %d %d", w.frame.y, w.paddingTop, y, w.paddingBottom, w.parentFrame.y, w.parentFrame.height)
				x = 0
			}
			if w.paddingTop+y+w.paddingBottom >= w.frame.height {
				return
			}
			if w.frame.y+w.paddingTop+y-1+w.paddingBottom >= w.parentFrame.y+w.parentFrame.height {
				return
			}
			if w.frame.y+w.paddingTop+y < w.parentFrame.y {
				x += width
				continue
			}
			w.put(cell{Char: r, Width: width, Style: as.Style}, x, y)
			if width == 2 {
				if as.Style.hasCursor {
					s := as.Style
					s.hasCursor = false
					w.put(cell{Char: ' ', Width: 0, Style: s}, x+1, y)
				} else {
					w.put(cell{Char: ' ', Width: 0, Style: as.Style}, x+1, y)
				}
			}
			x += width
		}
	}
}

func (w *viewRenderer) putBorder(s style) {
	if w.frame.y >= w.parentFrame.y {
		w.renderer.rows[w.frame.y][w.frame.x] = cell{Char: '╭', Width: 1, Style: s}
		for x := 1; x < w.frame.width-1; x++ {
			c := cell{Char: '─', Width: 1, Style: s}
			w.renderer.rows[w.frame.y][w.frame.x+x] = c
		}
		w.renderer.rows[w.frame.y][w.frame.x+w.frame.width-1] = cell{Char: '╮', Width: 1, Style: s}
	}
	for y := 1; y < w.frame.height-1; y++ {
		if w.frame.y+y < w.parentFrame.y {
			continue
		}
		if w.frame.y+y >= w.parentFrame.y+w.parentFrame.height {
			return
		}
		c := cell{Char: '│', Width: 1, Style: s}
		w.renderer.rows[w.frame.y+y][w.frame.x] = c
		w.renderer.rows[w.frame.y+y][w.frame.x+w.frame.width-1] = c
	}
	if w.frame.y+w.frame.height-1 >= w.parentFrame.y+w.parentFrame.height {
		return
	}
	if w.frame.y+w.frame.height-1 < w.parentFrame.y {
		return
	}
	w.renderer.rows[w.frame.y+w.frame.height-1][w.frame.x] = cell{Char: '╰', Width: 1, Style: s}
	for x := 1; x < w.frame.width-1; x++ {
		c := cell{Char: '─', Width: 1, Style: s}
		w.renderer.rows[w.frame.y+w.frame.height-1][w.frame.x+x] = c
	}
	w.renderer.rows[w.frame.y+w.frame.height-1][w.frame.x+w.frame.width-1] = cell{Char: '╯', Width: 1, Style: s}
}

func (w *viewRenderer) putTitle(slice []text) {
	x := 2 - w.paddingTop
	for _, as := range slice {
		for _, r := range as.Str {
			if r == '\n' {
				return
			}
			width := RuneWidth(r)
			if x+width > w.frame.width-w.paddingLeading-w.paddingTrailing {
				return
			}
			w.put(cell{Char: r, Width: width, Style: as.Style}, x, -w.paddingTop)
			if width == 2 {
				w.put(cell{Char: ' ', Width: 0}, x+1, -w.paddingTop)
			}
			x += width
		}
	}
}

func (w *viewRenderer) fill(c cell) {
	for y := 0; y < w.frame.height; y++ {
		if w.frame.y+y <= w.parentFrame.y {
			continue
		}
		if w.frame.y+y >= w.parentFrame.y+w.parentFrame.height {
			return
		}
		for x := 0; x < w.frame.width; x++ {
			if w.frame.x+x > 0 && w.renderer.rows[w.frame.y+y][w.frame.x+x-1].Width == 2 {
				w.renderer.rows[w.frame.y+y][w.frame.x+x-1] =
					cell{' ', 1, w.renderer.rows[w.frame.y+y][w.frame.x+x-1].Style}
			}
			w.renderer.rows[w.frame.y+y][w.frame.x+x] = c
		}
	}
}

func (w *viewRenderer) put(c cell, x, y int) {
	w.renderer.put(c, w.frame.x+x+w.paddingLeading, w.frame.y+y+w.paddingTop)
}

func RuneWidth(r rune) int {
	switch r {
	case '※', '│':
		return 1
	}
	return runewidth.RuneWidth(r)
}
