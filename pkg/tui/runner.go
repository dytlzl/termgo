package tui

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/dytlzl/tervi/pkg/key"
)

type Terminate struct{}

func Run(views []View, option Option, channel chan interface{}) error {
	r, err := NewRenderer()
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer r.Close()

	viewMap := map[reflect.Type]View{}
	for _, model := range views {
		viewMap[reflect.TypeOf(model)] = model
	}

	var state = &GlobalState{}
	state.FocusedModelType = reflect.TypeOf(views[0])

	go func() {
		for event := range channel {
			if _, ok := event.(Terminate); ok {
				state.ShouldTerminate = true
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
		state.ShouldSkipRendering = false
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
			state.ShouldSkipRendering = true
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
			viewMap[state.FocusedModelType].HandleEvent(ch, state)
			if state.ShouldTerminate {
				return nil
			}
		}
		select {
		case event := <-r.eventChan:
			viewMap[state.FocusedModelType].HandleEvent(event, state)
			if option.Footer != nil {
				option.Footer.HandleEvent(event, state)
			}
			if state.ShouldTerminate {
				return nil
			}
			state.ShouldSkipRendering = false
		case <-time.After(time.Millisecond * 10):
			viewMap[state.FocusedModelType].HandleEvent(nil, state)
			if option.Footer != nil {
				option.Footer.HandleEvent(nil, state)
			}
		}
		if changed, _ := r.UpdateTerminalSize(); changed {
			state.Width = r.width
			state.Height = r.height
			state.ShouldSkipRendering = false
		}
		if state.ShouldSkipRendering {
			continue
		}

		// Clear
		r.Fill(option.Style)

		// Footer Widget
		footerHeight := 0
		if option.Footer != nil {
			footerHeight = 1
			footerWidget, err := newWidget(r, 0, r.height-footerHeight, r.width, footerHeight, 1, 0)
			if err != nil {
				return fmt.Errorf("failed to create footer widget: %w", err)
			}
			footerWidget.fill(Cell{Char: ' ', Width: 1, Style: option.Footer.Style()})
			footerWidget.putBody(option.Footer.Text())
		}

		// Main Widget
		mainWidget, err := newWidget(r, 0, 0, r.width, r.height-footerHeight, 2, 2)
		if err != nil {
			return fmt.Errorf("failed to create main widget: %w", err)
		}
		if borderStyle := viewMap[state.FocusedModelType].BorderStyle(); borderStyle != nil {
			mainWidget.putBorder(*borderStyle)
		} else {
			mainWidget.putBorder(option.Style)
		}
		mainWidget.putTitle([]Text{{Str: " " + viewMap[state.FocusedModelType].Title() + " ", Style: option.Style}})
		mainWidget.putBody(viewMap[state.FocusedModelType].Body(true, state))

		// Sub Widget
		for _, subView := range viewMap[state.FocusedModelType].SubViews() {
			w1, w2 := subView.Width()
			subWidget, err := newWidget(r, r.width*(w2-w1)/w2, 2, r.width*w1/w2-2, r.height-footerHeight-4, 2, 2)
			if err != nil {
				return fmt.Errorf("failed to create sub widget: %w", err)
			}
			subWidget.fill(Cell{' ', 1, option.Style})
			if borderStyle := subView.BorderStyle(); borderStyle != nil {
				subWidget.putBorder(*borderStyle)
			} else {
				subWidget.putBorder(option.Style)
			}
			subWidget.putTitle([]Text{{Str: " " + subView.Title() + " ", Style: option.Style}})
			subWidget.putBody(subView.Body(true, state))
		}

		// Draw
		r.Draw()
	}
}
