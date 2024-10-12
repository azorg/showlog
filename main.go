// File: "main.go"

package main

import (
	"log"

	"fyne.io/fyne/v2" // fyne
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	//"fyne.io/fyne/v2/widget"
)

func main() {
	//os.Setenv("FYNE_SCALE", "0.8") // FIXME
	a := app.New()

	// Read config
	conf := NewConf()
	conf.Read()

	// Create window
	w := a.NewWindow("Show log")
	w.SetFullScreen(conf.Full)
	if conf.Center {
		w.CenterOnScreen()
	}

	if conf.W > 0 && conf.H > 0 {
		w.Resize(fyne.NewSize(conf.W, conf.H))
	}

	w.SetOnClosed(func() {
		log.Printf("onCLose: full=%v", w.FullScreen())
	})

	// Create widgets
	//...

	// Show window and run
	//w.SetContent(vbox)
	w.ShowAndRun()

	// Save configuration
	//conf.Full = w.FullScreen()
	size := w.Canvas().Size()
	conf.W = size.Width
	conf.H = size.Height
	conf.Write()
}

// EOF: "main.go"
