package main

import (
	"fmt"
	"strings"

	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	err := tui.Run(func() *tui.View {
		return tui.VMapN(50, func(i int) *tui.View {
			return tui.String(strings.Repeat(fmt.Sprintf("%3d", i), 50)).AbsoluteSize(0, 1)
		}).Border()
	})
	if err != nil {
		panic(err)
	}
}
