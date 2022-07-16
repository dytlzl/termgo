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

func renderBlocks() *tui.View {
	cursor, setCursor := tui.UseState(0)
	return tui.VMapN(8, func(i int) *tui.View {
		return tui.String(strings.Repeat(fmt.Sprintf("%03d;", i), 150)).AbsoluteSize(0, 5).Border().Padding(1, 2)
	}).Border().Padding(1, 2).OffsetY(cursor).KeyHandler(func(r rune) any {
		switch r {
		case key.ArrowUp:
			setCursor(cursor+1)
		case key.ArrowDown:
			setCursor(cursor-1)
		}
		return true
	}).AllowOverflow()
}

