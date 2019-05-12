package rserver

import "github.com/sirupsen/logrus"

type RServer struct {
	log     *logrus.Logger
	version string
}

// New returns a new republique server
func New(log *logrus.Logger, version string) *RServer {
	return &RServer{log, version}
}

// Run runs a republique server
func (r *RServer) Run(port int) {
	println("and here")
	r.log.WithField("version", r.version).Println("Starting RServer")
	r.log.SetFormatter(&logrus.JSONFormatter{})

	// Load DB

	// Setup REST endpoints
}
