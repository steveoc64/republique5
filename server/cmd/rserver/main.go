package rserver

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/server/rserver"
)

const version = "0.1"

func main() {
	log := logrus.New()
	s := rserver.New(log, version)
	s.Run()
}
