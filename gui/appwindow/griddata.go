package appwindow

import (
	"fyne.io/fyne"
	"image"
	"image/color"
	"math/rand"

	rp "github.com/steveoc64/republique5/proto"
)

type gridIcon struct {
	cmd      *rp.Command
	rect     image.Rectangle
	selected bool
}
type gridForces struct {
	commands []*gridIcon
	units    []*rp.Unit
}

func (g *gridForces) CommandAt(p fyne.Position) *rp.Command {
	pp := image.Point{X: int(p.X), Y: int(p.Y)}
	for _, v := range g.commands {
		if pp.In(v.rect) {
			return v.cmd
		}
	}
	return nil
}

func (g *gridForces) Select(id int32) bool {
	gotSome := false
	for _, v := range g.commands {
		v.selected = (v.cmd.Id == id)
		if v.selected {
			gotSome = true
		}
	}
	return gotSome
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

func (g *gridData) CommandAt(p fyne.Position) *rp.Command {
	for _, v := range g.units {
		if c := v.CommandAt(p); c != nil {
			return c
		}
	}
	return nil
}

func (g *gridData) Select(id int32) bool {
	gotSome := false
	for _, v := range g.units {
		gotSome = gotSome || v.Select(id)
	}
	return gotSome
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
	g.units[i].commands = append(g.units[i].commands, &gridIcon{cmd: c})
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
