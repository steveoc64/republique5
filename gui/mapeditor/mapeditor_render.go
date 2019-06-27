package mapeditor

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type mapEditorRender struct {
	render     *canvas.Raster
	m          *MapEditorWidget
	objects    []fyne.CanvasObject
	background *image.RGBA
	dirty      bool
}

func newMapEditorRender(m *MapEditorWidget) *mapEditorRender {
	r := &mapEditorRender{
		m: m,
	}
	render := canvas.NewRaster(r.getImage)
	r.render = render
	r.objects = []fyne.CanvasObject{render}
	return r
}

func (r *mapEditorRender) Scale() float32 {
	return fyne.CurrentApp().Driver().CanvasForObject(r.render).Scale()
}

// ApplyTheme applies the theme
func (r *mapEditorRender) ApplyTheme() {
	// noop
}

// BackgroundColor returns the background color for our map
func (r *mapEditorRender) BackgroundColor() color.Color {
	return color.RGBA{183, 211, 123, 1}
}

// Destroy removes any resources we have on this renderer
func (r *mapEditorRender) Destroy() {
	// noop
}

// Layout does .. the layout ?
func (r *mapEditorRender) Layout(size fyne.Size) {
	r.render.Resize(size)
}

// MinSize returns the minimum size for this renderer
func (r *mapEditorRender) MinSize() fyne.Size {
	return fyne.Size{
		Width:  int(r.m.x * 64),
		Height: int(r.m.y * 64),
	}
}

// Objects returns the slice of objects that we own
func (r *mapEditorRender) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh paints the map
func (r *mapEditorRender) Refresh() {
	canvas.Refresh(r.render)
}

func (r *mapEditorRender) getImage(w, h int) image.Image {
	if r.dirty || r.background == nil || r.background.Bounds().Size().X != w || r.background.Bounds().Size().Y != h {
		r.background = r.generateBackground(w, h)
	}
	return r.background
}

func (r *mapEditorRender) generateBackground(w, h int) *image.RGBA {
	r.dirty = false
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	dx := float64(w / r.m.x)
	dy := float64(h / r.m.y)
	mx := r.m.x
	my := r.m.y
	gc := draw2dimg.NewGraphicContext(img)
	twopi := math.Pi * 2

	paintBlock := func(x, y int) {
		//i := x + mx*y
		c := color.RGBA{uint8(rand.Intn(40) + 160), uint8(rand.Intn(40) + 180), uint8(rand.Intn(40) + 100), 200}
		draw.Draw(img,
			image.Rectangle{
				image.Point{x * int(dx), y * int(dy)},
				image.Point{(x+1)*int(dx) - 1, (y+1)*int(dy) - 1},
			},
			&image.Uniform{c},
			image.Point{0, 0},
			draw.Src)
	}

	// paint the background blocks
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			paintBlock(x, y)
		}
	}

	// grid overlays, hills and trees
	i := 0
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			fx := float64(x) * dx
			fy := float64(y) * dy
			mapChar := r.m.data[i]
			switch mapChar {
			case 'h':
				// small hills
				r.hills(gc, fx, fy, dx, dy, 1)
			case 'H':
				// larger hills
				r.hills(gc, fx, fy, dx, dy, 2)
			case 'w':
				for t := 0; t < 10+rand.Intn(20); t++ {
					gc.SetFillColor(map_woods_fill)
					gc.SetLineWidth(2)
					cx := fx + dx*float64(rand.Intn(100))/100.0
					cy := fy + dy*float64(rand.Intn(100))/100.0
					rr := 2 + rand.Intn(10)
					rx := dx * float64(rr) / 100.0
					ry := dy * float64(rr) / 100.0
					gc.BeginPath()
					gc.ArcTo(cx, cy, rx, ry, 0.0, twopi)
					gc.Close()
					gc.Fill()
				}
			case 'W':
				for t := 0; t < 25+rand.Intn(40); t++ {
					gc.SetFillColor(map_woods_fill)
					gc.SetLineWidth(2)
					cx := fx + dx*float64(rand.Intn(100))/100.0
					cy := fy + dy*float64(rand.Intn(100))/100.0
					rr := 2 + rand.Intn(10)
					rx := dx * float64(rr) / 100.0
					ry := dy * float64(rr) / 100.0
					gc.BeginPath()
					gc.ArcTo(cx, cy, rx, ry, 0.0, twopi)
					gc.Close()
					gc.Fill()
				}
			case 'T':
				for t := 0; t < 25+rand.Intn(40); t++ {
					gc.SetFillColor(map_town_fill)
					gc.SetStrokeColor(map_town_stroke)
					gc.SetLineWidth(2)
					x1 := fx + dx*float64(rand.Intn(100))/100.0
					y1 := fy + dy*float64(rand.Intn(100))/100.0
					x2 := x1 + dx*float64(2+rand.Intn(10))/100.0
					y2 := y1 + dy*float64(2+rand.Intn(10))/100.0
					gc.BeginPath()
					draw2dkit.Rectangle(gc, x1, y1, x2, y2)
					gc.Close()
					gc.Fill()
					gc.FillStroke()
				}
			case 't':
				for t := 0; t < 10+rand.Intn(20); t++ {
					gc.SetFillColor(map_town_fill)
					gc.SetStrokeColor(map_town_stroke)
					gc.SetLineWidth(2)
					x1 := fx + dx*float64(rand.Intn(100))/100.0
					y1 := fy + dy*float64(rand.Intn(100))/100.0
					x2 := x1 + dx*float64(2+rand.Intn(10))/100.0
					y2 := y1 + dy*float64(2+rand.Intn(10))/100.0
					gc.BeginPath()
					draw2dkit.Rectangle(gc, x1, y1, x2, y2)
					gc.Close()
					gc.Fill()
					gc.FillStroke()
				}
			}
			i++
		}
	}

	// draw rivers
	i = 0
	gc.SetFillColor(map_deep_blue)
	gc.SetStrokeColor(map_blue)
	gc.SetLineWidth(20)
	gc.BeginPath()
	gotriver := false
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			fx := float64(x) * dx
			fy := float64(y) * dy
			switch r.m.data[i] {
			case 'r':
				if !gotriver {
					gotriver = true
					gc.MoveTo(fx+0.5*dx, fy+0.5*dy)
				} else {
					gc.LineTo(fx+0.5*dx, fy+0.5*dy)
				}
			}
			i++
		}
	}
	if gotriver {
		gc.FillStroke()
	}

	// major grid lines - vertical
	gc.SetStrokeColor(map_grid)
	gc.SetLineWidth(2)
	for x := 0; x < mx; x++ {
		gc.BeginPath()
		gc.MoveTo(float64(x)*dx, 0.0)
		gc.LineTo(float64(x)*dx, float64(h))
		gc.Close()
		gc.FillStroke()
	}
	// major grid lines - horizontal
	for y := 0; y < my; y++ {
		gc.BeginPath()
		gc.MoveTo(0.0, float64(y)*dy)
		gc.LineTo(float64(w), float64(y)*dy)
		gc.Close()
		gc.FillStroke()
	}

	return img
}

func (r *mapEditorRender) hills(gc *draw2dimg.GraphicContext, fx, fy, dx, dy float64, count int) {
	gc.SetFillColor(map_hill_fill)
	gc.SetStrokeColor(map_hill_stroke)
	gc.SetLineWidth(2)
	switch count {
	case 1:
		gc.BeginPath()
		gc.MoveTo(fx+0.1*dx, fy+0.5*dy)
		gc.QuadCurveTo(fx+0.3*dx, fy+0.2*dy, fx+0.7*dx, fy+0.5*dy)
		//gc.Close()
		gc.FillStroke()
		gc.MoveTo(fx+0.3*dx, fy+0.7*dy)
		gc.QuadCurveTo(fx+0.5*dx, fy+0.4*dy, fx+0.9*dx, fy+0.7*dy)
		//gc.Close()
		gc.FillStroke()
	case 2:
		gc.BeginPath()
		gc.MoveTo(fx+0.4*dx, fy+0.3*dy)
		gc.QuadCurveTo(fx+0.6*dx, fy+0.1*dy, fx+0.8*dx, fy+0.3*dy)
		//gc.Close()
		gc.FillStroke()
		gc.BeginPath()
		gc.MoveTo(fx+0.1*dx, fy+0.5*dy)
		gc.QuadCurveTo(fx+0.3*dx, fy+0.2*dy, fx+0.7*dx, fy+0.5*dy)
		//gc.Close()
		gc.FillStroke()
		gc.MoveTo(fx+0.3*dx, fy+0.7*dy)
		gc.QuadCurveTo(fx+0.5*dx, fy+0.4*dy, fx+0.9*dx, fy+0.7*dy)
		//gc.Close()
		gc.FillStroke()
	}
}

func (r *mapEditorRender) ConvertToGrid(event *fyne.PointEvent) (int32, int32) {
	if r.background == nil {
		return 0, 0
	}
	scale := r.Scale()
	size := r.background.Bounds()
	dx := (float32(size.Max.X) / scale) / float32(r.m.x)
	dy := (float32(size.Max.Y) / scale) / float32(r.m.y)
	x := int32(float32(event.Position.X) / dx)
	y := int32(float32(event.Position.Y) / dy)
	return x + 1, y + 1
}
