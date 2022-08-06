package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(renderList)
	if err != nil {
		panic(err)
	}
	if name != "" {
		fmt.Println(name)
	}
}

var name = ""

func renderList() *tui.View {
	selected := tui.UseRef(0)
	files, _ := ioutil.ReadDir(".")
	return tui.ZStack(tui.ListMap(selected, files, func(file fs.FileInfo) *tui.View {
		return tui.HStack(tui.String(file.Name()).AbsoluteSize(20, 1), tui.Fmt("%d", file.Size()))
	}).Border().Padding(1, 2)).KeyHandler(func(r rune) any {
		switch r {
		case key.Enter:
			name = files[*selected].Name()
			return tui.Terminate
		default:
			return nil
		}
	})
}

func renderBlocks() *tui.View {
	cursor, setCursor := tui.UseState(0)
	return tui.VMapN(8, func(i int) *tui.View {
		return tui.String(strings.Repeat(fmt.Sprintf("%03d;", i), 150)).AbsoluteSize(0, 5).Border().Padding(1, 2)
	}).Border().Padding(1, 2).OffsetY(cursor).KeyHandler(func(r rune) any {
		switch r {
		case key.ArrowUp:
			setCursor(cursor + 1)
		case key.ArrowDown:
			setCursor(cursor - 1)
		}
		return true
	}).AllowOverflow()
}
