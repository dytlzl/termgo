package main

import (
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(rootView, tui.OptionEventHandler(handleEvent))
	if err != nil {
		panic(err)
	}
}

var codes = make([]rune, 0)

func rootView() *tui.View {
	return tui.InlineMap(codes[:], func(r rune) *tui.View {
		return tui.Fmt(" %d", int(r))
	}).Border().RelativeSize(9, 9)
}

func handleEvent(event any) any {
	switch typed := event.(type) {
	case rune:
		codes = append([]rune{typed}, codes...)
	}
	if len(codes) > 150 {
		codes = codes[:150]
	}
	return nil
}
