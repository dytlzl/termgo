package tui

type Style struct {
	F256          uint8
	B256          uint8
	invert        bool
	bold          bool
	italic        bool
	strikethrough bool
	underline     bool
	HasCursor     bool
}
