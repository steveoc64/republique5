package main

import (
	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/login"
)

func main() {
	app := app.New()

	login.Show(app)
	app.Run()
}
