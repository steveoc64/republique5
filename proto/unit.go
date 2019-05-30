package republique

import (
	"fmt"
	"strings"
)

func (unit *Unit) LabelString() string {
	nn := ""
	adds := ""
	if unit.Strength > 1 {
		adds = "s"
	}
	switch unit.Arm {
	case Arm_CAVALRY:
		nn = fmt.Sprintf("%d base%s (%d horse)", unit.Strength, adds, unit.Strength*300)
	case Arm_INFANTRY:
		nn = fmt.Sprintf("%d base%s (%d men)", unit.Strength, adds, unit.Strength*550)
		if unit.SkirmisherMax > 0 {
			nn = nn + fmt.Sprintf(" [%d sk]", unit.SkirmisherMax)
		}
	case Arm_ARTILLERY:
		nn = fmt.Sprintf("%d Bty", unit.Strength)
	}
	return fmt.Sprintf("%s %s %s - %s",
		unit.Name,
		strings.ToLower(unit.Grade.String()),
		strings.Replace(strings.ToLower(unit.UnitType.String()), "_", " ", 1),
		nn,
	)
}
