package main

import (
	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/login"
	"os"
)

func main() {
	app := app.New()

	servername := "localhost:1815"
	if len(os.Args) == 2 {
		servername = os.Args[1]
	}
	login.Show(app, servername)
	app.Run()
}
