package main

import (
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/login"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := app.New()
	app.SetIcon(resourceIconJpg)

	servername := "localhost:1815"
	if len(os.Args) == 2 {
		servername = os.Args[1]
	}
	login.Show(app, servername)
	app.Run()
}
