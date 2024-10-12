// File: "main.go"

package main

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2" // fyne
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	//os.Setenv("FYNE_SCALE", "0.8") // FIXME
	a := app.New()

	// Read config
	conf := NewConf()
	conf.Read()

	// Create main window
	w := a.NewWindow("Show log")
	w.SetFullScreen(conf.Full)
	if conf.Center {
		w.CenterOnScreen()
	}
	if conf.W > 0 && conf.H > 0 {
		w.Resize(fyne.NewSize(conf.W, conf.H))
	}
	w.SetOnClosed(func() {
		log.Printf("window closed")
	})

	// Create widgets
	lw := NewLogWidget("")
	
	btnQuit := widget.NewButton("Quit", func() {
		log.Print("button Quit pressed")
		a.Quit()
		os.Exit(0)
	})

	btnOpen := widget.NewButton("Open", func() {
		msg := "button Open pressed"
		log.Print(msg)
		text := lw.Text() + time.Now().Format(time.DateTime) + " " + msg + "\n"
		lw.SetText(text)
	})

	btnClear := widget.NewButton("Clear", func() {
		log.Print("button Clear pressed")
		lw.SetText("")
	})

	// Create containers/spacers/layout
	spacer := layout.NewSpacer
	vbox := container.NewBorder(
		nil, // top
		container.NewHBox(btnOpen, btnClear, spacer(), btnQuit), // bottom
		nil, nil, // left, right
		lw, // center
	)


	// Show window and run
	w.SetContent(vbox)
	w.ShowAndRun()

	// Save configuration
	size := w.Canvas().Size()
	conf.W = size.Width
	conf.H = size.Height
	conf.Write()
}

// EOF: "main.go"
