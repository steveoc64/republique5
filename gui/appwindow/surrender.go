package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type SurrenderPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

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
