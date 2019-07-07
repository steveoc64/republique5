package republique

import (
	"fmt"
	"math"
)

type Waypoint struct {
	X        int32
	Y        int32
	X2       int32
	Y2       int32
	Speed    float64
	Elapsed  float64
	Turns    int
	Going    string
	Distance float64
	Path     string
}

func (m *MapData) GetValue(x, y int32) byte {
	offset := x + y*m.X
	if offset < 0 || offset > m.X*m.Y {
		return ' '
	}
	return m.Data[offset]
}

func (m *MapData) GetWaypoints(command *Command) []Waypoint {
	waypoints := []Waypoint{}
	x := command.GetGameState().GetGrid().GetX()
	y := command.GetGameState().GetGrid().GetY()

	switch command.GetGameState().GetOrders() {
	case Order_RESTAGE, Order_FIRE, Order_NO_ORDERS, Order_RALLY:
		return waypoints
	}
	for k, v := range command.GetGameState().GetObjective() {
		if k > 0 {
			distance := math.Sqrt(float64((v.X-x)*(v.X-x) + (v.Y-y)*(v.Y-y)))
			speed := 1.0
			switch {
			case command.Rank == Rank_CORPS, command.Rank == Rank_ARMY:
				speed = 2.0
			case command.Arm == Arm_CAVALRY:
				speed = 1.5
			}

			//fromGrid := o.app.mapPanel.mapWidget.grid.Value(int32(x-1), int32(y-1))
			fromGrid := m.GetValue(x-1, y-1)
			switch fromGrid {
			case 't', 'w':
				speed *= 0.8
			case 'T', 'W', 'h', 'r':
				speed *= 0.6
			case 'H':
				speed *= 0.5
			}
			//toGrid := o.app.mapPanel.mapWidget.grid.Value(int32(v.X-1), int32(v.Y-1))
			toGrid := m.GetValue(v.X-1, v.Y-1)
			switch toGrid {
			case 't', 'w':
				speed *= 0.7
			case 'T', 'W', 'h':
				speed *= 0.5
			case 'H', 'r':
				speed *= 0.4
			}
			going := "at a good march"
			if command.Arm == Arm_CAVALRY {
				if command.GetGameState().GetOrders() == Order_CHARGE {
					speed *= 1.5
				}
				going = "at the trot"
				switch {
				case speed >= 2.0:
					going = "at the gallop"
				case speed >= 1.5:
					going = "at the trot"
				case speed <= 0.4:
					going = "very slow going"
				case speed <= 0.5:
					going = "difficult going"
				case speed <= 0.6:
					going = "at a slow walk"
				case speed <= 0.7:
					going = "at the walk"
				case speed <= 0.8:
					going = "at the walk"
				}
			} else {
				switch {
				case speed >= 1.5:
					going = "with great speed"
				case speed <= 0.4:
					going = "very slow"
				case speed <= 0.5:
					going = "harsh terrain"
				case speed <= 0.6:
					going = "slow going"
				case speed <= 0.7:
					going = "with some delays"
				case speed <= 0.8:
					going = "with minor delays"
				}
			}

			elapsed := ((distance / (speed * 3.0)) * 60.0) // 20mins to the mile in good order
			turns := int(elapsed) / 20
			turnss := ""
			if turns != 1 {
				turnss = "s"
			}
			path := fmt.Sprintf("-> %d,%d  (%0.1f miles %s, about %d turn%s)", v.X, v.Y, distance, going, turns, turnss)
			waypoints = append(waypoints, Waypoint{
				X:        x,
				Y:        y,
				X2:       v.X,
				Y2:       v.Y,
				Distance: distance,
				Going:    going,
				Turns:    turns,
				Elapsed:  elapsed,
				Speed:    speed,
				Path:     path,
			})
			x = v.X
			y = v.Y
		}
	}
	return waypoints

}
