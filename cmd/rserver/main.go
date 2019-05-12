package main

import (
	"flag"
	"github.com/steveoc64/republique5/rserver"

	"github.com/sirupsen/logrus"
)

const version = "0.1"

func main() {
	port := flag.Int("port", 1815, "port number to run on")
	log := logrus.New()
	s := rserver.New(log, version, *port)
	s.Run()
}
