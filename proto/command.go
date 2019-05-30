package republique

import (
	"fmt"
	"strings"
)

func (c *Command) LabelString() string {
	s := fmt.Sprintf("%s - %s [%s]", c.Name, c.CommanderName, strings.ToLower(c.Arrival.GetPosition().String()))
	if c.Notes != "" {
		s = s + " (" + c.Notes + ")"
	}
	return s
}
