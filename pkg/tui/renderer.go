package tui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type renderer struct {
	isAlternative       bool
	width               int
	height              int
	rows                [][]cell
	cursorX             int
	cursorY             int
	eventChan           chan any
	ttyin               *os.File
	oldState            *term.State
	shouldSkipRendering bool
}

func newRenderer(isAlternative bool) (*renderer, error) {
	ttyin := os.Stdin
	state, err := term.MakeRaw(int(ttyin.Fd()))
	if err != nil {
		return nil, err
	}
	height := 15
	if !isAlternative {
		push(strings.Repeat("\n", height-2))
		flush()
	}
	initRenderer(isAlternative)
	return &renderer{
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

func (r *renderer) close(isAlternative bool) error {
	showCursor()
	if isAlternative {
		rmcup()
		csi("u")
		flush()
		err := term.Restore(r.fd(), r.oldState)
		if err != nil {
			return err
		}
	} else {
		csi(fmt.Sprintf("%dB", r.height-r.cursorY-2))
		push("\n")
		flush()
		csi("1A")
		flush()
		err := term.Restore(r.fd(), r.oldState)
		push("\n")
		flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) fd() int {
	return int(r.ttyin.Fd())
}

var buffer string

func flush() {
	_, _ = fmt.Fprint(os.Stderr, buffer)
	buffer = ""
}

func push(s string) {
	buffer += s
}

func csi(s string) {
	buffer += "\x1b[" + s
}

func smcup() {
	csi("?1049h")
}

func rmcup() {
	csi("?1049l")
}

func showCursor() {
	csi("?25h")
}

func hideCursor() {
	csi("?25l")
}

func clearAll() {
	csi("2J")
}

func origin() {
	csi("1000A")
	push("\r")
}
