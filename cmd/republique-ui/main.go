package main

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Republique")

	gameCode := widget.NewPasswordEntry()
	teamCode := widget.NewPasswordEntry()
	playerCode := widget.NewPasswordEntry()
	form := &widget.Form{
		OnCancel: func() {
			w.Close()
		},
		OnSubmit: func() {
			fmt.Println("Login Submitted", gameCode.Text, teamCode.Text, playerCode.Text)
		},
	}

	gameCode.OnChanged = func(s string) {
		println("gamecode changed", s)
	}
	form.Append("Game", gameCode)
	form.Append("Team", teamCode)
	form.Append("Player", playerCode)
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Login To Game Server"),
		form,
	))

	w.ShowAndRun()
}
