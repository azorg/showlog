// File: "logwidget.go"

package main
	
import (
	"fyne.io/fyne/v2/widget"
)

type LogWidget struct {
	*widget.TextGrid
}

func NewLogWidget(file string) *LogWidget {
	tg := widget.NewTextGrid()
	tg.SetText("hello!")
	return &LogWidget{TextGrid: tg}
}

// EOF: "logwidget.go"
