package example

import (
	"fmt"

	"github.com/dytlzl/tervi/pkg/tui"
)

type KeyCodeView struct {
	codes []rune
}

func (k *KeyCodeView) Body(bool, tui.Size) []tui.Text {
	slice := make([]tui.Text, 0, len(k.codes))
	for i := len(k.codes) - 1; i >= 0 && i > len(k.codes)-5000; i-- {
		slice = append(slice, tui.Text{Str: fmt.Sprintf(" %d", int(k.codes[i])), Style: style})
	}
	return slice
}

func (k *KeyCodeView) HandleEvent(event interface{}) string {
	if k.codes == nil {
		k.codes = make([]rune, 0)
	}
	switch typed := event.(type) {
	case rune:
		k.codes = append(k.codes, typed)
	}
	return ""
}

func (*KeyCodeView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title: "Key Code",
	}
}
