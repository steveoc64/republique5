package main

import (
	"flag"
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
		gamename := ""
		port := 1815
		web := 8080
		flag.IntVar(&port, "port", 1815, "port number to run RPC server on")
		flag.IntVar(&web, "web", 8080, "port number to run web app on")
		flag.StringVar(&gamename, "game", "game.game", "filename of the game to run")
		flag.Parse()

		println("set gamename to", gamename)
		println("set port to", port)
		println("set web to", web)
		serve(log, gamename, port, web)
	default:
		usage()
	}
}
