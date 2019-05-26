package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"strings"
)

type ActionsPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newActionsPanel(app *App) *ActionsPanel {
	h := &ActionsPanel{
		app:    app,
		Header: widget.NewLabel("Actions for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel("No Actions Yet ..."),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
