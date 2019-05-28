package republique

import (
	"context"
	"github.com/steveoc64/republique5"
)

func (s *Server) GameTime(c context.Context, req *republique5.StringMessage) (*republique5.GameTimeResponse, error) {
	if err := s.Auth(req.Value); err != nil {
		return nil, err
	}
	return &republique5.GameTimeResponse{
		Phase:     s.game.Phase,
		GameTime:  s.game.GameTime,
		StopWatch: s.stopWatch,
	}, nil
}
