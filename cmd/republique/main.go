package main

import (
	"flag"
	"github.com/steveoc64/republique5/republique"
	"os"

	"github.com/sirupsen/logrus"
)

const version = "0.1"

func main() {
	// get the command verb and process
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "compile":
		compile(os.Args[2:])
	case "serve":
		serve(os.Args[2:])
	}
	port := flag.Int("port", 1815, "port number to run on")
	webport := flag.Int("webport", 8080, "port number to run on")
	flag.Parse()
	log := logrus.New()
	s := republique.New(log, version, *port, *webport)
	s.Run()
}
