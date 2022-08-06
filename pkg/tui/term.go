package tui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func TermSize() (width int, height int, err error) {
	return term.GetSize(int(os.Stdin.Fd()))
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
