package republique

import (
	"fmt"
	"strings"
)

// LabesString returns a formatted string for rendering the whole command
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
	c.GameState = &CommandGameState{
		Position:  pos,
		Formation: form,
		CanOrder:  true,
	}
}
