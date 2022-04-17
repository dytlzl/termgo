package main

import (
	"github.com/dytlzl/tervi/example/internal/example"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	menu := &example.MenuView{}
	models := []tui.View{
		menu,
		&example.ColorCodeView{},
		&example.DograMagraView{},
		&example.KeyCodeView{},
		&example.NoteView{},
	}
	menu.Tabs = models[1:]
	tui.Run(
		models,
		tui.Option{Style: tui.CellStyle{F256: 255, B256: 53}, Footer: &example.Footer{}}, nil)
}
