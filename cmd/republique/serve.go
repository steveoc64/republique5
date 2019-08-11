package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/republique"
	"os"
)

func serve(log *logrus.Logger, args []string) error {
	f := flag.NewFlagSet("serve args", flag.ContinueOnError)
	gamename := ""
	port := 1815
	web := 0
	f.IntVar(&port, "port", 1815, "port number to run RPC server on")
	f.IntVar(&web, "web", 0, "port number to run web app on")
	f.StringVar(&gamename, "game", "game.game", "filename of the game to run")
	err := f.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if err := info(log, gamename, true, false); err != nil {
		return err
	}
	s, err := republique.NewServer(log, version, gamename, port, web)
	if err != nil {
		return err
	}
	s.Serve()
	return nil
}
