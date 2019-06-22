package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// FooterBar is the UI for the footer
type FooterBar struct {
	app *App
	Box *widget.Box

	onDone     func(bool)
	Done       bool
	PhaseLabel *widget.Label
	DoneBtn    *widget.Button
	StopWatch  *widget.Label
}

// CanvasObject returns the top level UI element in the footer
func (f *FooterBar) CanvasObject() fyne.CanvasObject {
	return f.Box
}

// newFooterBar creates a new FooterBar and returns it
func newFooterBar(app *App, onDone func(bool)) *FooterBar {
	h := &FooterBar{
		app:        app,
		onDone:     onDone,
		PhaseLabel: widget.NewLabel(app.Phase),
		StopWatch:  widget.NewLabel("00:00"),
	}
	h.DoneBtn = widget.NewButtonWithIcon("End Turn", theme.CheckButtonIcon(), h.ToggleDone)
	h.DoneBtn.Style = widget.PrimaryButton
	h.Box = widget.NewHBox(
		h.PhaseLabel,
		layout.NewSpacer(),
		h.DoneBtn,
		layout.NewSpacer(),
		h.StopWatch,
	)
	return h
}

// ToggleDone toggles the Done flag to signal that the player is ready for the next phase
func (f *FooterBar) ToggleDone() {
	f.Done = !f.Done
	switch f.Done {
	case true:
		f.DoneBtn.SetIcon(theme.CheckButtonCheckedIcon())
	case false:
		f.DoneBtn.SetIcon(theme.CheckButtonIcon())
	}
	f.onDone(f.Done)
}

// NotDone sets the done flag back to not done. Use this whenever the player does an action
// that un-does the fact that they are finished with the turn
func (f *FooterBar) NotDone() {
	f.Done = false
	f.DoneBtn.SetIcon(theme.CheckButtonIcon())
	f.onDone(false)
}
