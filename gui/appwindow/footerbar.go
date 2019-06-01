package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type FooterBar struct {
	app *App
	Box *widget.Box

	onDone     func(bool)
	Done       bool
	PhaseLabel *widget.Label
	DoneBtn    *widget.Button
	StopWatch  *widget.Label
}

func (f *FooterBar) CanvasObject() fyne.CanvasObject {
	return f.Box
}

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

func (f *FooterBar) NotDone() {
	f.Done = false
	f.DoneBtn.SetIcon(theme.CheckButtonIcon())
	f.onDone(false)
}
