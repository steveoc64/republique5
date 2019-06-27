package main

import (
	"math/rand"
	"time"

	"fyne.io/fyne/app"
	"github.com/steveoc64/republique5/gui/mapeditor"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	a := app.New()
	mapeditor.New(a, 6, 4, "")
	a.Run()
}
