package republique

import (
	"context"
	"github.com/steveoc64/memdebug"
	"time"

	rp "github.com/steveoc64/republique5/proto"
)

// GetUnits returns all the units owned by the player with the given token
func (s *Server) GetUnits(c context.Context, req *rp.TokenMessage) (*rp.Units, error) {
	s.RLock()
	defer s.RUnlock()
	player, err := s.Auth(req.Id)
	if err != nil {
		return nil, err
	}
	t1 := time.Now()
	team := s.game.GetTeam(player)
	s.log.WithField("token", req.Id).Println("GetUnits gRPC")
	commands := &rp.Units{}
	for _, commander := range player.GetCommanders() {
		c := team.GetCommandByCommanderName(commander)
		if c != nil {
			if c.Arrival.ComputedTurn <= s.game.TurnNumber {
				commands.Commands = append(commands.Commands, c)
			}
		}
	}
	memdebug.Print(t1, "Got units")
	return commands, nil
}
