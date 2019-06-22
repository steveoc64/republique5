package appwindow

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	republique "github.com/steveoc64/republique5/proto"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/llgcode/draw2d/draw2dimg"
)

type mapRender struct {
	render  *canvas.Raster
	mw      *MapWidget
	objects []fyne.CanvasObject
	img     *image.RGBA
	dirty   bool
}

func newMapRender(mw *MapWidget) *mapRender {
	r := &mapRender{mw: mw}
	render := canvas.NewRaster(r.getImage)
	//render := canvas.NewRasterWithPixels(r.paint2)
	r.render = render
	r.objects = []fyne.CanvasObject{render}
	//r.ApplyTheme()
	return r
}

func (r *mapRender) Scale() float32 {
	return fyne.CurrentApp().Driver().CanvasForObject(r.render).Scale()
}

// ApplyTheme applies the theme
func (r *mapRender) ApplyTheme() {
	// noop
}

// BackgroundColor returns the background color for our map
func (r *mapRender) BackgroundColor() color.Color {
	return color.RGBA{183, 211, 123, 1}
}

// Destroy removes any resources we have on this renderer
func (r *mapRender) Destroy() {
	// noop
}

// Layout does .. the layout ?
func (r *mapRender) Layout(size fyne.Size) {
	r.render.Resize(size)
}

// MinSize returns the minimum size for this renderer
func (r *mapRender) MinSize() fyne.Size {
	return fyne.Size{
		Width:  int(r.mw.mapData.X * 64),
		Height: int(r.mw.mapData.Y * 64),
	}
}

// Objects returns the slice of objects that we own
func (r *mapRender) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh paints the map
func (r *mapRender) Refresh() {
	r.dirty = true
	canvas.Refresh(r.render)
}

func (r *mapRender) getImage(w, h int) image.Image {
	if r.dirty || r.img == nil || r.img.Bounds().Size().X != w || r.img.Bounds().Size().Y != h {
		r.img = r.generateImage(w, h)
		r.dirty = false
	}
	if r.mw.hidden {
		return &image.RGBA{}
	}
	return r.img
}

func (r *mapRender) generateImage(w, h int) *image.RGBA {
	scale := float64(r.Scale())
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	dx := float64(w / int(r.mw.mapData.X))
	dy := float64(h / int(r.mw.mapData.Y))
	mx := int(r.mw.mapData.X)
	my := int(r.mw.mapData.Y)
	gc := draw2dimg.NewGraphicContext(img)
	twopi := math.Pi * 2

	paintBlock := func(x, y int) {
		i := x + mx*y
		if i < 0 || i > (len(r.mw.grid.back)-1) {
			println("error", i, mx, x, y)
			return
		}
		c := r.mw.grid.back[i]
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
			mapChar := r.mw.mapData.Data[i]
			switch mapChar {
			case 'h':
				// small hills
				r.hills(gc, fx, fy, dx, dy, 1)
			case 'H':
				// larger hills
				r.hills(gc, fx, fy, dx, dy, 2)
			case 'w':
				// draw a little circle for woods
				gc.SetFillColor(map_woods_fill)
				gc.SetStrokeColor(map_woods_stroke)
				gc.SetLineWidth(1)
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.1, dy*.1,
					0, twopi)
				gc.Close()
				gc.FillStroke()
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.15, dy*.15,
					0, twopi)
				gc.Close()
				gc.FillStroke()
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.1, dy*.1,
					0, twopi)
				gc.Close()
				gc.FillStroke()
			case 'W':
			case 't', 'T':
				// draw a little boxes for towns
				mt := 8
				if mapChar == 'T' {
					mt = 16
				}
				cc := &image.Uniform{map_town_fill}
				for ii := 0; ii < mt; ii++ {
					tx := x*int(dx) + rand.Intn(int(dx))
					ty := y*int(dy) + rand.Intn(int(dy))
					draw.Draw(img,
						image.Rectangle{
							image.Point{tx, ty},
							image.Point{tx + rand.Intn(16), ty + rand.Intn(16)},
						},
						cc,
						image.Point{0, 0},
						draw.Src,
					)
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
			switch r.mw.mapData.Data[i] {
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

	// paint the units
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			forces := r.mw.grid.Units(int32(x), int32(y))
			/*
				if len(forces.commands) > 0 || len(forces.units) > 0 {
					println("there are ", len(forces.commands), "commands and", len(forces.units), "units at", x+1, y+1)
				}
			*/
			// draw the commands - 3 per line
			if cnt := len(forces.commands); cnt > 0 {
				fx := float64(x) * dx
				fy := float64(y) * dy
				blocksize := dx / 3.0
				radius := dx / 10
				xx := 0.0
				yy := 0.0
				for i := 0; i < cnt; i++ {
					if (i % 3) == 0 {
						xx = 0.0
						yy += blocksize
					}
					icon := forces.commands[i]
					gc.SetFillColor(map_unit_fill)
					gc.SetStrokeColor(map_unit_stroke)
					gc.SetLineWidth(2)
					gc.BeginPath()
					if icon.selected {
						gc.SetFillColor(map_unit_selected_fill)
						gc.SetStrokeColor(map_unit_selected_stroke)
					}
					forces.commands[i].rect = image.Rectangle{
						Min: image.Point{X: int((fx + xx) / scale), Y: int((fy + yy) / scale)},
						Max: image.Point{X: int((fx + xx + blocksize) / scale), Y: int((fy + yy + blocksize) / scale)},
					}
					draw2dkit.RoundedRectangle(gc,
						fx+xx+2, fy+yy+2,
						fx+xx+blocksize-4, fy+yy+blocksize,
						radius, radius)
					gc.Close()
					gc.FillStroke()

					// denote the type
					gc.SetStrokeColor(denote_unit)
					gc.SetFillColor(denote_unit)
					gc.SetLineWidth(dx / 30)
					gc.SetLineCap(draw2d.RoundCap)
					cmd := icon.cmd
					if len(cmd.Units) > 0 {
						switch cmd.Arm {
						case republique.Arm_CAVALRY:
							gc.MoveTo(fx+xx+blocksize-6, fy+yy+4)
							gc.LineTo(fx+xx+4, fy+yy+blocksize-6)
							gc.Stroke()
						case republique.Arm_INFANTRY:
							gc.MoveTo(fx+xx+blocksize-6, fy+yy+4)
							gc.LineTo(fx+xx+4, fy+yy+blocksize-6)
							gc.Stroke()
							gc.MoveTo(fx+xx+4, fy+yy+4)
							gc.LineTo(fx+xx+blocksize-6, fy+yy+blocksize-6)
							gc.Stroke()
						case republique.Arm_ARTILLERY:
							draw2dkit.Circle(gc, fx+xx+(blocksize/2)+2, fy+yy+(blocksize/2)+2, dx/20)
							gc.Fill()
						}
					}
					xx += blocksize
				}
			}
		}
	}

	return img
}

func (r *mapRender) hills(gc *draw2dimg.GraphicContext, fx, fy, dx, dy float64, count int) {
	gc.SetFillColor(map_hill_fill)
	gc.SetStrokeColor(map_hill_stroke)
	gc.SetLineWidth(2)
	count = rand.Intn(2) + 1
	switch count {
	case 1:
		gc.BeginPath()
		gc.MoveTo(fx+0.1*dx, fy+0.5*dy)
		gc.QuadCurveTo(fx+0.3*dx, fy+0.2*dy, fx+0.7*dx, fy+0.5*dy)
		gc.Close()
		gc.FillStroke()
		gc.MoveTo(fx+0.3*dx, fy+0.7*dy)
		gc.QuadCurveTo(fx+0.5*dx, fy+0.4*dy, fx+0.9*dx, fy+0.7*dy)
		gc.Close()
		gc.FillStroke()
	case 2:
		gc.BeginPath()
		gc.MoveTo(fx+0.4*dx, fy+0.3*dy)
		gc.QuadCurveTo(fx+0.6*dx, fy+0.1*dy, fx+0.8*dx, fy+0.3*dy)
		gc.Close()
		gc.FillStroke()
		gc.BeginPath()
		gc.MoveTo(fx+0.1*dx, fy+0.5*dy)
		gc.QuadCurveTo(fx+0.3*dx, fy+0.2*dy, fx+0.7*dx, fy+0.5*dy)
		gc.Close()
		gc.FillStroke()
		gc.MoveTo(fx+0.3*dx, fy+0.7*dy)
		gc.QuadCurveTo(fx+0.5*dx, fy+0.4*dy, fx+0.9*dx, fy+0.7*dy)
		gc.Close()
		gc.FillStroke()
	}
}
