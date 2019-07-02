package appwindow

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/steveoc64/memdebug"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	republique "github.com/steveoc64/republique5/proto"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/llgcode/draw2d/draw2dimg"
)

type mapRender struct {
	render     *canvas.Raster
	mw         *MapWidget
	objects    []fyne.CanvasObject
	background *image.RGBA
	foreground *image.RGBA
	dirty      bool
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
	if r.background == nil || r.background.Bounds().Size().X != w || r.background.Bounds().Size().Y != h {
		r.background = r.generateBackground(w, h)
		r.dirty = true
	}
	if r.dirty {
		r.foreground = r.generateForeground(w, h)
		r.dirty = false
	}
	if r.mw.hidden {
		return &image.RGBA{}
	}
	return r.foreground
}

func (r *mapRender) generateBackground(w, h int) *image.RGBA {
	t1 := time.Now()
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
	mapSide := r.mw.app.MapData.Side

	convertXY := func(x, y int) (int, int) {
		switch mapSide {
		case republique.MapSide_TOP:
			return mx - x - 1, my - y - 1
		}
		return x, y
	}

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
	for ya := 0; ya < my; ya++ {
		for xa := 0; xa < mx; xa++ {
			x, y := convertXY(xa, ya)
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
			case 'w', 'W':
				// draw a little circle for woods
				for _, v := range r.mw.grid.things[i] {
					gc.SetFillColor(map_woods_fill)
					//gc.SetStrokeColor(map_woods_stroke)
					gc.SetLineWidth(2)
					cx := fx + dx*float64(v.x)/100.0
					cy := fy + dy*float64(v.y)/100.0
					rx := dx * float64(v.size) / 100.0
					ry := dy * float64(v.size) / 100.0
					gc.BeginPath()
					gc.ArcTo(cx, cy, rx, ry, 0.0, twopi)
					gc.Close()
					gc.Fill()
				}
			case 't', 'T':
				// draw a little boxes for towns
				for _, v := range r.mw.grid.things[i] {
					gc.SetFillColor(map_town_fill)
					gc.SetStrokeColor(map_town_stroke)
					gc.SetLineWidth(1)
					x1 := fx + dx*float64(v.x)/100.0
					y1 := fy + dy*float64(v.y)/100.0
					x2 := x1 + dx*float64(v.size)/100.0
					y2 := y1 + dy*float64(v.size)/100.0
					gc.BeginPath()
					draw2dkit.Rectangle(gc, x1, y1, x2, y2)
					gc.FillStroke()
				}
			}
			i++
		}
	}

	// draw rivers
	if false {
		i = 0
		gc.SetFillColor(map_deep_blue)
		gc.SetStrokeColor(map_blue)
		gc.SetLineWidth(20)
		gc.BeginPath()
		gotriver := false
		for ya := 0; ya < my; ya++ {
			for xa := 0; xa < mx; xa++ {
				x, y := convertXY(xa, ya)
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
	}

	// do the river segments
	gc.SetFillColor(map_deep_blue)
	gc.SetStrokeColor(map_blue)
	gc.SetLineWidth(20)
	for k, v := range r.mw.rivers {
		gc.BeginPath()
		fx := (float64(k.x) + 0.5) * dx
		fy := (float64(k.y) + 0.5) * dy
		if len(v.adjacent) == 0 {
			draw2dkit.Ellipse(gc, fx, fy, dx/3, dy/5)
			gc.FillStroke()
			continue
		}
		for kk, _ := range v.adjacent {
			// double check that this segment isnt already done
			if other, ok := r.mw.rivers[riverPoint{kk.x, kk.y}]; ok {
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
		case k.x == 0, k.x == r.mw.grid.x-1:
			fx2 := 0.0
			if k.x == r.mw.grid.x-1 {
				fx2 = float64(r.mw.grid.x) * dx
			}
			bump := dx / -4
			for kk, _ := range v.adjacent {
				switch {
				case kk.y < k.y:
					bump += dy / 2
				case kk.y > k.y:
					bump += dy / -2
				}
				gc.MoveTo(fx, fy)
				gc.LineTo(fx2, fy+bump)
				gc.Stroke()
				break
			}
		case k.y == 0, k.y == r.mw.grid.y-1:
			fy2 := 0.0
			if k.y == r.mw.grid.y-1 {
				fy2 = float64(r.mw.grid.y) * dy
			}
			bump := dx / -4
			for kk, _ := range v.adjacent {
				switch {
				case kk.x < k.x:
					bump += dx / 2
				case kk.x > k.x:
					bump += dx / -2
				}
				gc.MoveTo(fx, fy)
				gc.LineTo(fx+bump, fy2)
				gc.Stroke()
				break
			}
		}
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

	memdebug.Print(t1, "rendered terrain")
	return img
}

func (r *mapRender) generateForeground(w, h int) *image.RGBA {
	t1 := time.Now()
	scale := float64(r.Scale())
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}
	// blit the background across to the foreground
	draw.Draw(img,
		r.background.Rect,
		r.background,
		image.Point{0, 0},
		draw.Src)
	memdebug.Print(t1, "blit terrain")
	dx := float64(w / int(r.mw.mapData.X))
	dy := float64(h / int(r.mw.mapData.Y))
	mx := int(r.mw.mapData.X)
	my := int(r.mw.mapData.Y)
	gc := draw2dimg.NewGraphicContext(img)
	mapSide := r.mw.app.MapData.Side

	convertXY := func(x, y int) (int, int) {
		switch mapSide {
		case republique.MapSide_TOP:
			return mx - x - 1, my - y - 1
		}
		return x, y
	}

	// paint the units
	for ya := 0; ya < my; ya++ {
		for xa := 0; xa < mx; xa++ {
			x, y := convertXY(xa, ya)
			forces := r.mw.grid.Units(int32(xa), int32(ya))
			// draw the commands - 3 per line
			if cnt := len(forces.commands); cnt > 0 {
				fx := float64(x) * dx
				fy := float64(y) * dy
				blocksize := dx / 3.0
				radius := dx / 10
				xx := 0.0
				yy := 0.0
				for i := 0; i < cnt; i++ {
					if i > 0 && (i%3) == 0 {
						xx = 0.0
						yy += blocksize
					}
					icon := forces.commands[i]
					if icon.cmd == nil {
						continue
					}
					// draw the lines of action
					gc.SetStrokeColor(map_unit_orders_stroke)
					gc.SetLineWidth(18)
					gc.SetLineJoin(draw2d.BevelJoin)
					gc.SetLineCap(draw2d.RoundCap)
					gc.MoveTo(fx+xx+(blocksize/2), fy+yy+(blocksize/2))
					doPath := true
					switch icon.cmd.GetGameState().GetOrders() {
					case republique.Order_MOVE, republique.Order_MARCH:
						gc.SetStrokeColor(map_unit_orders_march)
						//gc.SetLineDash([]float64{0.2, 0.4, 0.6, 0.8}, 0.0)
						gc.SetLineWidth(dx / 6)
					case republique.Order_ATTACK:
						gc.SetStrokeColor(map_unit_orders_attack)
						gc.SetLineWidth(dx / 3)
					case republique.Order_ENGAGE:
						gc.SetStrokeColor(map_unit_orders_engage)
						gc.SetLineWidth(dx / 4)
					case republique.Order_CHARGE:
						gc.SetStrokeColor(map_unit_orders_charge)
						gc.SetLineWidth(dx / 2)
					case republique.Order_FIRE:
						gc.SetFillColor(map_unit_orders_fire)
						if len(icon.cmd.GetGameState().GetObjective()) == 2 {
							target := icon.cmd.GameState.Objective[1]
							tx, ty := convertXY(int(target.X), int(target.Y))
							gc.LineTo(float64(tx-1)*dx, float64(ty-1)*dy+0.5*dy)
							gc.LineTo(float64(tx)*dx, float64(ty-1)*dy+0.5*dy)
							gc.Close()
							gc.Fill()
							doPath = false
						}
					case republique.Order_PURSUIT:
						gc.SetStrokeColor(map_unit_orders_pursuit)
						gc.SetLineWidth(dx / 6)
					}
					if doPath {
						for k, path := range icon.cmd.GetGameState().Objective {
							if k > 0 { // burn the first one
								px, py := int(path.X), int(path.Y)
								switch mapSide {
								case republique.MapSide_TOP:
									px = mx - px + 1
									py = my - py + 1
								}
								gc.LineTo(float64(px)*dx-(0.5*dx), float64(py)*dy-(0.5*dy))
							}
						}
						gc.Stroke()
					}

					// draw the basic rect
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

					// denote order status
					if icon.cmd.GetGameState().GetCan().GetOrder() {
						gc.SetFillColor(map_unit_can_order)
						if icon.cmd.GetGameState().GetHas().GetOrder() &&
							icon.cmd.GetGameState().GetOrders() != republique.Order_RESTAGE {
							gc.SetFillColor(map_unit_has_order)
						}
						draw2dkit.Rectangle(gc,
							fx+xx+2, fy+yy+blocksize-10,
							fx+xx+blocksize-4, fy+yy+blocksize-4)
						gc.Fill()
					}

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

	// paint the enemy
	for ya := 0; ya < my; ya++ {
		for xa := 0; xa < mx; xa++ {
			x, y := convertXY(xa, ya)
			e := []*republique.Command{}
			for _, enemy := range r.mw.app.Enemy {
				if enemy.GameState.Grid.X == int32(xa+1) && enemy.GameState.Grid.Y == int32(ya+1) {
					e = append(e, enemy)
				}
			}
			fx := float64(x) * dx
			fy := float64(y) * dy
			blocksize := dx / 3.0
			radius := dx / 10
			xx := 0.0
			yy := 0.0
			for i, enemy := range e {
				if i > 0 && (i%3) == 0 {
					xx = 0.0
					yy += blocksize
				}
				// draw the basic rect
				gc.SetFillColor(map_enemy_fill)
				gc.SetStrokeColor(map_unit_stroke)
				gc.SetLineWidth(2)
				gc.BeginPath()
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
				switch enemy.Arm {
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
				xx += blocksize
			}
		}
	}
	memdebug.Print(t1, "rendered units")
	return img
}

func (r *mapRender) hills(gc *draw2dimg.GraphicContext, fx, fy, dx, dy float64, count int) {
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

func (r *mapRender) ConvertToGrid(event *fyne.PointEvent) (int32, int32) {
	if r.background == nil {
		return 0, 0
	}
	scale := r.Scale()
	size := r.background.Bounds()
	dx := (float32(size.Max.X) / scale) / float32(r.mw.grid.x)
	dy := (float32(size.Max.Y) / scale) / float32(r.mw.grid.y)
	x := int32(float32(event.Position.X) / dx)
	y := int32(float32(event.Position.Y) / dy)
	return x + 1, y + 1
}
