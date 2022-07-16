package tui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dytlzl/tervi/pkg/key"
)

func Print(createView func() *View, options ...option) error {
	cfg := config{}
	for _, opt := range options {
		err := opt(&cfg)
		if err != nil {
			return err
		}
	}

	isAlternative := false

	r, err := newRenderer(isAlternative)
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer func() {
		r.close(isAlternative)
		for _, obj := range bufferForDebug {
			fmt.Printf("%#v\n", obj)
		}
	}()

	if changed, _ := r.updateTerminalSize(); changed {
		r.shouldSkipRendering = false
	}

	// Clear
	r.fill(style{})

	v := ZStack(createView())

	// Render views
	err = computeView(r, v, &cfg, rect{0, 0, r.width, r.height}, rect{0, 0, r.width, r.height}, style{})
	if err != nil {
		return fmt.Errorf("failed to render view: %w", err)
	}

	// Draw
	r.draw()

	return nil
}

func Run(createView func() *View, options ...option) error {
	cfg := config{}
	for _, opt := range options {
		err := opt(&cfg)
		if err != nil {
			return err
		}
	}

	isAlternative := true

	r, err := newRenderer(isAlternative)
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer func() {
		r.close(isAlternative)
		for _, line := range bufferForDebug {
			fmt.Println(line)
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
		for {
			ch, size := readBuffer(keyBuffer)
			if size == 0 {
				break
			}
			keyBuffer = keyBuffer[size:]
			if ch == key.CtrlC {
				return nil
			}
			if cfg.focusedView != nil {
				switch cfg.focusedView.keyHandler(ch).(type) {
				case terminate:
					return nil
				case nil:
					if cfg.eventHandler != nil {
						switch cfg.eventHandler(ch).(type) {
						case terminate:
							return nil
						}
					}
				}
			} else if cfg.eventHandler != nil {
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

		if changed, _ := r.updateTerminalSize(); changed {
			r.shouldSkipRendering = false
		}

		if shouldTerminate {
			return nil
		}

		if r.shouldSkipRendering {
			continue
		}

		// Clear
		r.fill(style{})

		v := ZStack(createView()).AbsoluteSize(r.width, r.height)
		cfg.focusedView = nil

		// Render views
		err = computeView(r, v, &cfg, rect{0, 0, r.width, r.height}, rect{0, 0, r.width, r.height}, style{})
		if err != nil {
			return fmt.Errorf("failed to render view: %w", err)
		}

		// Draw
		r.draw()

	}
}

type terminate struct{}

var Terminate = terminate{}

func computeView(r *renderer, v *View, cfg *config, frame rect, parentFrame rect, defaultStyle style) error {
	vr, err := newViewRenderer(
		r,
		frame,
		parentFrame,
		int(v.paddingTop), int(v.paddingLeading), int(v.paddingBottom), int(v.paddingTrailing),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to create viewRenderer: %w", err)
	}
	if v.style == nil {
		v.style = new(style)
	}
	v.style.merge(defaultStyle)
	if v.border != nil || v.title != "" || v.renderer != nil || v.style.b256 != 0 {
		vr.fill(cell{' ', 1, *v.style})
	}
	if v.border != nil {
		v.border.merge(*v.style)
		vr.putBorder(*v.border)
	}
	if v.title != "" {
		vr.putTitle([]text{{Str: " " + v.title + " ", Style: *v.style}})
	}
	if v.renderer != nil {
		vr.putBody(v.renderer(), *v.style)
	}
	if v.keyHandler != nil && (cfg.focusedView == nil || v.priority >= cfg.focusedView.priority) {
		cfg.focusedView = v
	}

	availableWidth := frame.width - int(v.paddingLeading) - int(v.paddingTrailing)
	availableHeight := frame.height - int(v.paddingTop) - int(v.paddingBottom)

	accumulatedX := frame.x + int(v.paddingLeading)
	accumulatedY := frame.y + int(v.paddingTop)

	remainedWidth := availableWidth
	remainedHeight := availableHeight
	numberOfAutoWidth := 0
	numberOfAutoHeight := 0

	children := func() []*View {
		if v.children == nil {
			return nil
		} else {
			return v.children()
		}
	}()

	for idx := range children {
		if children[idx] == nil {
			continue
		}
		// calculate absolute size from relative size
		if children[idx].absoluteWidth == 0 {
			children[idx].absoluteWidth = availableWidth * int(children[idx].relativeWidth) / 12
		}
		if children[idx].absoluteHeight == 0 {
			children[idx].absoluteHeight = availableHeight * int(children[idx].relativeHeight) / 12
		}

		remainedWidth -= children[idx].absoluteWidth
		remainedHeight -= children[idx].absoluteHeight

		// count auto-sizing view
		if children[idx].absoluteWidth == 0 {
			numberOfAutoWidth++
		}
		if children[idx].absoluteHeight == 0 {
			numberOfAutoHeight++
		}
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		// calculate size of auto-sizing view
		if child.absoluteWidth == 0 {
			if v.dir == horizontal {
				child.absoluteWidth = remainedWidth / numberOfAutoWidth
				numberOfAutoWidth--
				remainedWidth -= child.absoluteWidth
			} else {
				child.absoluteWidth = availableWidth
			}
		}
		if child.absoluteHeight == 0 {
			if v.dir == vertical {
				child.absoluteHeight = remainedHeight / numberOfAutoHeight
				numberOfAutoHeight--
				remainedHeight -= child.absoluteHeight
			} else {
				child.absoluteHeight = availableHeight
			}
		}

		x := frame.x + int(v.paddingLeading) + (availableWidth-child.absoluteWidth)/2
		if v.dir == horizontal {
			x = accumulatedX
		}
		y := frame.y + int(v.paddingTop) + (availableHeight-child.absoluteHeight)/2
		if v.dir == vertical {
			y = accumulatedY
		}

		if y+int(v.paddingBottom) >= frame.y+frame.height {
			break
		}

		err = computeView(r, child, cfg,
			rect{
				x,
				y,
				child.absoluteWidth,
				child.absoluteHeight,
			},
			rect{
				frame.x + int(v.paddingLeading),
				frame.y + int(v.paddingTop),
				v.absoluteWidth - int(v.paddingLeading) - int(v.paddingTrailing),
				v.absoluteHeight - int(v.paddingTop) - int(v.paddingBottom),
			},
			*v.style)
		if err != nil {
			return err
		}
		if v.dir == horizontal {
			accumulatedX += child.absoluteWidth
		}
		if v.dir == vertical {
			accumulatedY += child.absoluteHeight
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

var bufferForDebug = make([]string, 0)

func debugf(format string, a ...any) {
	bufferForDebug = append(bufferForDebug, fmt.Sprintf(format, a...))
}
