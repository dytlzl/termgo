package main

import (
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(rootView, tui.OptionEventHandler(handleEvent))
	if err != nil {
		panic(err)
	}
}

var position = 0
var input = ""

func rootView() *tui.View {
	r, size := utf8.DecodeRuneInString(input[position:])
	return tui.VStack(
		tui.If(position == len(input),
			tui.InlineStack(
				tui.String(input[:position]),
				tui.Cursor(" ").FGColor(-1),
			),
			tui.If(r == '\n',
				tui.InlineStack(
					tui.String(input[:position]),
					tui.Cursor(" ").FGColor(-1),
					tui.String(input[position:]),
				),
				tui.InlineStack(
					tui.String(input[:position]),
					tui.Cursor(input[position:position+size]).FGColor(-1),
					tui.String(input[position+size:]),
				),
			),
		),
	).Border().Title("Note").RelativeSize(9, 9)
}

func handleEvent(event any) any {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.Enter:
			input += "\n"
			position++
		case key.ArrowLeft:
			if position > 0 {
				_, size := utf8.DecodeLastRuneInString(input[:position])
				position -= size
			}
		case key.ArrowRight:
			if position < len(input) {
				_, size := utf8.DecodeRuneInString(input[position:])
				position += size
			}
		case key.ArrowUp, key.ArrowDown:
		case key.Del:
			if input != "" {
				_, size := utf8.DecodeLastRuneInString(input[:position])
				input = input[:position-size] + input[position:]
				position -= size
			}
		default:
			input = input[:position] + string(typed) + input[position:]
			position += utf8.RuneLen(typed)
		}
	}
	return nil
}
