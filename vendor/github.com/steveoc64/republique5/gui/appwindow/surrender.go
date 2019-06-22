package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// SurrenderPanel is the UI for editting surrender terms
type SurrenderPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

// CanvasObject returns the top level UI element for the SurrenderPanel
func (s *SurrenderPanel) CanvasObject() fyne.CanvasObject {
	return s.Box
}

func newSurrenderPanel(app *App) *SurrenderPanel {
	h := &SurrenderPanel{
		app:    app,
		Header: widget.NewLabelWithStyle("Surrender", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes:  widget.NewLabel("Really surrender ?"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
