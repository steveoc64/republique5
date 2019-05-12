package rserver

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5"
)

const version = "0.1"

func main() {
	println("here 1")
	log := logrus.New()
	println("here 1")
	s := republique5.New(log, version)
	println("here 1")
	port := flag.Int("port", 1815, "port number to run on")
	println("here 1")
	s.Run(*port)
}
