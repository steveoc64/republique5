package appwindow

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/davecgh/go-spew/spew"
	"github.com/llgcode/draw2d/draw2dimg"
	"k8s.io/apimachinery/pkg/util/rand"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
)

type mapRender struct {
	render  *canvas.Raster
	mw      *MapWidget
	objects []fyne.CanvasObject
	img     *image.RGBA
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
	println("renderer layout", size.Width, size.Height)
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
	println("maprender refresh")
	canvas.Refresh(r.mw)
}

func (r *mapRender) getImage(w, h int) image.Image {
	println("been asked to getImage and hidden =", r.mw.hidden, w, h)
	if r.img != nil {
		spew.Dump(r.img.Bounds())
	}
	if r.img == nil || r.img.Bounds().Size().X != w || r.img.Bounds().Size().Y != h {
		r.img = r.generateImage(w, h)
	}
	if r.mw.hidden {
		return &image.RGBA{}
	}
	return r.img
}

func (r *mapRender) generateImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	dx := float64(w / int(r.mw.mapData.X))
	dy := float64(h / int(r.mw.mapData.Y))
	mx := int(r.mw.mapData.X)
	my := int(r.mw.mapData.Y)
	gc := draw2dimg.NewGraphicContext(img)

	// grid overlays
	i := 0
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			fx := float64(x) * dx
			fy := float64(y) * dy
			c := color.RGBA{uint8(rand.Intn(64)), 200, 50, uint8(rand.Intn(16))}
			draw.Draw(img,
				image.Rectangle{image.Point{x * int(dx), y * int(dy)},
					image.Point{(x + 1) * int(dx), (y + 1) * int(dy)}},
				&image.Uniform{c},
				image.Point{0, 0},
				draw.Src)

			mapChar := r.mw.mapData.Data[i]
			switch mapChar {
			case 'H':
				// small hills
				r.hills(gc, fx, fy, dx, dy, 1)
			case 'h':
				// larger hills
				r.hills(gc, fx, fy, dx, dy, 2)
			case 'w':
				// draw a little circle for woods
				gc.SetFillColor(map_woods_fill)
				gc.SetStrokeColor(map_woods_stroke)
				gc.SetLineWidth(5)
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.1, dy*.1,
					0, 6)
				gc.Close()
				gc.FillStroke()
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.15, dy*.15,
					0, 6)
				gc.Close()
				gc.FillStroke()
				gc.BeginPath()
				gc.ArcTo(fx+float64(rand.Intn(100))/100.0*dx,
					fy+float64(rand.Intn(100))/100.0*dy,
					dx*.1, dy*.1,
					0, 6)
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
	gc.SetFillColor(map_water)
	gc.SetStrokeColor(map_blue)
	gc.SetLineWidth(8)
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

	// grid lines - vertical
	gc.SetStrokeColor(map_grid)
	gc.SetLineWidth(1)
	for x := 0; x < mx; x++ {
		gc.BeginPath()
		gc.MoveTo(float64(x)*dx, 0.0)
		gc.LineTo(float64(x)*dx, float64(h))
		gc.Close()
		gc.FillStroke()
	}
	// grid lines - horizontal
	for y := 0; y < my; y++ {
		gc.BeginPath()
		gc.MoveTo(0.0, float64(y)*dy)
		gc.LineTo(float64(w), float64(y)*dy)
		gc.Close()
		gc.FillStroke()
	}

	return img
}

func (r *mapRender) hills(gc *draw2dimg.GraphicContext, fx, fy, dx, dy float64, count int) {
	gc.SetFillColor(map_hill_fill)
	gc.SetStrokeColor(map_hill_stroke)
	gc.SetLineWidth(5)
	switch count {
	case 1:
		gc.BeginPath()
		gc.MoveTo(fx+0.2*dx, fy+0.6*dy)
		gc.LineTo(fx+0.5*dx, fy+0.4*dy)
		gc.LineTo(fx+0.8*dx, fy+0.6*dy)
		gc.Close()
		gc.FillStroke()
	case 2:
		gc.BeginPath()
		gc.MoveTo(fx+0.1*dx, fy+0.5*dy)
		gc.LineTo(fx+0.4*dx, fy+0.3*dy)
		gc.LineTo(fx+0.7*dx, fy+0.5*dy)
		gc.Close()
		gc.FillStroke()
		gc.BeginPath()
		gc.MoveTo(fx+0.3*dx, fy+0.7*dy)
		gc.LineTo(fx+0.6*dx, fy+0.5*dy)
		gc.LineTo(fx+0.9*dx, fy+0.7*dy)
		gc.Close()
		gc.FillStroke()
	}
}
