// File: "main.go"

package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2" // fyne
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
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

	onExit := func() {
		// Save configuration
		size := w.Canvas().Size()
		conf.W = size.Width
		conf.H = size.Height
		conf.Write()
	}

	// Create widgets
	label := widget.NewLabel("")
	label.TextStyle.Bold = true
	lw := NewLogWidget("")

	if conf.File != "" {
		// Open file
		u := storage.NewFileURI(conf.File)
		log.Printf("open: scheme=%s path=%s", u.Scheme(), u.Path())
		label.SetText(u.Name())
		lw.Open(u.Path())
	}

	// Create buttons
	btnQuit := widget.NewButton("Quit", func() {
		log.Print("button Quit pressed")
		//a.Quit()
		onExit()
		os.Exit(0)
	})

	btnOpen := widget.NewButton("Open", func() {
		msg := "button Open pressed"
		log.Print(msg)
		d := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				log.Print("can't select file:", err)
				return
			}
			u := uri.URI()
			log.Printf("open: scheme=%s path=%s", u.Scheme(), u.Path())
			label.SetText(u.Name())
			conf.File = u.Path()
			lw.Open(conf.File)
			lw.TextGrid.SetText("")
		}, w)

		uri := storage.NewFileURI(".")
		uris, err := storage.ListerForURI(uri)
		if err != nil {
			log.Print("can't set start location:", err)
		} else {
			d.SetLocation(uris)
		}
		d.Show()
	})

	btnClear := widget.NewButton("Clear", func() {
		log.Print("button Clear pressed")
		lw.SetText("")
	})

	// Create containers/spacers/layout
	spacer := layout.NewSpacer
	vbox := container.NewBorder(
		container.NewCenter(label),                              // top
		container.NewHBox(btnOpen, btnClear, spacer(), btnQuit), // bottom
		nil, nil, // left, right
		lw.Scroll, // center
	)

	// Show window and run
	w.SetContent(vbox)
	w.ShowAndRun()
	onExit()
}

// EOF: "main.go"
