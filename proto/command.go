package republique

import (
	"fmt"
	"math/rand"
	"strings"
)

// LabelString returns a formatted string for rendering the whole command
// in the GUI app
func (c *Command) LabelString() string {
	s := fmt.Sprintf("%s - %s [%s]", c.Name, c.CommanderName, upString(strings.ToLower(c.GetGameState().GetPosition().String())))
	if c.Notes != "" {
		s = s + " (" + c.Notes + ")"
	}
	if c.Subcommands == nil && c.Arm == Arm_INFANTRY && len(c.Units) > 2 {
		s = s +
			" by " +
			upString(strings.Replace(strings.ToLower(c.GetGameState().GetFormation().String()), "_", " ", 1)) +
			"s of Bde's"
	}
	return s
}

// LongDoscription returns a long description string for the command
func (c *Command) LongDescription() string {
	s := fmt.Sprintf("[%d] %s - %s (+%d) %s",
		c.Id,
		c.Name,
		c.CommanderName,
		c.CommanderBonus,
		c.GetCommandStrengthLabel())

	return s
}

// CommandStrength is the number of troops in a command
type CommandStrength struct {
	Infantry int32
	Cavalry  int32
	Guns     int32
}

func (c *Command) GetCommandStrength() CommandStrength {
	s := CommandStrength{}
	for _, unit := range c.Units {
		switch unit.Arm {
		case Arm_INFANTRY:
			s.Infantry += unit.Strength * 550
		case Arm_CAVALRY:
			s.Cavalry += unit.Strength * 300
		case Arm_ARTILLERY:
			s.Guns += unit.Strength * 12
		}
	}
	return s
}

func (c *Command) GetCommandStrengthLabel() string {
	retval := []string{}
	s := c.GetCommandStrength()
	if s.Infantry > 0 {
		retval = append(retval, fmt.Sprintf("%d Bayonets", s.Infantry))
	}
	if s.Cavalry > 0 {
		retval = append(retval, fmt.Sprintf("%d Sabres", s.Cavalry))
	}
	if s.Guns > 0 {
		retval = append(retval, fmt.Sprintf("%d Guns", s.Guns))
	}
	return strings.Join(retval, ", ")
}

// BattleFormation returns the default battle formation
// for a command, based on its drill rating
func (c *Command) BattleFormation() Formation {
	switch c.Arm {
	case Arm_ARTILLERY:
		return Formation_LINE
	case Arm_CAVALRY:
		return Formation_COLUMN
	}
	switch c.Drill {
	case Drill_LINEAR:
		if c.Grade > UnitGrade_REGULAR {
			return Formation_DOUBLE_LINE
		}
		return Formation_LINE
	case Drill_MASSED:
		return Formation_COLUMN
	case Drill_RAPID:
		return Formation_DOUBLE_LINE
	}
	return Formation_DEBANDE
}

func (c *Command) initState(parent *Command, standDown bool, side MapSide, mx, my int32) {
	pos := BattlefieldPosition_OFF_BOARD
	form := Formation_MARCH_COLUMN
	if parent != nil {
		pos = parent.GetGameState().GetPosition()
		form = parent.GetGameState().GetFormation()
	}
	if c.Arrival == nil && parent.GetArrival() != nil {
		c.Arrival = parent.GetArrival()
	}
	pos = c.Arrival.GetPosition()
	switch {
	case c.GetArrival().GetFrom() > 0:
		// offboard units are in march column heading to the battle
		pos = BattlefieldPosition_OFF_BOARD
		form = Formation_MARCH_COLUMN
		c.Arrival.ComputedTurn = int32(rand.Intn(int(c.Arrival.To-c.Arrival.From)) + int(c.Arrival.From))
	case c.GetArrival().GetContact():
		// on-board units in contact are in battle formation
		form = c.BattleFormation()
	case pos == BattlefieldPosition_REAR:
		// units in the rear echelon are in reserve formation
		form = Formation_RESERVE
	}
	if standDown {
		pos = BattlefieldPosition_REAR
		form = Formation_COLUMN
	}
	if form != Formation_MARCH_COLUMN {
		switch c.Arm {
		case Arm_CAVALRY:
			form = Formation_COLUMN
		case Arm_ARTILLERY:
			form = Formation_LINE
		}
	}

	randomizeCentre := func(x int32) int32 {
		allow := x - 4
		if allow < 2 {
			return 2
		}
		return int32(3 + rand.Intn(int(allow)))
	}
	randomizeLow := func(x int32) int32 {
		if x < 5 {
			return 1
		}
		return int32(1 + rand.Intn(2))
	}
	randomizeHigh := func(x int32) int32 {
		if x < 5 {
			return x
		}
		return int32(x - int32(rand.Intn(2)))
	}

	var x, y int32
	switch c.Arrival.Position {
	case BattlefieldPosition_CENTRE, BattlefieldPosition_REAR:
		switch side {
		case MapSide_FRONT:
			x = randomizeCentre(mx)
			y = my
		case MapSide_TOP:
			x = randomizeCentre(mx)
			y = 1
		case MapSide_RIGHT_FLANK:
			x = mx
			y = randomizeCentre(my)
		case MapSide_LEFT_FLANK:
			x = 1
			y = randomizeCentre(my)
		}
	case BattlefieldPosition_RIGHT:
		switch side {
		case MapSide_FRONT:
			x = randomizeHigh(mx)
			y = my
		case MapSide_TOP:
			x = randomizeLow(mx)
			y = 1
		case MapSide_RIGHT_FLANK:
			x = mx
			y = randomizeLow(my)
		case MapSide_LEFT_FLANK:
			x = 1
			y = randomizeHigh(my)
		}
	case BattlefieldPosition_LEFT:
		switch side {
		case MapSide_FRONT:
			x = randomizeLow(mx)
			y = my
		case MapSide_TOP:
			x = randomizeHigh(mx)
			y = 1
		case MapSide_RIGHT_FLANK:
			x = mx
			y = randomizeLow(my)
		case MapSide_LEFT_FLANK:
			x = 1
			y = randomizeHigh(my)
		}
	}
	c.GameState = &CommandGameState{
		Position:  pos,
		Formation: form,
		Grid: &Grid{
			X: x,
			Y: y,
		},
		Can: &UnitAction{
			Order: true,
		},
	}
}
