package republique

import (
	"context"
	rp "github.com/steveoc64/republique5/proto"
)

func (s *Server) GameTime(c context.Context, req *rp.StringMessage) (*rp.GameTimeResponse, error) {
	if err := s.Auth(req.Value); err != nil {
		return nil, err
	}
	return &rp.GameTimeResponse{
		Phase:     s.game.Phase,
		GameTime:  s.game.GameTime,
		StopWatch: s.stopWatch,
	}, nil
}
