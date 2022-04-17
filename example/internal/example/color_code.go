package example

import (
	"fmt"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

var style = tui.CellStyle{F256: 218, B256: 53}

type ColorCodeView struct {
	tui.DefaultView
	position int
}

func (m *ColorCodeView) HandleEvent(event interface{}) string {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.ArrowLeft:
			if m.position > 0 {
				m.position--
			}
		case key.ArrowRight:
			if m.position < 255 {
				m.position++
			}
		case key.ArrowUp:
			if m.position > 15 {
				m.position -= 16
			}
		case key.ArrowDown:
			if m.position < 240 {
				m.position += 16
			}
		case key.Esc:
			return (*MenuView).Options(nil).Title
		}
	}
	return ""
}

func (m *ColorCodeView) Body(bool, tui.Size) []tui.Text {
	var slice []tui.Text
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			seq := i*16 + j
			if seq == m.position {
				slice = append(slice, tui.Text{Str: " ", Style: style})
				slice = append(slice, tui.Text{Str: fmt.Sprintf("%3d", i*16+j), Style: tui.CellStyle{F256: seq, B256: 103}})
			} else {
				slice = append(slice, tui.Text{Str: fmt.Sprintf("%4d", i*16+j), Style: tui.CellStyle{F256: seq, B256: style.B256}})
			}
		}
		slice = append(slice, tui.Text{Str: "\n", Style: style})
	}
	return slice
}

func (*ColorCodeView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title: "256 Color Codes",
	}
}
