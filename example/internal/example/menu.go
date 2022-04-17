package example

import (
	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

type MenuView struct {
	Tabs             []tui.View
	Channel          chan interface{}
	currentTabNumber int
}

func (m *MenuView) Body(bool, tui.Size) []tui.Text {
	var style = tui.CellStyle{F256: 218, B256: 53}
	var whiteStyle = tui.CellStyle{F256: 255, B256: 53}
	var slice []tui.Text
	for i, element := range m.Tabs {
		if i == m.currentTabNumber {
			slice = append(slice, tui.Text{Str: "> ", Style: style})
			slice = append(slice, tui.Text{Str: element.Options().Title, Style: whiteStyle})
		} else {
			slice = append(slice, tui.Text{Str: "  ", Style: style})
			slice = append(slice, tui.Text{Str: element.Options().Title, Style: style})
		}
		slice = append(slice, tui.Text{Str: "\n", Style: style})
	}
	return slice
}

func (m *MenuView) HandleEvent(event interface{}) string {
	switch typed := event.(type) {
	case rune:
		switch typed {
		case key.ArrowUp:
			if m.currentTabNumber > 0 {
				m.currentTabNumber--
			}
		case key.ArrowDown:
			if m.currentTabNumber < len(m.Tabs)-1 {
				m.currentTabNumber++
			}
		case key.Enter:
			return m.Tabs[m.currentTabNumber].Options().Title
		case key.Esc:
			m.Channel <- tui.Terminate
		}
	}
	return ""
}

func (m *MenuView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title:    "Menu",
		SubViews: []tui.View{&ColorCodeView{}},
	}
}
