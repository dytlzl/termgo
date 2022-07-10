package main

import (
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Print(rootView)
	if err != nil {
		panic(err)
	}
}

func rootView() *tui.View {
	return tui.InlineMapN(16, func(i int) *tui.View {
		return tui.InlineStack(
			tui.InlineMapN(16, func(j int) *tui.View {
				seq := i*16 + j
				return tui.Fmt("%4d", seq).FGColor(uint8(seq))
			}),
			tui.Break(),
		)
	}).AbsoluteSize(69, 20).Border()
}
