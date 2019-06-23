package appwindow

import (
	"image"
	"image/color"

	"github.com/vdobler/chart/imgg"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/vdobler/chart"
)

type plotRender struct {
	render  *canvas.Raster
	p       *PlotWidget
	objects []fyne.CanvasObject
	img     *image.RGBA
}

func newPlotRender(p *PlotWidget) *plotRender {
	r := &plotRender{p: p}
	render := canvas.NewRaster(r.getImage)
	//render := canvas.NewRasterWithPixels(r.paint2)
	r.render = render
	r.objects = []fyne.CanvasObject{render}
	return r
}

func (r *plotRender) Scale() float32 {
	return fyne.CurrentApp().Driver().CanvasForObject(r.render).Scale()
}

// ApplyTheme applies the theme
func (r *plotRender) ApplyTheme() {
	// noop
}

// BackgroundColor returns the background color for our map
func (r *plotRender) BackgroundColor() color.Color {
	return color.RGBA{183, 211, 123, 1}
}

// Destroy removes any resources we have on this renderer
func (r *plotRender) Destroy() {
	// noop
}

// Layout does .. the layout ?
func (r *plotRender) Layout(size fyne.Size) {
	r.render.Resize(size)
}

// MinSize returns the minimum size for this renderer
func (r *plotRender) MinSize() fyne.Size {
	return fyne.Size{
		Width:  300,
		Height: 300,
	}
}

// Objects returns the slice of objects that we own
func (r *plotRender) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh paints the map
func (r *plotRender) Refresh() {
	canvas.Refresh(r.render)
}

func (r *plotRender) getImage(w, h int) image.Image {
	if r.img == nil || r.img.Bounds().Size().X != w || r.img.Bounds().Size().Y != h {
		r.img = r.generateImage(w, h)
	}
	if r.p.hidden {
		return &image.RGBA{}
	}
	return r.img
}

func (r *plotRender) generateImage(w, h int) *image.RGBA {
	img := r.img
	img = image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	ebit := chart.BarChart{Title: r.p.title}
	ebit.XRange.Category = []string{"06:00", "06:30", "07:00", "07:30", "08:00", "08:30", "09:00", "09:30"}
	ebit.XRange.Label, ebit.YRange.Label = "Game Performance", "% Status"
	ebit.Key.Pos, ebit.Key.Cols, ebit.Key.Border = "otc", 2, -1
	ebit.YRange.ShowZero = true
	ebit.ShowVal = 0
	ebit.AddDataPair("Strength", []float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{100, 100, 80, 78, 78, 78, 65, 65},
		chart.Style{Symbol: '#', LineColor: color.NRGBA{0x30, 0x30, 0xff, 0xff}, LineWidth: 2, FillColor: color.NRGBA{0x07, 0x71, 0x3c, 0xff}})
	ebit.AddDataPair("Morale", []float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{100, 100, 40, 60, 65, 70, 75, 90},
		chart.Style{Symbol: 'O', LineColor: color.NRGBA{0x44, 0x11, 0x44, 0xff}, LineWidth: 2, FillColor: color.NRGBA{0x5a, 0x22, 0x8a, 0xff}})
	igr := imgg.AddTo(img, 0, 0, w, h, plot_background, nil, nil)
	ebit.Plot(igr)
	r.img = img
	return r.img

	/*
		if len(r.p.data) == 0 {
			return img
		}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}
		l, err := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
		if err != nil {
			panic(err)
		}
		p.Add(l)

		// Draw the plot to an in-memory image.
		c := vgimg.NewWith(vgimg.UseImage(img))
		p.Draw(vdraw.New(c))
		return img
		/*

		p.Title.Text = r.p.title
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		for _,v := range r.p.data {
			err = plotutil.AddLinePoints(p, v.name, v)
			if err != nil {
				panic(err)
			}
		}

		// Save the image.
		f, err := os.Create("plot.png")
		if err != nil {
			panic(err)
		}
		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}

		return img

	*/
}
