package main

import (
	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/appwindow"
	"github.com/steveoc64/republique5/gui/login"
	"github.com/steveoc64/republique5/republique"
	"os"
)

func main() {
	app := app.New()

	servername := "localhost:1815"
	if len(os.Args) == 2 {
		servername = os.Args[1]
	}
	s := &republique.Session{
		ServerName: servername,
	}
	login.Show(s, app, servername, func() {
		println("session", s.String())
		appwindow.Show(s, app)
	})
	app.Run()
}
