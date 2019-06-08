package republique

import (
	"context"
	rp "github.com/steveoc64/republique5/proto"
)

// GameTime is an RPC that returns the current game time
func (s *Server) GameTime(c context.Context, req *rp.TokenMessage) (*rp.GameTimeResponse, error) {
	if _, err := s.Auth(req.Id); err != nil {
		return nil, err
	}
	return &rp.GameTimeResponse{
		Phase:     s.game.Phase,
		GameTime:  s.game.GameTime,
		StopWatch: s.stopWatch,
	}, nil
}
