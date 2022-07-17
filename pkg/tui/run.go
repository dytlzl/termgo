package tui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dytlzl/tervi/pkg/key"
)

func Run(createView func() *View, options ...option) error {
	cfg := config{
		viewPQ: newQueue(),
	}
	for _, opt := range options {
		err := opt(&cfg)
		if err != nil {
			return err
		}
	}

	isAlternative := true

	w, err := newGeneralCellWriter(isAlternative)
	if err != nil {
		return fmt.Errorf("failed to init renderer: %w", err)
	}
	defer func() {
		w.close(isAlternative)
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
			w.eventChan <- event
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
		shouldSkipRendering := false
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
			shouldSkipRendering = true
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
			handled := false
			pq := newQueue()
			for _, v := range cfg.viewPQ {
				pq.PushView(v)
			}
		Depth2:
			for pq.Len() > 0 {
				v := pq.PopView()
				switch v.keyHandler(ch).(type) {
				case terminate:
					return nil
				case nil:
				default:
					handled = true
					break Depth2
				}
			}
			if cfg.eventHandler != nil && !handled {
				switch cfg.eventHandler(ch).(type) {
				case terminate:
					return nil
				}
			}
		}
		event := func() any {
			select {
			case event := <-w.eventChan:
				shouldSkipRendering = false
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

		if changed, _ := w.updateTerminalSize(); changed {
			shouldSkipRendering = false
		}

		if shouldTerminate {
			return nil
		}

		if shouldSkipRendering {
			continue
		}

		// Clear
		w.fill(style{})

		v := ZStack(createView()).AbsoluteSize(w.width, w.height)

		// Render views
		cfg.viewPQ = newQueue()
		err = moldView(w, v, &cfg, rect{0, 0, w.width, w.height}, rect{0, 0, w.width, w.height}, style{}, false)
		if err != nil {
			return fmt.Errorf("failed to render view: %w", err)
		}

		// Draw
		w.draw()

	}
}

type terminate struct{}

var Terminate = terminate{}

var bufferForDebug = make([]string, 0)

func debugf(format string, a ...any) {
	bufferForDebug = append(bufferForDebug, fmt.Sprintf(format, a...))
}
