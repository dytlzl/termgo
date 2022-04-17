package main

import (
	"github.com/dytlzl/tervi/example/internal/example"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	menu := &example.MenuView{
		Channel: make(chan interface{}),
	}
	menu.Tabs = []tui.View{
		&example.ColorCodeView{},
		&example.DograMagraView{},
		&example.KeyCodeView{},
		&example.NoteView{},
	}
	views := map[string]tui.View{
		menu.Options().Title: menu,
	}
	for _, view := range menu.Tabs {
		views[view.Options().Title] = view
	}
	tui.Run(
		views,
		tui.Options{
			DefaultViewName: menu.Options().Title,
			Style:           tui.CellStyle{F256: 255, B256: 53},
			Footer:          &example.Footer{},
		},
		menu.Channel,
	)
}
