package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type HeaderBar struct {
	Box *widget.Box
	app *App

	TeamName *widget.Label
	GameName *widget.Label
	GameTime *widget.Label
	Toolbar  *widget.Toolbar
}

func (h *HeaderBar) CanvasObject() fyne.CanvasObject {
	return h.Box
}

func newHeaderBar(app *App) *HeaderBar {
	h := &HeaderBar{
		app:      app,
		TeamName: widget.NewLabel(app.TeamName),
		GameName: widget.NewLabel(app.GameName),
		GameTime: widget.NewLabel(app.GameTime.Format(republique.DateFormat)),
		Toolbar: widget.NewToolbar(
			widget.NewToolbarAction(theme.ContentCopyIcon(), app.ToggleTheme),
		),
	}
	h.TeamName.TextStyle = fyne.TextStyle{Bold: true}
	h.GameTime.TextStyle = fyne.TextStyle{Italic: true}
	//h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(3),
	//h.TeamName, h.GameName, h.GameTime)
	h.Box = widget.NewVBox(
		widget.NewHBox(
			h.TeamName,
			layout.NewSpacer(),
			h.GameName,
			layout.NewSpacer(),
			h.GameTime,
		),
		h.Toolbar,
	)
	return h
}

func (h *HeaderBar) Recalc() *HeaderBar {
	h.GameName.SetText(h.app.GameName)
	h.GameTime.SetText(h.app.GameTime.Format(republique.TimeFormat))
	return h
}
