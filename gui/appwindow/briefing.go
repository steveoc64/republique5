package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"strings"
)

type BriefingPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newBriefingPanel(app *App) *BriefingPanel {
	h := &BriefingPanel{
		app:    app,
		Header: widget.NewLabel("Briefing for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel(app.Briefing),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
