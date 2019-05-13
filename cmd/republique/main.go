package main

import (
	"flag"
	"github.com/steveoc64/republique5/republique"

	"github.com/sirupsen/logrus"
)

const version = "0.1"

func main() {
	port := flag.Int("port", 1815, "port number to run on")
	webport := flag.Int("webport", 8080, "port number to run on")
	flag.Parse()
	log := logrus.New()
	s := republique.New(log, version, *port, *webport)
	s.Run()
}
