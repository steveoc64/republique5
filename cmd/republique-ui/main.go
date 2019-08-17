package main

import (
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/login"
	"github.com/steveoc64/republique5/gui/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := app.New()
	app.SetIcon(resourceIconJpg)

	servername := "localhost:1815"
	if len(os.Args) == 2 {
		servername = os.Args[1]
	}

	store := store.NewStore()
	login.Show(app, servername, store)
	app.Run()
}
