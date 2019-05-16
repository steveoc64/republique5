package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
)

func serve(log *logrus.Logger, gamename string, port int, web int) {

	s := republique.NewRServer(log, version, gamename, port, web)
	s.Serve()
}
