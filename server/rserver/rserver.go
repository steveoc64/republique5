package rserver

import "github.com/sirupsen/logrus"

type RServer struct {
	log     *logrus.Logger
	version string
}

func New(log *logrus.Logger, version string) *RServer {
	return &RServer{log, version}
}

func (r *RServer) Run() {
	r.log.WithField("version", r.version).Println("Starting RServer")
	r.log.SetFormatter(&logrus.JSONFormatter{})
}
