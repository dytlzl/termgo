package main

import (
	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/component"
	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	isQuitMenuOpen := false
	isConfirmedTermination := false
	position := 0
	input := ""
	err := tui.Run(func() *tui.View {
		return tui.ZStack(
			component.TextField(&input, &position).
				Border(tui.BorderOptionFGColor(color.RGB(100, 100, 100))).
				Title("Note").
				RelativeSize(9, 9),
			component.QuitView(&isQuitMenuOpen, &isConfirmedTermination).
				BGColor(color.RGB(145, 0, 145)),
		)
	}, tui.OptionEventHandler(func(event any) any {
		switch typed := event.(type) {
		case rune:
			switch typed {
			case key.Esc:
				isQuitMenuOpen = true
			}
		}
		return nil
	}))
	if err != nil {
		panic(err)
	}
}
