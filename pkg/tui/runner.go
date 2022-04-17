package tui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dytlzl/tervi/pkg/key"
)

type terminate struct{}

var Terminate = terminate{}

func Run(views map[string]View, options Options, channel chan interface{}) error {
	r, err := NewRenderer()
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer r.Close()

	focusedViewName := options.DefaultViewName

	var shouldTerminate = false

	go func() {
		for event := range channel {
			if event == Terminate {
				shouldTerminate = true
			}
			r.eventChan <- event
		}
	}()

	keyChannel := make(chan rune, 128)
	keyBuffer := make([]rune, 0)
	go func() {
		reader := bufio.NewReaderSize(os.Stdin, 10)
		for {
			ch, _, err := reader.ReadRune()
			if err != nil {
				panic(fmt.Errorf("failed to read keyboard input: %w", err))
			}
			keyChannel <- ch
		}
	}()

	for {
		r.shouldSkipRendering = false
	Depth1:
		for {
			select {
			case k := <-keyChannel:
				keyBuffer = append(keyBuffer, k)
			case <-time.After(time.Millisecond):
				break Depth1
			}
		}
		if len(keyBuffer) == 0 {
			r.shouldSkipRendering = true
		}
		for {
			ch, size := readBuffer(keyBuffer)
			if size == 0 {
				break
			}
			keyBuffer = keyBuffer[size:]
			if ch == key.CtrlC {
				return nil
			}
			switch value := views[focusedViewName].HandleEvent(ch).(type) {
			case string:
				if value != "" {
					focusedViewName = value
				}
			case terminate:
				return nil
			}
		}
		event := func() interface{} {
			select {
			case event := <-r.eventChan:
				r.shouldSkipRendering = false
				return event
			case <-time.After(time.Millisecond * 10):
				return nil
			}
		}()
		switch value := views[focusedViewName].HandleEvent(event).(type) {
		case string:
			if value != "" {
				focusedViewName = value
			}
		case terminate:
			return nil
		}
		if options.Footer != nil {
			switch value := options.Footer.HandleEvent(event).(type) {
			case string:
				if value != "" {
					focusedViewName = value
				}
			case terminate:
				return nil
			}
		}

		if changed, _ := r.UpdateTerminalSize(); changed {
			r.shouldSkipRendering = false
		}

		if shouldTerminate {
			return nil
		}

		if r.shouldSkipRendering {
			continue
		}

		// Clear
		r.Fill(options.Style)

		// Footer Widget
		footerHeight := 0
		if options.Footer != nil {
			footerHeight = 1
			footerWidget, err := newWidget(r, 0, r.height-footerHeight, r.width, footerHeight, 1, 0)
			if err != nil {
				return fmt.Errorf("failed to create footer widget: %w", err)
			}
			footerWidget.fill(Cell{Char: ' ', Width: 1, Style: options.Footer.Style()})
			footerWidget.putBody(options.Footer.Text())
		}

		size := Size{Width: r.width, Height: r.height}

		// Main Widget
		mainWidget, err := newWidget(r, 0, 0, r.width, r.height-footerHeight, 2, 2)
		if err != nil {
			return fmt.Errorf("failed to create main widget: %w", err)
		}
		viewOptions := views[focusedViewName].Options()
		if viewOptions.BorderStyle != nil {
			mainWidget.putBorder(*viewOptions.BorderStyle)
		} else {
			mainWidget.putBorder(options.Style)
		}
		if viewOptions.Title != "" {
			mainWidget.putTitle([]Text{{Str: " " + viewOptions.Title + " ", Style: options.Style}})
		}
		mainWidget.putBody(views[focusedViewName].Body(true, size))

		// Sub Widget
		for _, subView := range viewOptions.SubViews {
			viewOptions := subView.Options()
			if viewOptions.Width == nil {
				viewOptions.Width = NewFraction(2, 3)
			}
			w1, w2 := viewOptions.Width.Numer, viewOptions.Width.Denom
			subWidget, err := newWidget(r, r.width*(w2-w1)/w2, 2, r.width*w1/w2-2, r.height-footerHeight-4, 2, 2)
			if err != nil {
				return fmt.Errorf("failed to create sub widget: %w", err)
			}
			subWidget.fill(Cell{' ', 1, options.Style})
			if viewOptions.BorderStyle != nil {
				subWidget.putBorder(*viewOptions.BorderStyle)
			} else {
				subWidget.putBorder(options.Style)
			}
			if viewOptions.Title != "" {
				subWidget.putTitle([]Text{{Str: " " + viewOptions.Title + " ", Style: options.Style}})
			}
			subWidget.putBody(subView.Body(true, size))
		}

		// Draw
		r.Draw()
	}
}
