package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const version = "0.1"

func usage() {
	fmt.Println(`
Usage: republique command

commands:
	compile [filenames]
	serve -port RPCPort -web WebPort -game FileName
`)
}

func main() {
	// get the command verb and process
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	log := logrus.New()
	switch os.Args[1] {
	case "compile":
		compile(log, os.Args[2:])
	case "serve":
		err := serve(log, os.Args[2:])
		if err != nil {
			log.WithError(err).Println("Incorrect Server Args", os.Args[2:])
			usage()
			os.Exit(1)
		}
	default:
		usage()
	}
}
