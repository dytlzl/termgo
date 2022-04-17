package tui

import (
	"fmt"

	"golang.org/x/term"
)

type Cell struct {
	Char  rune
	Width int
	Style CellStyle
}

type Text struct {
	Str   string
	Style CellStyle
}

type CellStyle struct {
	Foreground int
	Background int
	F256       int
	B256       int
	HasCursor  bool
}

type Option struct {
	Style  CellStyle
	Footer FooterView
}

var DefaultStyle = CellStyle{}

func (r *Renderer) UpdateTerminalSize() (bool, error) {
	width, height, err := term.GetSize(r.fd())
	if err != nil {
		return false, err
	}
	hasChanged := r.width != width || r.height != height
	if hasChanged {
		r.width = width
		r.height = height
		r.rows = make([][]Cell, r.height)
		for y := 0; y < r.height; y++ {
			r.rows[y] = make([]Cell, r.width)
		}
	}
	return hasChanged, nil
}

func (r *Renderer) Fill(style CellStyle) {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.rows[y][x] = Cell{' ', 1, style}
		}
	}
}

func (r *Renderer) put(cell Cell, x, y int) {
	if r.rows[y][x] != cell {
		r.rows[y][x] = cell
	}
}

func (r *Renderer) Draw() {
	origin()
	lastStyle := DefaultStyle
	for y := 0; y < r.height; y++ {
		if y != 0 {
			csi("1B")
			push("\r")
		}
		for x := 0; x < r.width; x++ {
			style := r.rows[y][x].Style
			if lastStyle.Foreground != style.Foreground ||
				lastStyle.Background != style.Background ||
				lastStyle.F256 != style.F256 ||
				lastStyle.B256 != style.B256 {
				push("\033[1;0m")
				if style.F256 != 0 {
					push(fmt.Sprintf("\033[38;5;%dm", style.F256))
				} else if style.Foreground != 0 {
					push(fmt.Sprintf("\033[1;%dm", style.Foreground))
				}
				if style.B256 != 0 {
					push(fmt.Sprintf("\033[48;5;%dm", style.B256))
				} else if style.Background != 0 {
					push(fmt.Sprintf("\033[1;%dm", style.Background))
				}
				lastStyle = style
			}
			if style.HasCursor {
				r.cursorY = y
				r.cursorX = x
			}
			if !(r.rows[y][x].Width == 0 && r.rows[y][x-1].Width == 2) {
				push(string(r.rows[y][x].Char))
			}
		}
	}
	push("\033[1;0m") // Reset Style
	origin()
	csi(fmt.Sprintf("%dB", r.cursorY))
	csi(fmt.Sprintf("%dC", r.cursorX))
	flush()
}
