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
	m          *Widget
	objects    []fyne.CanvasObject
	background *image.RGBA
	dirty      bool
	bgcolors   []color.RGBA
	x          int
	y          int
}

func newMapEditorRender(m *Widget) *mapEditorRender {
	r := &mapEditorRender{
		m: m,
	}
	render := canvas.NewRaster(r.getImage)
	r.render = render
	r.objects = []fyne.CanvasObject{render}
	r.generateBGColors()
	return r
}

func (r *mapEditorRender) generateBGColors() {
	println("generating new color set", r.m.x, r.m.y)
	r.bgcolors = []color.RGBA{}
	for i := 0; i < r.m.x*r.m.y; i++ {
		c := color.RGBA{uint8(rand.Intn(40) + 160), uint8(rand.Intn(40) + 180), uint8(rand.Intn(40) + 100), 200}
		r.bgcolors = append(r.bgcolors, c)
	}
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
	if r.x != r.m.x || r.y != r.m.y {
		r.generateBGColors()
	}
	if r.dirty || r.background == nil || r.background.Bounds().Size().X != w || r.background.Bounds().Size().Y != h {
		r.background = r.generateBackground(w, h)
	}
	r.x = r.m.x
	r.y = r.m.y
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
		i := x + (y * x)
		c := r.bgcolors[i]
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
					gc.SetFillColor(mapWoodsFill)
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
					gc.SetFillColor(mapWoodsFill)
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
					gc.SetFillColor(mapTownFill)
					gc.SetStrokeColor(mapTownStroke)
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
					gc.SetFillColor(mapTownFill)
					gc.SetStrokeColor(mapTownStroke)
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

	// do the river segments
	gc.SetFillColor(mapDeepBlue)
	gc.SetStrokeColor(mapBlue)
	gc.SetLineWidth(20)
	for k, v := range r.m.rivers {
		gc.BeginPath()
		fx := (float64(k.x) + 0.5) * dx
		fy := (float64(k.y) + 0.5) * dy
		if len(v.adjacent) == 0 {
			draw2dkit.Ellipse(gc, fx, fy, dx/3, dy/5)
			gc.FillStroke()
			continue
		}
		for kk := range v.adjacent {
			// double check that this segment isnt already done
			if other, ok := r.m.rivers[riverPoint{kk.x, kk.y}]; ok {
				if toMe, okk := other.adjacent[riverPoint{k.x, k.y}]; okk {
					if !toMe {
						gc.MoveTo(fx, fy)
						fx2 := (float64(kk.x) + 0.5) * dx
						fy2 := (float64(kk.y) + 0.5) * dy
						gc.LineTo(fx2, fy2)
						v.adjacent[kk] = true
					}
				}
			}
		}
		gc.Stroke()

		// if we are on an edge, then complete the river
		switch {
		case k.x == 0, k.x == r.m.x-1:
			fx2 := 0.0
			if k.x == r.m.x-1 {
				fx2 = float64(r.m.x) * dx
			}
			bump := 0.0
			for kk := range v.adjacent {
				switch {
				case kk.y < k.y:
					bump = dy / 2
				case kk.y > k.y:
					bump = dy / -2
				}
				gc.MoveTo(fx, fy)
				gc.LineTo(fx2, fy+bump)
				gc.Stroke()
				break
			}
		case k.y == 0, k.y == r.m.y-1:
			fy2 := 0.0
			if k.y == r.m.y-1 {
				fy2 = float64(r.m.y) * dy
			}
			bump := 0.0
			for kk := range v.adjacent {
				switch {
				case kk.x < k.x:
					bump = dx / 2
				case kk.x > k.x:
					bump = dx / -2
				}
				gc.MoveTo(fx, fy)
				gc.LineTo(fx+bump, fy2)
				gc.Stroke()
				break
			}
		}
	}

	// major grid lines - vertical
	gc.SetStrokeColor(mapGrid)
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

	// highlight current grid
	gc.SetStrokeColor(mapSelect)
	gc.SetLineWidth(8)
	x1 := float64(r.m.cx-1) * dx
	y1 := float64(r.m.cy-1) * dy
	draw2dkit.Rectangle(gc, x1, y1, x1+dx, y1+dy)
	gc.Stroke()

	return img
}

func (r *mapEditorRender) hills(gc *draw2dimg.GraphicContext, fx, fy, dx, dy float64, count int) {
	gc.SetFillColor(mapHillFill)
	gc.SetStrokeColor(mapHillStroke)
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

func (r *mapEditorRender) ConvertToGrid(event *fyne.PointEvent) (int, int) {
	if r.background == nil {
		return 0, 0
	}
	scale := r.Scale()
	size := r.background.Bounds()
	dx := (float32(size.Max.X) / scale) / float32(r.m.x)
	dy := (float32(size.Max.Y) / scale) / float32(r.m.y)
	x := int(float32(event.Position.X) / dx)
	y := int(float32(event.Position.Y) / dy)
	return x + 1, y + 1
}
