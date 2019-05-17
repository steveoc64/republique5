package republique

import context "context"

func (s *Server) Echo(c context.Context, req *StringMessage) (*StringMessage, error) {
	s.log.WithField("value", req.Value).Println("Echo gRPC")
	return req, nil
}
