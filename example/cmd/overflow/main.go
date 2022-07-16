package main

import (
	"fmt"
	"strings"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(renderBlocks)
	if err != nil {
		panic(err)
	}
}

var cursor = 0

func renderBlocks() *tui.View {

	return tui.VMapN(5, func(i int) *tui.View {
		return tui.String(strings.Repeat(fmt.Sprintf("%03d;", i), 150)).AbsoluteSize(0, 5).Border().Padding(1, 2)
	}).Border().Padding(1, 2).OffsetY(cursor).KeyHandler(func(r rune) any {
		switch r {
		case key.ArrowUp:
			cursor--
		case key.ArrowDown:
			cursor++
		}
		return true
	})
}

func renderInlines() *tui.View {
	cursor := 0
	return tui.InlineMapN(50, func(i int) *tui.View {
		return tui.String(strings.Repeat(fmt.Sprintf("%03d;", i), 150)).AbsoluteSize(0, 5).Border().Padding(1, 2)
	}).Border().Padding(1, 2).OffsetY(cursor).KeyHandler(func(r rune) any {
		switch r {
		case key.ArrowUp:
			cursor--
		case key.ArrowDown:
			cursor++
		}
		return true
	})
}
