package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type AdvancePanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func (a *AdvancePanel) CanvasObject() fyne.CanvasObject {
	return a.Box
}

func newAdvancePanel(app *App) *AdvancePanel {
	h := &AdvancePanel{
		app:    app,
		Header: widget.NewLabelWithStyle("General Advance", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes:  widget.NewLabel("you sure ?"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
