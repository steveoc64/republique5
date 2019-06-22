package appwindow

import (
	"image"
	"image/color"
	"math/rand"

	"fyne.io/fyne"

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

func (g *gridForces) Select(id int32) (gotCommand *rp.Command, gotSome bool) {
	for _, v := range g.commands {
		v.selected = (v.cmd.Id == id)
		if v.selected {
			gotSome = true
			gotCommand = v.cmd
		}
	}
	return gotCommand, gotSome
}

type gridData struct {
	x, y   int32
	back   []color.RGBA
	value  []byte
	units  []gridForces
	things []thingDataArray
}

type thingData struct {
	x    int
	y    int
	size int
}

type thingDataArray []thingData

func newGridData(x, y int32, values string) *gridData {
	g := &gridData{
		x:      x,
		y:      y,
		back:   make([]color.RGBA, x*y),
		value:  []byte(values),
		units:  make([]gridForces, x*y),
		things: make([]thingDataArray, x*y),
	}
	for i := 0; i < int(x*y); i++ {
		g.back[i] = color.RGBA{uint8(rand.Intn(40) + 160), uint8(rand.Intn(40) + 180), uint8(rand.Intn(40) + 100), 200}
		switch g.value[i] {
		case 'w', 't':
			trees := []thingData{}
			for t := 0; t < 10+rand.Intn(20); t++ {
				trees = append(trees, thingData{
					x:    rand.Intn(100),
					y:    rand.Intn(100),
					size: 2 + rand.Intn(10),
				})
			}
			g.things[i] = trees
		case 'W', 'T':
			trees := []thingData{}
			for t := 0; t < 25+rand.Intn(40); t++ {
				trees = append(trees, thingData{
					x:    rand.Intn(100),
					y:    rand.Intn(100),
					size: 2 + rand.Intn(10),
				})
			}
			g.things[i] = trees
		}
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

func (g *gridData) Select(id int32) (gotCommand *rp.Command, gotSome bool) {
	for _, v := range g.units {
		if cmd, ok := v.Select(id); ok {
			gotSome = true
			gotCommand = cmd
		}
	}
	return gotCommand, gotSome
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
