package republique

import (
	"context"
	"github.com/steveoc64/republique5"
)

func (s *Server) Version(c context.Context, req *republique5.EmptyMessage) (*republique5.StringMessage, error) {
	s.log.Println("Version gRPC")
	return &republique5.StringMessage{
		Value: s.version,
	}, nil
}
