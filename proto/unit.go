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

// GetStrengthLabel returns a description of the unit strength
func (u *Unit) GetStrengthLabel() string {
	if u == nil {
		return ""
	}
	nn := ""
	adds := ""
	if u.Strength > 1 {
		adds = "s"
	}
	switch u.Arm {
	case Arm_CAVALRY:
		nn = fmt.Sprintf("%d base%s (%d horse)", u.Strength, adds, u.Strength*300)
	case Arm_INFANTRY:
		nn = fmt.Sprintf("%d base%s (%d men)", u.Strength, adds, u.Strength*550)
		if u.SkirmisherMax > 0 {
			skd := "-"
			if u.GameState.SkirmishersDeployed > 0 {
				skd = fmt.Sprintf("%d", u.GameState.SkirmishersDeployed)
			}
			nn = nn + fmt.Sprintf(" (%s/%d sk)", skd, u.SkirmisherMax)
		}
	case Arm_ARTILLERY:
		nn = fmt.Sprintf("%d Bty", u.Strength)
		if u.GetGameState().GunsDeployed {
			nn = nn + " [Ready to Fire]"
		} else {
			nn = nn + " [Limbered]"
		}
	}
	return nn
}

// GetSKLabel returns a description of the unit skirmisher state
func (u *Unit) GetSKLabel() string {
	if u == nil {
		return ""
	}
	nn := ""
	if u.SkirmisherMax == 0 {
		return "N/A"
	}

	skd := ""
	if u.GameState.SkirmishersDeployed > 0 {
		skd = fmt.Sprintf("(%d deployed)", u.GameState.SkirmishersDeployed)
	}
	nn = nn + fmt.Sprintf("%d %s %s", u.SkirmisherMax, upString(u.SkirmishRating.String()), skd)
	return nn
}

// LabelString returns a formatted string suitable for labelling the unit in the GUI
func (u *Unit) LabelString() string {
	if u == nil {
		return ""
	}
	nn := ""
	ff := ""
	adds := ""
	if u.Strength > 1 {
		adds = "s"
	}
	switch u.Arm {
	case Arm_CAVALRY:
		nn = fmt.Sprintf("%d base%s (%d horse)", u.Strength, adds, u.Strength*300)
		ff = "in " + upString(strings.Replace(strings.ToLower(u.GetGameState().GetFormation().String()), "_", " ", 1))
	case Arm_INFANTRY:
		nn = fmt.Sprintf("%d base%s (%d men)", u.Strength, adds, u.Strength*550)
		if u.SkirmisherMax > 0 {
			skd := "-"
			if u.GameState.SkirmishersDeployed > 0 {
				skd = fmt.Sprintf("%d", u.GameState.SkirmishersDeployed)
			}
			nn = nn + fmt.Sprintf(" (%s/%d sk)", skd, u.SkirmisherMax)
		}
		ff = "in " + upString(strings.Replace(strings.ToLower(u.GetGameState().GetFormation().String()), "_", " ", 1))
	case Arm_ARTILLERY:
		nn = fmt.Sprintf("%d Bty", u.Strength)
		if u.GetGameState().GunsDeployed {
			nn = nn + " [Ready to Fire]"
		} else {
			nn = nn + " [Limbered]"
		}
	}
	return fmt.Sprintf("%s, %s %s - %s %s",
		u.Name,
		strings.ToLower(u.Grade.String()),
		strings.Replace(strings.ToLower(u.UnitType.String()), "_", " ", 1),
		nn,
		ff,
	)
}

// BattleFormation returns the default battle formation
// for a unit, based on its drill
func (u *Unit) BattleFormation() Formation {
	if u.Arm == Arm_CAVALRY {
		switch u.UnitType {
		case UnitType_CAVALRY_LIGHT:
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
		if u.UnitType == UnitType_INFANTRY_LIGHT && u.Grade >= UnitGrade_REGULAR {
			return Formation_LINE
		}
		return Formation_CLOSED_COLUMN
	case Drill_RAPID:
		if u.UnitType == UnitType_INFANTRY_LIGHT {
			if u.Grade > UnitGrade_REGULAR {
				return Formation_DEBANDE
			}
			return Formation_LINE
		}
		return Formation_ATTACK_COLUMN
	}
	return Formation_DEBANDE
}

func (u *Unit) initState(parent *Command, standDown bool) {
	sk := int32(0)
	form := parent.GetGameState().GetFormation()
	if parent.GetArrival().GetContact() {
		form = u.BattleFormation()
		if u.UnitType == UnitType_INFANTRY_LIGHT {
			sk = u.SkirmisherMax
		}
	}
	guns := parent.GetArrival().GetContact()
	if standDown {
		form = Formation_MARCH_COLUMN
		guns = false
		sk = 0
	}
	u.GameState = &UnitGameState{
		SkirmishersDeployed: sk,
		Formation:           form,
		GunsDeployed:        guns,
	}
}
