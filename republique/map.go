package republique

import (
	"context"
	rp "github.com/steveoc64/republique5/proto"
)

// GetMap returns the map data relative to the player with the given token
func (s *Server) GetMap(c context.Context, req *rp.TokenMessage) (*rp.MapData, error) {
	s.RLock()
	defer s.RUnlock()
	player, err := s.Auth(req.Id)
	if err != nil {
		return nil, err
	}
	team := s.game.GetTeam(player)
	if team == nil {
		return nil, errUnauthorised
	}
	return &rp.MapData{
		X:    s.game.TableX,
		Y:    s.game.TableY,
		Data: s.game.TableLayout,
		Side: team.GetSide(),
	}, nil
}
