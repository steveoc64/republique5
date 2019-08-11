package republique

import (
	"fmt"
	"math"
)

// Waypoint data for each segment of a move order
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
	Prep     bool
}

// GetValue gets the grid contents at x,y
func (m *MapData) GetValue(x, y int32) byte {
	offset := x + y*m.X
	if offset < 0 || offset > m.X*m.Y {
		return ' '
	}
	return m.Data[offset]
}

// GetWaypoints gets a slice of waypoints for a given command on this map
func (m *MapData) GetWaypoints(command *Command) []Waypoint {
	waypoints := []Waypoint{}
	x := command.GetGameState().GetGrid().GetX()
	y := command.GetGameState().GetGrid().GetY()

	switch command.GetGameState().GetOrders() {
	case Order_RESTAGE, Order_NO_ORDERS, Order_RALLY:
		return waypoints
	case Order_FIRE:
		// do they need to deploy ?
		waypoint := Waypoint{
			X:        x,
			Y:        y,
			X2:       x,
			Y2:       y,
			Speed:    1.0,
			Elapsed:  5.0,
			Turns:    0,
			Going:    "ready",
			Distance: 0.0,
			Path:     "Deployed, Ready to fire",
		}
		for _, v := range command.GetUnits() {
			if v.Arm == Arm_ARTILLERY && !v.GetGameState().GetGunsDeployed() {
				waypoint.Turns = 1
				waypoint.Going = "deploying"
				waypoint.Path = "-> Deploy for bombardment  (half a turn)"
				waypoint.Prep = true
				waypoint.Elapsed = 10.0
				break
			}
		}
		return append(waypoints, waypoint)
	case Order_ENGAGE:
		switch command.GetGameState().GetFormation() {
		case Formation_MARCH_COLUMN, Formation_RESERVE:
			waypoints = append(waypoints, Waypoint{
				X:        x,
				Y:        y,
				X2:       x,
				Y2:       y,
				Speed:    1.0,
				Elapsed:  20.0,
				Turns:    1,
				Going:    "preparing",
				Distance: 0.0,
				Path:     "-> Form line of Battle (about 1 turn)",
				Prep:     true,
			})
		}
	case Order_MARCH:
		switch command.Arm {
		case Arm_ARTILLERY:
			for _, v := range command.GetUnits() {
				if v.Arm == Arm_ARTILLERY && v.GetGameState().GetGunsDeployed() {
					waypoints = append(waypoints, Waypoint{
						X:        x,
						Y:        y,
						X2:       x,
						Y2:       y,
						Speed:    1.0,
						Elapsed:  10.0,
						Turns:    1,
						Going:    "limbering",
						Distance: 0.0,
						Path:     "-> Limber ready to move (half a turn)",
						Prep:     true,
					})
					break
				}
			}
		case Arm_INFANTRY:
			if command.GetRank() != Rank_CORPS && command.GetRank() != Rank_ARMY {
				switch command.GetGameState().GetFormation() {
				case Formation_MARCH_COLUMN, Formation_COLUMN:
				default:
					waypoints = append(waypoints, Waypoint{
						X:        x,
						Y:        y,
						X2:       x,
						Y2:       y,
						Speed:    1.0,
						Elapsed:  10.0,
						Turns:    1,
						Going:    "prepare to march",
						Distance: 0.0,
						Path:     "-> Form march columns (about 1 turn)",
						Prep:     true,
					})
				}
			}
		}
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
			case command.Arm == Arm_INFANTRY && command.GetGameState().GetOrders() == Order_MARCH:
				speed = 1.4
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
			turns := int(elapsed+10.0) / 20
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
	// then add a deploy option at the end
	switch command.GetGameState().GetOrders() {
	case Order_MARCH:
		switch command.Arm {
		case Arm_ARTILLERY:
			waypoints = append(waypoints, Waypoint{
				X:        x,
				Y:        y,
				X2:       x,
				Y2:       y,
				Speed:    1.0,
				Elapsed:  10.0,
				Turns:    1,
				Going:    "deploying",
				Distance: 0.0,
				Path:     "-> Unlimber, prepare for action (half a turn)",
				Prep:     true,
			})
		case Arm_INFANTRY:
			waypoints = append(waypoints, Waypoint{
				X:        x,
				Y:        y,
				X2:       x,
				Y2:       y,
				Speed:    1.0,
				Elapsed:  10.0,
				Turns:    1,
				Going:    "deploying",
				Distance: 0.0,
				Path:     "-> Deploy back to " + upString(command.GetGameState().GetFormation().String()) + " (about 1 turn)",
				Prep:     true,
			})
		}
	}
	return waypoints

}
