package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type MapPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newMapPanel(app *App) *MapPanel {
	h := &MapPanel{
		app:    app,
		Header: widget.NewLabelWithStyle("Map", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes:  widget.NewLabel("draw a map here"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
