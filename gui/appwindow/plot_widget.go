package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

type timedata struct {
	name   string
	values []int
}

func (t timedata) Len() int {
	return len(t.values)
}

func (t timedata) XY(i int) (x, y float64) {
	return float64(i), float64(t.values[i])
}

// PlotWidget plots a set of data
type PlotWidget struct {
	app      *App
	size     fyne.Size
	position fyne.Position
	hidden   bool
	data     []timedata
	title    string
}

func newPlotWidget(app *App, title string) *PlotWidget {
	p := &PlotWidget{
		app:   app,
		title: title,
	}
	return p
}

// Size returns the current size of the plotWidget
func (p *PlotWidget) Size() fyne.Size {
	return p.size
}

// Resize resizes the plotWidget
func (p *PlotWidget) Resize(size fyne.Size) {
	p.size = size
	widget.Renderer(p).Layout(p.size)
	canvas.Refresh(p)
}

// Position returns the current position of the plotWidget
func (p *PlotWidget) Position() fyne.Position {
	return p.position
}

// Move orders the plotWidget to be moved
func (p *PlotWidget) Move(pos fyne.Position) {
	p.position = pos
	widget.Renderer(p).Layout(p.size)
}

// MinSize returns the minSize of the plotWidget
func (p *PlotWidget) MinSize() fyne.Size {
	return widget.Renderer(p).MinSize()
}

// Visible returns whether the plotWidget is visible or not
func (p *PlotWidget) Visible() bool {
	return !p.hidden
}

// Show sets the plotWidget to be visible
func (p *PlotWidget) Show() {
	p.hidden = false
	for _, obj := range widget.Renderer(p).Objects() {
		obj.Show()
	}
}

// Hide sets the plotWidget to be not visible
func (p *PlotWidget) Hide() {
	p.hidden = true
	for _, obj := range widget.Renderer(p).Objects() {
		obj.Hide()
	}
}

// ApplyTheme applies the theme to the plotWidget
func (p *PlotWidget) ApplyTheme() {
	widget.Renderer(p).ApplyTheme()
}

// CreateRenderer builds a new renderer
func (p *PlotWidget) CreateRenderer() fyne.WidgetRenderer {
	return newPlotRender(p)
}

// Clear removes all plot data
func (p *PlotWidget) Clear() {
	p.data = []timedata{}
}

// AddSet adds a set of data to the plot
func (p *PlotWidget) AddSet(data timedata) {
	p.data = append(p.data, data)
	widget.Renderer(p).Refresh()
}
