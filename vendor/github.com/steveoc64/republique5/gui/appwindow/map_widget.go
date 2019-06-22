package appwindow

import (
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
	unitDesc *widget.Label
}

func newMapWidget(app *App, mapData *rp.MapData, unitDesc *widget.Label) *MapWidget {
	mw := &MapWidget{
		app:      app,
		mapData:  mapData,
		grid:     newGridData(mapData.X, mapData.Y, mapData.Data),
		unitDesc: unitDesc,
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
	mw.Select(mw.grid.CommandAt(event.Position).GetId())
}

// TappedSecondary is called when the user right-taps the map widget
func (mw *MapWidget) TappedSecondary(event *fyne.PointEvent) {
	if cmd := mw.grid.CommandAt(event.Position); cmd != nil {
		mw.app.unitsPanel.ShowCommand(cmd)
		mw.app.Tab(TAB_UNITS)
	}
}

// Select selects the command with the given ID
func (mw *MapWidget) Select(id int32) (*rp.Command, bool) {
	if cmd, ok := mw.grid.Select(id); ok {
		widget.Renderer(mw).Refresh()
		mw.app.mapPanel.SetCommand(cmd)
		return cmd, ok
	}
	mw.app.mapPanel.SetCommand(nil)
	return nil, false
}