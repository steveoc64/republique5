package appwindow

import (
	"github.com/davecgh/go-spew/spew"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// MapWidget is a complete map viewer widget
// ... or will be when it grows up
type MapWidget struct {
	app      *App
	size     fyne.Size
	position fyne.Position
	hidden   bool
	mapData  *rp.MapData
	grid     *gridData
}

func newMapWidget(app *App, mapData *rp.MapData) *MapWidget {
	mw := &MapWidget{
		app:     app,
		mapData: mapData,
		grid:    newGridData(mapData.X, mapData.Y),
	}

	// generate the forces list in the grid
	for _, c := range app.Commands {
		mw.grid.addCommand(c)
		for _, u := range c.Units {
			mw.grid.addUnit(u)
		}
		for _, s := range c.Subcommands {
			mw.grid.addCommand(s)
			for _, u := range s.Units {
				mw.grid.addUnit(u)
			}
		}
	}

	// set size
	mw.Resize(mw.MinSize())
	return mw
}

// Size returns the current size of the mapWidget
func (mw *MapWidget) Size() fyne.Size {
	return mw.size
}

// Resize resizes the mapWidget
func (mw *MapWidget) Resize(size fyne.Size) {
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
	mw.hidden = false
	for _, obj := range widget.Renderer(mw).Objects() {
		obj.Show()
	}
}

// Hide sets the mapWidget to be not visible
func (mw *MapWidget) Hide() {
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

// Tapped is called when the user taps the map widget
func (mw *MapWidget) Tapped(event *fyne.PointEvent) {
	cmd := mw.grid.CommandAt(event.Position)
	if cmd != nil {
		if mw.grid.Select(cmd.Id) {
			widget.Renderer(mw).Refresh()
		}
	}
}

// TappedSecondary is called when the user right-taps the map widget
func (t *MapWidget) TappedSecondary(event *fyne.PointEvent) {
	spew.Dump(event, "tappedSecondary")
}
