package tui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Renderer struct {
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

func NewRenderer() (*Renderer, error) {
	ttyin := os.Stdin
	state, err := term.MakeRaw(int(ttyin.Fd()))
	if err != nil {
		return nil, err
	}
	initRenderer()
	return &Renderer{
		ttyin:     ttyin,
		eventChan: make(chan any, 64),
		oldState:  state,
	}, nil
}

func initRenderer() {
	csi("s")
	smcup()
	hideCursor()
	flush()
}

func (r *Renderer) Close() error {
	showCursor()
	rmcup()
	csi("u")
	flush()
	err := term.Restore(r.fd(), r.oldState)
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) fd() int {
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

func origin() {
	csi("1000A")
	push("\r")
}
