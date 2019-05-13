package republique

import context "context"

func (s *RServer) Version(c context.Context, req *EmptyMessage) (*StringMessage, error) {
	s.log.Println("Version gRPC")
	return &StringMessage{
		Value: s.version,
	}, nil
}
