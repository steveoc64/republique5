package rserver

import (
	context "context"
	"flag"
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
}

// New returns a new republique server
func New(log *logrus.Logger, version string, port int) *RServer {
	return &RServer{log, version, port}
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterGameServiceServer(grpcServer, s)
	s.log.WithField("port", s.port).Println("Serving gRPC")
	grpcServer.Serve(lis)
}

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:1815", "endpoint of YourService")
)

func (s *RServer) rpcProxy() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := RegisterGameServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	s.log.WithField("port", "8080").Println("Starting gRPC Proxy Server")
	return http.ListenAndServe(":8080", mux)
}
