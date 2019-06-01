package republique

import (
	"fmt"
	"strings"
)

func upString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i < 1 {
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		}
		b.WriteByte(c)
	}
	return b.String()
}

// LabelString returns a formatted string suitable for labelling the unit in the GUI
func (unit *Unit) LabelString() string {
	nn := ""
	ff := ""
	adds := ""
	if unit.Strength > 1 {
		adds = "s"
	}
	switch unit.Arm {
	case Arm_CAVALRY:
		nn = fmt.Sprintf("%d base%s (%d horse)", unit.Strength, adds, unit.Strength*300)
		ff = "in " + upString(strings.Replace(strings.ToLower(unit.GetGameState().GetFormation().String()), "_", " ", 1))
	case Arm_INFANTRY:
		nn = fmt.Sprintf("%d base%s (%d men)", unit.Strength, adds, unit.Strength*550)
		if unit.SkirmisherMax > 0 {
			nn = nn + fmt.Sprintf(" (%d sk)", unit.SkirmisherMax)
		}
		ff = "in " + upString(strings.Replace(strings.ToLower(unit.GetGameState().GetFormation().String()), "_", " ", 1))
	case Arm_ARTILLERY:
		nn = fmt.Sprintf("%d Bty", unit.Strength)
		if unit.GetGameState().GunsDeployed {
			nn = nn + " [Ready to Fire]"
		} else {
			nn = nn + " [Limbered]"
		}
	}
	return fmt.Sprintf("%s %s %s - %s %s",
		unit.Name,
		strings.ToLower(unit.Grade.String()),
		strings.Replace(strings.ToLower(unit.UnitType.String()), "_", " ", 1),
		nn,
		ff,
	)
}

// BattleFormation returns the default battle formation
// for a unit, based on its drill
func (u *Unit) BattleFormation() Formation {
	if u.Arm == Arm_CAVALRY {
		switch u.UnitType {
		case UnitType_INFANTRY_LIGHT:
			return Formation_DEBANDE
		case UnitType_CAVALRY_HUSSAR:
			return Formation_ATTACK_COLUMN
		}
		return Formation_RESERVE
	}
	if u.Arm == Arm_ARTILLERY {
		return Formation_LINE
	}
	switch u.Drill {
	case Drill_LINEAR:
		return Formation_LINE
	case Drill_MASSED:
		return Formation_CLOSED_COLUMN
	case Drill_RAPID:
		if u.UnitType == UnitType_INFANTRY_LIGHT {
			if u.Grade > UnitGrade_REGULAR {
				return Formation_DEBANDE
			}
			return Formation_DOUBLE_LINE
		}
		return Formation_ATTACK_COLUMN
	}
	return Formation_DEBANDE
}

func (u *Unit) initState(parent *Command, standDown bool) {
	sk := int32(0)
	if parent.Arrival.Contact {
		sk = u.SkirmisherMax
	}

	form := parent.GetGameState().GetFormation()
	if parent.GetArrival().GetContact() {
		form = u.BattleFormation()
	}
	guns := parent.GetArrival().GetContact()
	if standDown {
		form = Formation_MARCH_COLUMN
		guns = false
	}
	u.GameState = &UnitGameState{
		SkirmishersDeployed: sk,
		Formation:           form,
		GunsDeployed:        guns,
	}
}
