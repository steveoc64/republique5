package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
	"image/color"
)

type gridData struct {
	x,y int
	back []color.RGBA
	value []byte
}

func newGridData(x,y int) *gridData {
	return &gridData{
		x: x,
		y: y,
		back: make([]color.RGBA, x*y*4),
		value: make([]byte, x*y),
	}
}

func (g *gridData) Color(x,y int) color.RGBA {
	i := y*g.x + x
	if i < 0 || y > len(g.back)-1 {
		return nil
	}
	return g.back[i]
}

func (g *gridData) Value(x,y int) byte {
	i := y*g.x + x
	if i < 0 || y > len(g.value)-1 {
		return nil
	}
	return g.value[i]
}

// MapWidget is a complete map viewer widget
// ... or will be when it grows up
type MapWidget struct {
	app *App
	size     fyne.Size
	position fyne.Position
	hidden   bool
	mapData  *rp.MapData

	grid *gridData
}

func newMapWidget(app *App, mapData *rp.MapData) *MapWidget {
	mw := &MapWidget{
		app:     app,
		mapData: mapData,
		grid: newGridData(mapData.X, mapData.Y)
	}
	mw.Resize(mw.MinSize())
	return mw
}

// Size returns the current size of the mapWidget
func (mw *MapWidget) Size() fyne.Size {
	return mw.size
}

// Resize resizes the mapWidget
func (mw *MapWidget) Resize(size fyne.Size) {
	println("mw resize to", size.Width, size.Height)
	mw.size = size
	widget.Renderer(mw).Layout(mw.size)
	canvas.Refresh(mw)
}

// Position returns the current position of the mapWidget
func (mw *MapWidget) Position() fyne.Position {
	return mw.position
}

// Move orders the mapWidget to be moved
func (mw *MapWidget) Move(pos fyne.Position) {
	mw.position = pos
	widget.Renderer(mw).Layout(mw.size)
}

// MinSize returns the minSize of the mapWitdget
func (mw *MapWidget) MinSize() fyne.Size {
	return widget.Renderer(mw).MinSize()
}

// Visible returns whether the mapWidget is visible or not
func (mw *MapWidget) Visible() bool {
	return !mw.hidden
}

// Show sets the mapWidget to be visible
func (mw *MapWidget) Show() {
	println("widget show")
	mw.hidden = false
	for _, obj := range widget.Renderer(mw).Objects() {
		obj.Show()
	}
}

// Hide sets the mapWidget to be not visible
func (mw *MapWidget) Hide() {
	println("widget hide")
	mw.hidden = true
	for _, obj := range widget.Renderer(mw).Objects() {
		obj.Hide()
	}
}

// ApplyTheme applies the theme to the mapWidget
func (mw *MapWidget) ApplyTheme() {
	widget.Renderer(mw).ApplyTheme()
}

// CreateRenderer builds a new renderer
func (mw *MapWidget) CreateRenderer() fyne.WidgetRenderer {
	return newMapRender(mw)
}
