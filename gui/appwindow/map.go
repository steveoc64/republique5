package appwindow

import (
	"context"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

// MapPanel is the UI for the map
type MapPanel struct {
	app       *App
	content   *fyne.Container
	mapWidget *MapWidget
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
	}

	m.mapWidget = newMapWidget(app, app.MapData)
	m.mapWidget.Hide()
	m.content = fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	m.content.AddObject(m.mapWidget)
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
