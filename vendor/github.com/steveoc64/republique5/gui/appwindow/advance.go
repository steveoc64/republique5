package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// AdvancePanel is the UI for ordering a general advance
type AdvancePanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

// CanvasObject returns the top level UI element for the AdvancePanel
func (a *AdvancePanel) CanvasObject() fyne.CanvasObject {
	return a.Box
}

// newAdvancePanel builds and returns a new AdvancePanel
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
