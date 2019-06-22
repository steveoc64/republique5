package appwindow

import (
	"context"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// MapPanel is the UI for the map
type MapPanel struct {
	app       *App
	content   *fyne.Container
	mapWidget *MapWidget
	hbox      *widget.Box
	unitDesc  *widget.Label
}

// CanvasObject returns the top level UI element for the map
func (m *MapPanel) CanvasObject() fyne.CanvasObject {
	return m.content
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func newMapPanel(app *App) *MapPanel {
	if err := app.GetMap(); err != nil {
		println("Failed to get map", err.Error())
		return nil
	}
	m := &MapPanel{
		app: app,
		unitDesc: widget.NewLabelWithStyle(
			"No Unit Selected",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Italic: true},
		),
	}

	m.hbox = widget.NewHBox(
		layout.NewSpacer(),
		m.unitDesc,
		layout.NewSpacer(),
	)

	m.mapWidget = newMapWidget(app, app.MapData, m.unitDesc)
	m.content = fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, m.hbox, nil, nil),
		m.mapWidget,
		m.hbox,
	)

	return m
}

// GetMap fetches the map from the server
func (a *App) GetMap() error {
	mapData, err := a.gameServer.GetMap(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.MapData = mapData
	return nil
}
