package republique

import (
	"context"
	"fmt"
	"github.com/steveoc64/republique5"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func (s *Server) grpcRun() {
	endpoint := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	republique5.RegisterGameServiceServer(grpcServer, s)
	s.log.WithFields(logrus.Fields{
		"port":     s.port,
		"endpoint": endpoint,
	}).Println("Serving gRPC")
	grpcServer.Serve(lis)
}

func (s *Server) rpcProxy() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	rpcendpoint := fmt.Sprintf(":%d", s.port)
	webendpoint := fmt.Sprintf(":%d", s.web)
	err := republique5.RegisterGameServiceHandlerFromEndpoint(ctx, mux, rpcendpoint, opts)
	if err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{
		"port":     s.web,
		"endpoint": webendpoint,
	}).Println("Serving REST Proxy")
	return http.ListenAndServe(webendpoint, mux)
}
