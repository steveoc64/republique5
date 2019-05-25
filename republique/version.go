package republique

import (
	"context"
	rp "github.com/steveoc64/republique5/republique/proto"
)

func (s *Server) Version(c context.Context, req *rp.EmptyMessage) (*rp.StringMessage, error) {
	s.log.Println("Version gRPC")
	return &rp.StringMessage{
		Value: s.version,
	}, nil
}
