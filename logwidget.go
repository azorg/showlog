// File: "logwidget.go"

package main

import (
	"context"
	//"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const PERIOD = time.Duration(500) * time.Millisecond

// Log widget
type LogWidget struct {
	*widget.Entry
	*container.Scroll
	*os.File
	*time.Ticker
	mx     sync.Mutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

// Create new log widgets
func NewLogWidget(file string) *LogWidget {
	e := widget.NewEntry()
	e.TextStyle.Monospace = true
	//e.MultiLine = true

	ctx, cancel := context.WithCancel(context.Background())
	lw := &LogWidget{
		Entry:  e,
		Scroll: container.NewScroll(e),
		Ticker: time.NewTicker(PERIOD),
		ctx:    ctx,
		cancel: cancel,
	}

	lw.wg.Add(1)
	go lw.goMonitor()
	return lw
}

// Open new log file
func (lw *LogWidget) Open(path string) {
	lw.mx.Lock()
	defer lw.mx.Unlock()

	if lw.File != nil {
		// Close old file
		err := lw.File.Close()
		if err != nil {
			log.Print("close error:", err)
		}
		lw.File = nil
	}

	// Open new file
	var err error
	lw.File, err = os.Open(path)
	if err != nil {
		log.Print("can't open file:", err)
		lw.SetText(err.Error())
		return
	}

	lw.SetText("")
	lw.Scroll.ScrollToTop()
	lw.update()
}

// Update entry from file
func (lw *LogWidget) update() {
	data, err := io.ReadAll(lw.File)

	if err != nil {
		log.Print("read error:", err)
		return
	}

	if len(data) == 0 {
		return
	}
	log.Printf("read %d bytes", len(data))

	// Append text
	text := lw.Text + string(data)
	lw.SetText(text)
}

// Cancel monitor
func (lw *LogWidget) Cancel() {
	lw.cancel()
	lw.Ticker.Stop()

	lw.mx.Lock()
	defer lw.mx.Unlock()

	if lw.File != nil {
		// CLose file
		err := lw.File.Close()
		if err != nil {
			log.Print("close error:", err)
		}
	}

	lw.wg.Wait()
}

func (lw *LogWidget) goMonitor() {
	defer lw.wg.Done()
	for {
		select {
		case <-lw.ctx.Done():
			log.Print("cancel")
			return

		case _, ok := <-lw.Ticker.C:
			if !ok {
				log.Print("ticker closed")
				return
			}
		} // select

		lw.mx.Lock()
		if lw.File == nil {
			lw.mx.Unlock()
			continue
		}

		lw.update()
		lw.mx.Unlock()
	} // for
}

// EOF: "logwidget.go"
