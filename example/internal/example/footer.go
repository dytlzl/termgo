package example

import "github.com/dytlzl/tervi/pkg/tui"

type Footer struct {
}

func (f Footer) Text() []tui.Text {
	return []tui.Text{
		{Str: "[Key Binding] Ctrl+C: Quit", Style: f.Style()},
	}
}

func (Footer) Style() tui.CellStyle {
	return tui.CellStyle{F256: 255, B256: 171}
}

func (Footer) HandleEvent(interface{}) string {
	return ""
}
