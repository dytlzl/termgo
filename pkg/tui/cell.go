package tui

import (
	"fmt"
	"os"
	"strings"

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

type cellWriter interface {
	size() (int, int)
	matrix() [][]cell
	put(c cell, x, y int)
}

type generalCellWriter struct {
	isAlternative bool
	width         int
	height        int
	rows          [][]cell
	cursorX       int
	cursorY       int
	eventChan     chan any
	ttyin         *os.File
	oldState      *term.State
}

func newGeneralCellWriter(isAlternative bool) (*generalCellWriter, error) {
	ttyin := os.Stdin
	state, err := term.MakeRaw(int(ttyin.Fd()))
	if err != nil {
		return nil, err
	}
	height := 20
	if !isAlternative {
		push(strings.Repeat("\n", height-2))
		flush()
	}
	initRenderer(isAlternative)
	return &generalCellWriter{
		isAlternative: isAlternative,
		height:        If(isAlternative, 0, height),
		cursorY:       If(isAlternative, 0, height-1),
		ttyin:         ttyin,
		eventChan:     make(chan any, 64),
		oldState:      state,
	}, nil
}

func initRenderer(isAlternative bool) {
	csi("s")
	if isAlternative {
		smcup()
	}
	hideCursor()
	flush()
}

func (w *generalCellWriter) size() (int, int) {
	return w.width, w.height
}

func (w *generalCellWriter) matrix() [][]cell {
	return w.rows
}

func (w *generalCellWriter) close(isAlternative bool) error {
	showCursor()
	if isAlternative {
		rmcup()
		csi("u")
		flush()
		err := term.Restore(w.fd(), w.oldState)
		if err != nil {
			return err
		}
	} else {
		csi(fmt.Sprintf("%dB", w.height-w.cursorY-2))
		push("\n")
		flush()
		csi("1A")
		flush()
		err := term.Restore(w.fd(), w.oldState)
		push("\n")
		flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *generalCellWriter) updateTerminalSize() (bool, error) {
	width, height, err := term.GetSize(w.fd())
	if !w.isAlternative {
		height = w.height
	}
	if err != nil {
		return false, err
	}
	hasChanged := w.width != width || w.height != height
	if hasChanged {
		w.width = width
		w.height = height
		w.rows = make([][]cell, w.height)
		for y := 0; y < w.height; y++ {
			w.rows[y] = make([]cell, w.width)
		}
	}
	return hasChanged, nil
}

func (w *generalCellWriter) fill(s style) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.rows[y][x] = cell{' ', 1, s}
		}
	}
}

func (w *generalCellWriter) put(c cell, x, y int) {
	if w.rows[y][x] != c {
		w.rows[y][x] = c
	}
}

func (w *generalCellWriter) draw() {
	if w.isAlternative {
		origin()
	} else {
		csi(fmt.Sprintf("%dA", w.cursorY-1))
		push("\r")
	}
	lastStyle := style{}
	for y := 0; y < w.height; y++ {
		if y != 0 {
			csi("1B")
			push("\r")
		}
		for x := 0; x < w.width; x++ {
			s := w.rows[y][x].Style
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
				w.cursorY = y
				w.cursorX = x
			}
			if !(w.rows[y][x].Width == 0 && w.rows[y][x-1].Width == 2) {
				push(string(w.rows[y][x].Char))
			}
		}
	}
	push("\033[1;0m") // Reset Style
	if w.isAlternative {
		origin()
	} else {
		csi(fmt.Sprintf("%dA", w.height-1))
		push("\r")
	}
	csi(fmt.Sprintf("%dB", w.cursorY))
	csi(fmt.Sprintf("%dC", w.cursorX))
	flush()
}

func (w *generalCellWriter) fd() int {
	return int(w.ttyin.Fd())
}
