package tui

import (
	"errors"

	"github.com/mattn/go-runewidth"
)

type widget struct {
	renderer *Renderer
	x        int
	y        int
	width    int
	height   int
	paddingH int
	paddingV int
}

type widgetStyle struct {
	x int
	y int
}

func newWidget(renderer *Renderer, x, y, width, height, paddingH, paddingV int) (*widget, error) {
	if x+width > renderer.width || y+height > renderer.height {
		return nil, errors.New("terminal size is too small")
	}
	return &widget{renderer, x, y, width, height, paddingH, paddingV}, nil
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
			if x+width > w.width-w.paddingH*2 {
				y++
				x = 0
			}
			if y >= w.height-w.paddingV*2 {
				return
			}
			w.put(Cell{Char: r, Width: width, Style: as.Style}, x, y)
			if width == 2 {
				if as.Style.HasCursor {
					style := as.Style
					style.HasCursor = false
					w.put(Cell{Char: ' ', Width: 0, Style: style}, x+1, y)
				} else {
					w.put(Cell{Char: ' ', Width: 0, Style: as.Style}, x+1, y)
				}
			}
			x += width
		}
	}
}

func (w *widget) putBorder(style CellStyle) {
	for x := 1; x < w.width-1; x++ {
		cell := Cell{Char: '─', Width: 1, Style: style}
		w.renderer.rows[w.y][w.x+x] = cell
		w.renderer.rows[w.y+w.height-1][w.x+x] = cell
	}
	for y := 1; y < w.height-1; y++ {
		cell := Cell{Char: '│', Width: 1, Style: style}
		w.renderer.rows[w.y+y][w.x] = cell
		w.renderer.rows[w.y+y][w.x+w.width-1] = cell
	}
	w.renderer.rows[w.y][w.x] = Cell{Char: '╭', Width: 1, Style: style}
	w.renderer.rows[w.y][w.x+w.width-1] = Cell{Char: '╮', Width: 1, Style: style}
	w.renderer.rows[w.y+w.height-1][w.x] = Cell{Char: '╰', Width: 1, Style: style}
	w.renderer.rows[w.y+w.height-1][w.x+w.width-1] = Cell{Char: '╯', Width: 1, Style: style}
}

func (w *widget) putTitle(slice []Text) {
	x := 2 - w.paddingH
	for _, as := range slice {
		for _, rune_ := range as.Str {
			if rune_ == '\n' {
				return
			}
			width := RuneWidth(rune_)
			if x+width > w.width-w.paddingH*2 {
				return
			}
			w.put(Cell{Char: rune_, Width: width, Style: as.Style}, x, -w.paddingV)
			if width == 2 {
				w.put(Cell{Char: ' ', Width: 0, Style: DefaultStyle}, x+1, -w.paddingV)
			}
			x += width
		}
	}
}

func (w *widget) fill(cell Cell) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			if w.x+x > 0 && w.renderer.rows[w.y+y][w.x+x-1].Width == 2 {
				w.renderer.rows[w.y+y][w.x+x-1] =
					Cell{' ', 1, w.renderer.rows[w.y+y][w.x+x-1].Style}
			}
			w.renderer.rows[w.y+y][w.x+x] = cell
		}
	}
}

func (w *widget) put(cell Cell, x, y int) {
	w.renderer.put(cell, w.x+x+w.paddingH, w.y+y+w.paddingV)
}

func RuneWidth(r rune) int {
	switch r {
	case '※', '│':
		return 1
	}
	return runewidth.RuneWidth(r)
}
