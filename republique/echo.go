package republique

import (
	"context"
	rp "github.com/steveoc64/republique5"
)

func (s *Server) Echo(c context.Context, req *rp.StringMessage) (*rp.StringMessage, error) {
	s.log.WithField("value", req.Value).Println("Echo gRPC")
	return req, nil
}
