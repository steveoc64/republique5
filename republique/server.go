package republique

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	log      *logrus.Logger
	version  string
	filename string
	port     int
	web      int
	game     *Game
	db       *DB
}

// New returns a new republique server
func NewServer(log *logrus.Logger, version string, filename string, port int, web int) (*Server, error) {
	// load the DB
	if !strings.HasSuffix(filename, ".db") {
		filename = filename + ".db"
	}
	db, err := OpenDB(log, filename)
	if err != nil {
		return nil, err
	}
	data := &Game{}
	err = db.Load(data)
	if err != nil {
		return nil, err
	}
	return &Server{
		log:      log,
		version:  version,
		filename: filename,
		port:     port,
		web:      web,
		game:     data,
		db:       db,
	}, nil
}

// Run runs a republique server
func (s *Server) Serve() {
	s.log.WithFields(logrus.Fields{
		"version":  s.version,
		"port":     s.port,
		"web":      s.web,
		"filename": s.filename,
	}).Println("Starting Republique 5.0 Server")
	s.log.SetFormatter(&logrus.JSONFormatter{})

	// Load DB

	// Setup REST endpoints
	go s.rpcProxy()

	// Load GPPC endpoints
	s.grpcRun()
}

func (s *Server) grpcRun() {
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

func (s *Server) rpcProxy() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	rpcendpoint := fmt.Sprintf(":%d", s.port)
	webendpoint := fmt.Sprintf(":%d", s.web)
	err := RegisterGameServiceHandlerFromEndpoint(ctx, mux, rpcendpoint, opts)
	if err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{
		"port":     s.web,
		"endpoint": webendpoint,
	}).Println("Serving REST Proxy")
	return http.ListenAndServe(webendpoint, mux)
}
