package main

import (
	"fmt"

	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(func() *tui.View {
		return tui.VStack(tui.CreateView(Body)).Border(tui.DefaultStyle).Title("RGB").AbsoluteSize(89, 40)
	})
	if err != nil {
		panic(err)
	}
}

var style = tui.Style{F256: 218, B256: 53}

func Body(tui.Size) []tui.Text {
	var slice []tui.Text
	values := []int{0, 95, 135, 175, 215, 255}
	for _, r := range values {
		for _, g := range values {
			for _, b := range values {
				slice = append(slice, tui.Text{Str: " ", Style: tui.DefaultStyle})
				fg := color.RGB(255, 255, 255)
				if color.RelativeBrightness(r, g, b) > 0.5 {
					fg = color.RGB(0, 0, 0)
				}
				slice = append(slice, tui.Text{Str: fmt.Sprintf("%3d, %3d, %3d", r, g, b), Style: tui.Style{F256: fg, B256: color.RGB(r, g, b)}})
			}
			slice = append(slice, tui.Text{Str: "\n", Style: tui.DefaultStyle})
		}
	}
	return slice
}
