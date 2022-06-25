package main

import (
	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(rootView)
	if err != nil {
		panic(err)
	}
}

func rootView() *tui.View {
	values := []int{0, 95, 135, 175, 215, 255}
	return tui.VMap(values, func(r int) *tui.View {
		return tui.VMap(values, func(g int) *tui.View {
			return tui.HMap(values, func(b int) *tui.View {
				return tui.P(
					tui.Span(" "),
					tui.Fmt("%3d, %3d, %3d", r, g, b).
						FGColor(tui.If(color.RelativeBrightness(r, g, b) > 0.5, color.RGB(0, 0, 0), color.RGB(255, 255, 255))).
						BGColor(color.RGB(r, g, b)),
				)
			}).AbsoluteSize(0, 1)
		})
	}).Border().Title("RGB").AbsoluteSize(89, 40)
}
