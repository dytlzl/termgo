package example

import (
	"reflect"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

type MenuView struct {
	tui.DefaultView
	Tabs             []tui.View
	currentTabNumber int
}

func (m MenuView) Title() string {
	return "Menu"
}

func (m MenuView) Body(bool, *tui.GlobalState) []tui.Text {
	var style = tui.CellStyle{F256: 218, B256: 53}
	var whiteStyle = tui.CellStyle{F256: 255, B256: 53}
	var slice []tui.Text
	for i, element := range m.Tabs {
		if i == m.currentTabNumber {
			slice = append(slice, tui.Text{Str: "> ", Style: style})
			slice = append(slice, tui.Text{Str: element.Title(), Style: whiteStyle})
		} else {
			slice = append(slice, tui.Text{Str: "  ", Style: style})
			slice = append(slice, tui.Text{Str: element.Title(), Style: style})
		}
		slice = append(slice, tui.Text{Str: "\n", Style: style})
	}
	return slice
}

func (m *MenuView) HandleEvent(event interface{}, state *tui.GlobalState) {
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
			state.FocusedModelType = reflect.TypeOf(m.Tabs[m.currentTabNumber])
		case key.Esc:
			state.ShouldTerminate = true
		}
	}
}

func (m *MenuView) SubViews() []tui.View {
	return []tui.View{&ColorCodeView{}}
}
