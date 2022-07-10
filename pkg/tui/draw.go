package tui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type cell struct {
	Char  rune
	Width int
	Style style
}

type text struct {
	Str   string
	Style style
}

func TermSize() (width int, height int, err error) {
	return term.GetSize(int(os.Stdin.Fd()))
}

func (r *renderer) updateTerminalSize() (bool, error) {
	width, height, err := term.GetSize(r.fd())
	if !r.isAlternative {
		height = r.height
	}
	if err != nil {
		return false, err
	}
	hasChanged := r.width != width || r.height != height
	if hasChanged {
		r.width = width
		r.height = height
		r.rows = make([][]cell, r.height)
		for y := 0; y < r.height; y++ {
			r.rows[y] = make([]cell, r.width)
		}
	}
	return hasChanged, nil
}

func (r *renderer) fill(s style) {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.rows[y][x] = cell{' ', 1, s}
		}
	}
}

func (r *renderer) put(c cell, x, y int) {
	if r.rows[y][x] != c {
		r.rows[y][x] = c
	}
}

func (r *renderer) draw() {
	if r.isAlternative {
		origin()
	} else {
		csi(fmt.Sprintf("%dA", r.cursorY-1))
		push("\r")
	}
	lastStyle := style{}
	for y := 0; y < r.height; y++ {
		if y != 0 {
			csi("1B")
			push("\r")
		}
		for x := 0; x < r.width; x++ {
			s := r.rows[y][x].Style
			if s != lastStyle {
				push("\033[1;0m")
				if s.f256 != 0 {
					push(fmt.Sprintf("\033[38;5;%dm", s.f256))
				}
				if s.b256 != 0 {
					push(fmt.Sprintf("\033[48;5;%dm", s.b256))
				}
				if s.bold {
					push("\033[1m")
				}
				if s.italic {
					push("\033[3m")
				}
				if s.underline {
					push("\033[4m")
				}
				if s.reverse {
					push("\033[7m")
				}
				if s.strikethrough {
					push("\033[9m")
				}

				lastStyle = s
			}

			if s.hasCursor {
				r.cursorY = y
				r.cursorX = x
			}
			if !(r.rows[y][x].Width == 0 && r.rows[y][x-1].Width == 2) {
				push(string(r.rows[y][x].Char))
			}
		}
	}
	push("\033[1;0m") // Reset Style
	if r.isAlternative {
		origin()
	} else {
		csi(fmt.Sprintf("%dA", r.height-1))
		push("\r")
	}
	csi(fmt.Sprintf("%dB", r.cursorY))
	csi(fmt.Sprintf("%dC", r.cursorX))
	flush()
}
