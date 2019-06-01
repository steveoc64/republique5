package appwindow

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"strings"
)

type BriefingPanel struct {
	app *App
	Box *widget.Box

	Header *widget.Label
	Notes  *widget.Label
}

func newBriefingPanel(app *App) *BriefingPanel {
	h := &BriefingPanel{
		app:    app,
		Header: widget.NewLabel("Briefing for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel(app.Briefing),
	}
	img := canvas.NewImageFromResource(resourceBannerJpg)
	img.FillMode = canvas.ImageFillOriginal
	h.Box = widget.NewVBox(
		img,
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
