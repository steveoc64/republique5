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
	canvas  fyne.CanvasObject
}

func newMapRender(mw *MapWidget) *mapRender {
	r := &mapRender{mw: mw}
	render := canvas.NewRaster(r.getImage)
	//render := canvas.NewRasterWithPixels(r.paint2)
	r.render = render
	r.objects = []fyne.CanvasObject{render}
	r.ApplyTheme()
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
	return r.img
}

func (r *mapRender) generateImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	dx := w / int(r.mw.mapData.X)
	dy := h / int(r.mw.mapData.Y)
	mx := int(r.mw.mapData.X)
	my := int(r.mw.mapData.Y)

	// grid overlays
	gc := draw2dimg.NewGraphicContext(img)
	i := 0
	for x := 0; x < mx; x++ {
		for y := 0; y < my; y++ {
			c := color.RGBA{uint8(rand.Intn(64)), 200, 50, uint8(rand.Intn(16))}
			draw.Draw(img,
				image.Rectangle{image.Point{x * dx, y * dy},
					image.Point{(x + 1) * dx, (y + 1) * dy}},
				&image.Uniform{c},
				image.Point{0, 0},
				draw.Src)

			gc.BeginPath()
			switch r.mw.mapData.Data[i] {
			case 'h':
				// draw a silly little triangle in the middle of the grid square
				gc.SetStrokeColor(color.RGBA{50, 40, 30, 64})
				gc.SetLineWidth(5)
				gc.MoveTo(float64(x*mx+mx/3), float64(y*my+2*my/3))
				gc.LineTo(float64(x*mx+mx/2), float64(y*my+my/3))
				gc.LineTo(float64(x*mx+2*mx/3), float64(y*my+2*my/3))
			case 'H':
			case 'w':
				gc.SetStrokeColor(color.RGBA{0, 40, 0, 64})
				gc.SetFillColor(color.RGBA{0, 80, 0, 64})
				gc.SetLineWidth(5)
				gc.C
			case 'W':
			case 't':
			case 'T':
			case 'r':
			case 'R':
			}
			gc.Close()
			gc.FillStroke()

		}
	}

	// grid lines - vertical
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0x08})
	gc.SetLineWidth(1)
	for x := 0; x < mx; x++ {
		gc.BeginPath()
		gc.MoveTo(float64(x*dx), 0.0)
		gc.LineTo(float64(x*dx), float64(h))
		gc.Close()
		gc.FillStroke()
	}
	// grid lines - horizontal
	for y := 0; y < my; y++ {
		gc.BeginPath()
		gc.MoveTo(0.0, float64(y*dy))
		gc.LineTo(float64(w), float64(y*dy))
		gc.Close()
		gc.FillStroke()
	}

	return img
}
