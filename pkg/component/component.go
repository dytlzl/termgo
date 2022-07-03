package component

import (
	"strings"
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func TextField(input *string, position *int) *tui.View {
	i := *input
	p := *position
	r, size := utf8.DecodeRuneInString(i[p:])
	return tui.VStack(
		tui.If(p == len(i),
			tui.InlineStack(
				tui.String(i[:p]),
				tui.Cursor(" ").Invert(true),
			),
			tui.If(r == '\n',
				tui.InlineStack(
					tui.String(i[:p]),
					tui.Cursor(" ").Invert(true),
					tui.String(i[p:]),
				),
				tui.InlineStack(
					tui.String(i[:p]),
					tui.Cursor(i[p:p+size]).Invert(true),
					tui.String(i[p+size:]),
				),
			),
		),
	).KeyHandler(func(r rune) any {
		switch r {
		case key.Esc:
			return nil
		case key.Enter:
			*input += "\n"
			*position++
		case key.ArrowLeft:
			if *position > 0 {
				_, size := utf8.DecodeLastRuneInString((*input)[:*position])
				*position -= size
			}
		case key.ArrowRight:
			if *position < len(*input) {
				_, size := utf8.DecodeRuneInString((*input)[*position:])
				*position += size
			}
		case key.ArrowUp, key.ArrowDown:
		case key.Del:
			if *input != "" {
				_, size := utf8.DecodeLastRuneInString((*input)[:*position])
				*input = (*input)[:*position-size] + (*input)[*position:]
				*position -= size
			}
		default:
			*input = (*input)[:*position] + string(r) + (*input)[*position:]
			*position += utf8.RuneLen(r)
		}
		return true
	})
}

func QuitView(isOpen, isConfirmed *bool) *tui.View {
	return tui.InlineStack(
		tui.Fmt("%sAre you sure to quit?\n\n%s     ",
			strings.Repeat(" ", (32-21)/2),
			strings.Repeat(" ", (32-21)/2)),
		tui.String(" Yes ").Invert(*isConfirmed),
		tui.String(" "),
		tui.String(" No ").Invert(!*isConfirmed),
	).
		AbsoluteSize(36, 7).
		Title("Quit").
		Border().
		KeyHandler(func(r rune) any {
			switch r {
			case key.Esc:
				*isConfirmed = false
				*isOpen = false
			case key.ArrowLeft:
				*isConfirmed = true
			case key.ArrowRight:
				*isConfirmed = false
			case key.Enter:
				if *isConfirmed {
					return tui.Terminate
				} else {
					*isOpen = false
				}
			}
			return true
		}).
		Hidden(!*isOpen)
}
