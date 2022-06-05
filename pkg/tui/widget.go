package tui

import (
	"errors"

	"github.com/mattn/go-runewidth"
)

type widget struct {
	renderer        *Renderer
	x               int
	y               int
	width           int
	height          int
	paddingTop      int
	paddingLeading  int
	paddingBottom   int
	paddingTrailing int
}

type Size struct {
	Width  int
	Height int
}

func newWidget(renderer *Renderer, x, y, width, height, paddingTop, paddingLeading, paddingBottom, paddingTrailing int) (*widget, error) {
	if x+width > renderer.width || y+height > renderer.height {
		return nil, errors.New("terminal size is too small")
	}
	return &widget{renderer, x, y, width, height, paddingTop, paddingLeading, paddingBottom, paddingTrailing}, nil
}

func (w *widget) putBody(slice []Text) {
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
			width := runewidth.RuneWidth(r)
			if x+width > w.width-w.paddingLeading-w.paddingTrailing {
				y++
				x = 0
			}
			if y >= w.height-w.paddingTop-w.paddingBottom {
				return
			}
			w.put(cell{Char: r, Width: width, Style: as.Style}, x, y)
			if width == 2 {
				if as.Style.HasCursor {
					style := as.Style
					style.HasCursor = false
					w.put(cell{Char: ' ', Width: 0, Style: style}, x+1, y)
				} else {
					w.put(cell{Char: ' ', Width: 0, Style: as.Style}, x+1, y)
				}
			}
			x += width
		}
	}
}

func (w *widget) putBorder(style Style) {
	for x := 1; x < w.width-1; x++ {
		c := cell{Char: '─', Width: 1, Style: style}
		w.renderer.rows[w.y][w.x+x] = c
		w.renderer.rows[w.y+w.height-1][w.x+x] = c
	}
	for y := 1; y < w.height-1; y++ {
		c := cell{Char: '│', Width: 1, Style: style}
		w.renderer.rows[w.y+y][w.x] = c
		w.renderer.rows[w.y+y][w.x+w.width-1] = c
	}
	w.renderer.rows[w.y][w.x] = cell{Char: '╭', Width: 1, Style: style}
	w.renderer.rows[w.y][w.x+w.width-1] = cell{Char: '╮', Width: 1, Style: style}
	w.renderer.rows[w.y+w.height-1][w.x] = cell{Char: '╰', Width: 1, Style: style}
	w.renderer.rows[w.y+w.height-1][w.x+w.width-1] = cell{Char: '╯', Width: 1, Style: style}
}

func (w *widget) putTitle(slice []Text) {
	x := 2 - w.paddingTop
	for _, as := range slice {
		for _, rune_ := range as.Str {
			if rune_ == '\n' {
				return
			}
			width := RuneWidth(rune_)
			if x+width > w.width-w.paddingLeading-w.paddingTrailing {
				return
			}
			w.put(cell{Char: rune_, Width: width, Style: as.Style}, x, -w.paddingTop)
			if width == 2 {
				w.put(cell{Char: ' ', Width: 0, Style: DefaultStyle}, x+1, -w.paddingTop)
			}
			x += width
		}
	}
}

func (w *widget) fill(c cell) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			if w.x+x > 0 && w.renderer.rows[w.y+y][w.x+x-1].Width == 2 {
				w.renderer.rows[w.y+y][w.x+x-1] =
					cell{' ', 1, w.renderer.rows[w.y+y][w.x+x-1].Style}
			}
			w.renderer.rows[w.y+y][w.x+x] = c
		}
	}
}

func (w *widget) put(c cell, x, y int) {
	w.renderer.put(c, w.x+x+w.paddingLeading, w.y+y+w.paddingTop)
}

func RuneWidth(r rune) int {
	switch r {
	case '※', '│':
		return 1
	}
	return runewidth.RuneWidth(r)
}
