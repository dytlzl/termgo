package example

import (
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

type NoteView struct {
	position int
	input    string
}

func (n *NoteView) Body(hasFocus bool, _ tui.Size) []tui.Text {
	style := tui.CellStyle{F256: 255, B256: 53}
	cursorStyle := tui.CellStyle{F256: style.B256, B256: style.F256, HasCursor: true}
	if hasFocus {
		if n.position == len(n.input) {
			return []tui.Text{
				{Str: n.input[:n.position], Style: style},
				{Str: " ", Style: cursorStyle},
			}
		}
		r, size := utf8.DecodeRuneInString(n.input[n.position:])
		if r == '\n' {
			return []tui.Text{
				{Str: n.input[:n.position], Style: style},
				{Str: " ", Style: cursorStyle},
				{Str: n.input[n.position:], Style: tui.CellStyle{F256: 255, B256: 53}},
			}
		}
		return []tui.Text{
			{Str: n.input[:n.position], Style: style},
			{Str: n.input[n.position : n.position+size], Style: cursorStyle},
			{Str: n.input[n.position+size:], Style: style},
		}
	} else {
		return []tui.Text{
			{Str: n.input, Style: style},
		}
	}
}

func (n *NoteView) HandleEvent(event interface{}) interface{} {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.Enter:
			n.input += "\n"
			n.position++
		case key.ArrowLeft:
			if n.position > 0 {
				_, size := utf8.DecodeLastRuneInString(n.input[:n.position])
				n.position -= size
			}
		case key.ArrowRight:
			if n.position < len(n.input) {
				_, size := utf8.DecodeRuneInString(n.input[n.position:])
				n.position += size
			}
		case key.Del:
			if n.input != "" {
				_, size := utf8.DecodeLastRuneInString(n.input[:n.position])
				n.input = n.input[:n.position-size] + n.input[n.position:]
				n.position -= size
			}
		case key.Esc:
			return (*MenuView).Options(nil).Title
		default:
			n.input = n.input[:n.position] + string(typed) + n.input[n.position:]
			n.position += utf8.RuneLen(typed)
		}
	}
	return nil
}

func (*NoteView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title: "Note",
	}
}
