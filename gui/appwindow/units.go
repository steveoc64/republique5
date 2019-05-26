package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"strings"
)

type UnitsPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newUnitsPanel(app *App) *UnitsPanel {
	h := &UnitsPanel{
		app:    app,
		Header: widget.NewLabel("Units for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel("No units Yet ..."),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
