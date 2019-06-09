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
	/*
		for k, v := range useString {
			c := color.RGBA{183, 211, 123, 1}
			switch v {
			case 't': // town small
				c = color.RGBA{122, 110, 69, 1}
			case 'T': // town big
				c = color.RGBA{83, 73, 50, 1}
			case 'w': // woods light
				c = color.RGBA{56, 138, 60, 1}
			case 'W': // woods thick
				c = color.RGBA{37, 77, 39, 1}
			case 'r': // river
				c = color.RGBA{58, 114, 157, 1}
			case 'h': // hill low
				c = color.RGBA{176, 148, 78, 1}
			case 'H': // hill high
				c = color.RGBA{104, 92, 61, 1}
			}
			r := NewTapRect(c, func() {
				println("tapped on rect", k)
			}, nil)
			r.rect.SetMinSize(fyne.Size{Width: 64, Height: 64})
			m.content.AddObject(r)
			m.rects[k] = r
		}

	*/
	m.mapWidget = newMapWidget(app, app.MapData)
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
