package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type FormationsPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newFormationsPanel(app *App) *FormationsPanel {
	h := &FormationsPanel{
		app:    app,
		Header: widget.NewLabelWithStyle("Formations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes:  widget.NewLabel("Info about formations here"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
