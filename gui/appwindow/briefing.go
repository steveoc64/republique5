package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"strings"
)

// BriefingPanel in the UI for the briefing
type BriefingPanel struct {
	app    *App
	Box    *widget.Box
	Scroll *widget.ScrollContainer

	Header *widget.Label
	Notes  *widget.Label
}

// CanvasObject returns the top level UI element
func (b *BriefingPanel) CanvasObject() fyne.CanvasObject {
	return b.Box
}

// newBriefingPanel creates a new BriefingPanel and returns it
func newBriefingPanel(app *App) *BriefingPanel {
	briefingText := strings.Replace(app.Briefing, "\n", "\n\n", -1)
	h := &BriefingPanel{
		app: app,
		Header: widget.NewLabelWithStyle(
			"Briefing for: "+strings.Join(app.Commanders, ", "),
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes: widget.NewLabel(briefingText),
	}
	img := canvas.NewImageFromResource(resourceBannerJpg)
	img.FillMode = canvas.ImageFillOriginal
	h.Box = widget.NewVBox(
		img,
		h.Header,
		h.Notes,
	)
	h.Scroll = widget.NewScrollContainer(h.Box)
	h.Scroll.Resize(app.MinSize())
	h.Box.Show()
	return h
}
