package appwindow

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

type gridForces struct {
	commands []*rp.Command
	units    []*rp.Unit
}

type gridData struct {
	x, y  int32
	back  []color.RGBA
	value []byte
	units []gridForces
}

func newGridData(x, y int32) *gridData {
	g := &gridData{
		x:     x,
		y:     y,
		back:  make([]color.RGBA, x*y),
		value: make([]byte, x*y),
		units: make([]gridForces, x*y),
	}
	for i := 0; i < int(x*y); i++ {
		g.back[i] = color.RGBA{uint8(rand.Intn(40) + 160), uint8(rand.Intn(40) + 180), uint8(rand.Intn(40) + 100), 200}
	}
	return g
}

func (g *gridData) Color(x, y int32) color.RGBA {
	i := y*g.x + x
	if i < 0 || i > int32(len(g.back))-1 {
		return color.RGBA{}
	}
	return g.back[i]
}

func (g *gridData) Value(x, y int32) byte {
	i := y*g.x + x
	if i < 0 || i > int32(len(g.value))-1 {
		return ' '
	}
	return g.value[i]
}

func (g *gridData) Units(x, y int32) gridForces {
	i := y*g.x + x
	if i < 0 || i > int32(len(g.units))-1 {
		return gridForces{}
	}
	return g.units[i]
}

func (g *gridData) addCommand(c *rp.Command) {
	x := c.GetGameState().GetGrid().GetX() - 1
	y := c.GetGameState().GetGrid().GetY() - 1
	i := y*g.x + x
	if i < 0 || i > int32(len(g.units))-1 {
		return
	}
	println("addCommand", i, c, len(g.units))
	g.units[i].commands = append(g.units[i].commands, c)
}

func (g *gridData) addUnit(c *rp.Unit) {
	x := c.GetGameState().GetGrid().GetX() - 1
	y := c.GetGameState().GetGrid().GetY() - 1
	i := y*g.x + x
	if i < 0 || i > int32(len(g.units))-1 {
		return
	}
	g.units[i].units = append(g.units[i].units, c)
}

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
	rand.Seed(time.Now().UnixNano())
	return newMapRender(mw)
}
