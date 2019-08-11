package main

import (
	"fmt"
	"github.com/steveoc64/memdebug"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const version = "0.1"

func usage() {
	fmt.Println(`
Usage: republique command

commands:
	compile [filenames]
	serve -port RPCPort -web WebPort -game FileName
	info GameName`)
}

func main() {
	memdebug.GCMode(false)
	rand.Seed(time.Now().UnixNano())
	//memdebug.Profile()
	//defer memdebug.WriteProfile()
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
	case "oob":
		err := oob(log, os.Args[2])
		if err != nil {
			println("Error:", err.Error())
			usage()
			os.Exit(1)
		}
	case "info":
		err := info(log, os.Args[2], true, false)
		if err != nil {
			println("Error:", err.Error())
			usage()
			os.Exit(1)
		}
	case "list":
		err := list(log, os.Args[2:])
		if err != nil {
			println("Error:", err.Error())
			usage()
			os.Exit(1)
		}
	default:
		usage()
	}
}
