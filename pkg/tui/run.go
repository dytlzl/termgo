package tui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dytlzl/tervi/pkg/key"
)

func Run(createView func() *View, options ...option) error {

	cfg := config{}
	for _, opt := range options {
		err := opt(&cfg)
		if err != nil {
			return err
		}
	}

	r, err := NewRenderer()
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer func() {
		r.Close()
		for _, obj := range bufferForDebug {
			fmt.Printf("%#v\n", obj)
		}
	}()
	var shouldTerminate = false

	go func() {
		for event := range cfg.channel {
			if event == Terminate {
				shouldTerminate = true
			}
			r.eventChan <- event
		}
	}()

	keyChannel := make(chan rune, 1024)
	keyBuffer := make([]rune, 0)
	go func() {
		reader := bufio.NewReaderSize(os.Stdin, 256)
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
		v := createView()
		for {
			ch, size := readBuffer(keyBuffer)
			if size == 0 {
				break
			}
			keyBuffer = keyBuffer[size:]
			if ch == key.CtrlC {
				return nil
			}
			if cfg.eventHandler != nil {
				switch cfg.eventHandler(ch).(type) {
				case terminate:
					return nil
				}
			}
		}
		event := func() any {
			select {
			case event := <-r.eventChan:
				r.shouldSkipRendering = false
				return event
			case <-time.After(time.Millisecond * 10):
				return nil
			}
		}()

		if cfg.eventHandler != nil {
			switch cfg.eventHandler(event).(type) {
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
		r.fill(*cfg.style)

		v = createView()

		// Render views
		err = renderView(r, v, cfg, rect{0, 0, r.width, r.height})
		if err != nil {
			return fmt.Errorf("failed to render view: %w", err)
		}

		// Draw
		r.draw()

	}
}

type terminate struct{}

var Terminate = terminate{}

var bufferForDebug = make([]any, 0)

func renderView(r *Renderer, v *View, cfg config, frame rect) error {
	w, err := newWidget(r, frame.x, frame.y, frame.width, frame.height, v.paddingTop, v.paddingLeading, v.paddingBottom, v.paddingTrailing)
	if err != nil {
		return fmt.Errorf("failed to create widget: %w", err)
	}
	if v.style == nil {
		v.style = cfg.style
	}
	if v.border != nil || v.title != "" || v.renderer != nil {
		w.fill(cell{' ', 1, *v.style})
	}
	if v.border != nil {
		w.putBorder(*v.border)
	}
	if v.title != "" {
		w.putTitle([]Text{{Str: " " + v.title + " ", Style: *v.style}})
	}
	if v.renderer != nil {
		w.putBody(v.renderer(Size{Width: frame.width - v.paddingLeading - v.paddingTrailing, Height: frame.height - v.paddingTop - v.paddingBottom}))
	}

	availableWidth := frame.width - v.paddingLeading - v.paddingTrailing
	availableHeight := frame.height - v.paddingTop - v.paddingBottom

	x := frame.x + v.paddingLeading
	y := frame.y + v.paddingTop
	if v.reverseH {
		x += availableWidth
	}
	if v.reverseV {
		y += availableHeight
	}
	for _, child := range v.children {
		if child == nil {
			return nil
		}
		if child.absoluteWidth == 0 {
			if child.relativeWidth == 0 {
				child.relativeWidth = 12
			}
			child.absoluteWidth = availableWidth * child.relativeWidth / 12
			if child.x+child.relativeWidth == 12 {
				child.absoluteWidth = availableWidth - availableWidth*child.x/12
			}
		}

		if child.absoluteHeight == 0 {
			if child.relativeHeight == 0 {
				child.relativeHeight = 12
			}
			child.absoluteHeight = availableHeight * child.relativeHeight / 12
			if child.y+child.relativeHeight == 12 {
				child.absoluteHeight = availableHeight - availableHeight*child.y/12
			}
		}

		x := frame.x + v.paddingLeading + availableWidth*child.x/12
		if v.reverseH {
			x = frame.x + v.paddingLeading + availableWidth - availableWidth*child.x/12 - child.absoluteWidth
		}
		y := frame.y + v.paddingTop + availableHeight*child.y/12
		if v.reverseV {
			y = frame.y + v.paddingTop + availableHeight - availableHeight*child.y/12 - child.absoluteHeight
		}

		err = renderView(r, child, cfg, rect{
			x,
			y,
			child.absoluteWidth,
			child.absoluteHeight,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type rect struct {
	x      int
	y      int
	width  int
	height int
}
