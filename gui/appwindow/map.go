package appwindow

import (
	"context"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	rp "github.com/steveoc64/republique5/proto"
)

type MapPanel struct {
	app     *App
	content *fyne.Container
}

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

	useString := app.MapData.Data
	topSize := app.MapData.X
	switch app.MapData.Side {
	case rp.MapSide_FRONT:
		// nothing to do here
	case rp.MapSide_RIGHT_FLANK, rp.MapSide_LEFT_FLANK:
		topSize = app.MapData.Y
		runes := make([]byte, app.MapData.X*app.MapData.Y)
		i := 0
		for x := 0; x < int(app.MapData.X); x++ {
			for y := 0; y < int(app.MapData.Y); y++ {
				println(i, "->", x, y, (y+1)*int(app.MapData.X)-x-1)
				runes[i] = useString[(y+1)*int(app.MapData.X)-x-1]
				i++
			}
		}
		useString = string(runes)
		if app.MapData.Side == rp.MapSide_RIGHT_FLANK {
			useString = reverse(useString)
		}
	case rp.MapSide_TOP:
		useString = reverse(useString)
	}

	m := &MapPanel{
		app:     app,
		content: fyne.NewContainerWithLayout(layout.NewGridLayout(int(topSize))),
	}
	for _, v := range useString {
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
		r := canvas.NewRectangle(c)
		r.SetMinSize(fyne.Size{Width: 64, Height: 64})
		m.content.AddObject(r)
	}
	m.content.Show()
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
