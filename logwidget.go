// File: "logwidget.go"

package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2" // fyne
	//"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/storage"
)

const PERIOD = time.Duration(500) * time.Millisecond

// Log widget
type LogWidget struct {
	*widget.Entry
	//*container.Scroll
	*time.Ticker
	reader fyne.URIReadCloser
	mx     sync.Mutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

// Create new log widgets
func NewLogWidget(file string) *LogWidget {
	//entry := widget.NewEntry()
	entry := widget.NewMultiLineEntry()
	//entry.MultiLine = true
	entry.TextStyle.Monospace = true

	ctx, cancel := context.WithCancel(context.Background())
	lw := &LogWidget{
		Entry:  entry,
		Ticker: time.NewTicker(PERIOD),
		ctx:    ctx,
		cancel: cancel,
	}

	lw.wg.Add(1)
	go lw.goMonitor()
	return lw
}

// Open new log file
func (lw *LogWidget) Open(u fyne.URI) {
	lw.mx.Lock()
	defer lw.mx.Unlock()

	if lw.reader != nil {
		// Close old file
		err := lw.reader.Close()
		if err != nil {
			log.Print("close error: ", err)
		}
		lw.reader = nil
	}

	// Open new file
	var err error
	lw.reader, err = storage.Reader(u)
	if err != nil {
		lw.reader = nil
		log.Print("can't open file: ", err)
		lw.SetText(err.Error())
		return
	}

	lw.SetText("")
	//lw.Scroll.ScrollToTop()
	lw.update()
}

// Update entry from file
func (lw *LogWidget) update() {
	data, err := io.ReadAll(lw.reader)

	if err != nil {
		log.Print("read error: ", err)
		return
	}

	if len(data) == 0 {
		return
	}
	log.Printf("read %d bytes", len(data))

	// Append text
	text := lw.Text + string(data)
	lw.SetText(text)
	
	//focused := lw.Entry.Canvas().Focused()
	//// If the user is not focused on the text area then scroll to the end
	//if focused == nil || focused != lw.Entry {
  //  lw.Entry.CursorRow = len(lw.Text) - 1 // ets the cursor to the end
	//}
}

// Cancel monitor
func (lw *LogWidget) Cancel() {
	lw.cancel()
	lw.Ticker.Stop()

	lw.mx.Lock()
	defer lw.mx.Unlock()

	if lw.reader != nil {
		// CLose file
		err := lw.reader.Close()
		if err != nil {
			log.Print("close error: ", err)
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
		if lw.reader != nil {
			lw.update()
		}
		lw.mx.Unlock()
	} // for
}

// EOF: "logwidget.go"
