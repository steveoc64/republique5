package republique

import (
	"fmt"
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

func (c *Command) initState(parent *Command, standDown bool) {
	pos := BattlefieldPosition_OFF_BOARD
	form := Formation_MARCH_COLUMN
	if parent != nil {
		pos = parent.GetGameState().GetPosition()
		form = parent.GetGameState().GetFormation()
	}
	if c.Arrival != nil {
		pos = c.Arrival.GetPosition()
	}
	switch {
	case c.GetArrival().GetFrom() > 0:
		// offboard units are in march column heading to the battle
		pos = BattlefieldPosition_OFF_BOARD
		form = Formation_MARCH_COLUMN
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
	c.GameState = &CommandGameState{
		Position:  pos,
		Formation: form,
		CanOrder:  true,
	}
}
