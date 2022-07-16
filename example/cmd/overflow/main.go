package main

import (
	"fmt"
	"strings"

	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(func() *tui.View {
		return tui.VMapN(50, func(i int) *tui.View {
			return tui.String(strings.Repeat(fmt.Sprintf("%03d;", i), 150)).AbsoluteSize(0, 5).Border().Padding(1, 2)
		}).Border().Padding(1, 2)
	})
	if err != nil {
		panic(err)
	}
}
