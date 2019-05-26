package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type HeaderBar struct {
	Box *widget.Box
	app *App

	TeamName *widget.Label
	GameName *widget.Label
	GameTime *widget.Label
}

func newHeaderBar(app *App) *HeaderBar {
	h := &HeaderBar{
		app:      app,
		TeamName: widget.NewLabel(app.TeamName),
		GameName: widget.NewLabel(app.GameName),
		GameTime: widget.NewLabel(app.GameTime.Format(republique.DateFormat)),
	}
	h.TeamName.TextStyle = fyne.TextStyle{Bold: true}
	h.GameTime.TextStyle = fyne.TextStyle{Italic: true}
	//h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(3),
	//h.TeamName, h.GameName, h.GameTime)
	h.Box = widget.NewHBox(
		h.TeamName,
		layout.NewSpacer(),
		h.GameName,
		layout.NewSpacer(),
		h.GameTime,
	)
	return h
}

func (h *HeaderBar) Recalc() *HeaderBar {
	h.GameName.SetText(h.app.GameName)
	h.GameTime.SetText(h.app.GameTime.Format(republique.TimeFormat))
	return h
}
