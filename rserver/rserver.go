package rserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RServer struct {
	log     *logrus.Logger
	version string
	port    int
	webport int
}

// New returns a new republique server
func New(log *logrus.Logger, version string, port int, webport int) *RServer {
	return &RServer{log, version, port, webport}
}

// Run runs a republique server
func (s *RServer) Run() {
	s.log.WithFields(logrus.Fields{
		"version": s.version,
		"port":    s.port,
	}).Println("Starting RServer")
	s.log.SetFormatter(&logrus.JSONFormatter{})

	// Load DB

	// Setup REST endpoints
	go s.rpcProxy()

	// Load GPPC endpoints
	s.grpcRun()
}

func (s *RServer) grpcRun() {
	endpoint := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterGameServiceServer(grpcServer, s)
	s.log.WithFields(logrus.Fields{
		"port":     s.port,
		"endpoint": endpoint,
	}).Println("Serving gRPC")
	grpcServer.Serve(lis)
}

func (s *RServer) rpcProxy() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	rpcendpoint := fmt.Sprintf(":%d", s.port)
	webendpoint := fmt.Sprintf(":%d", s.webport)
	err := RegisterGameServiceHandlerFromEndpoint(ctx, mux, rpcendpoint, opts)
	if err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{
		"port":     s.webport,
		"endpoint": webendpoint,
	}).Println("Serving REST Proxy")
	return http.ListenAndServe(webendpoint, mux)
}
