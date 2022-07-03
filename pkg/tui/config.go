package tui

type config struct {
	channel      chan any
	style        Style
	focusedView  *View
	eventHandler func(any) any
}

func OptionChannel(ch chan any) func(*config) error {
	return func(c *config) error {
		c.channel = ch
		return nil
	}
}

func OptionStyle(style Style) func(*config) error {
	return func(c *config) error {
		c.style = style
		return nil
	}
}

func OptionEventHandler(fn func(any) any) func(*config) error {
	return func(c *config) error {
		c.eventHandler = fn
		return nil
	}
}

type option = func(*config) error
