package appwindow

import (
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type FooterBar struct {
	Box     *widget.Box
	session *republique.Session

	PhaseLabel *widget.Label
	DoneBtn    *widget.Button
}

func newFooterBar(s *republique.Session) *FooterBar {
	h := &FooterBar{
		session:    s,
		PhaseLabel: widget.NewLabel(s.Phase),
	}
	h.DoneBtn = widget.NewButtonWithIcon("End Turn", theme.CheckButtonCheckedIcon(), h.done)
	h.Box = widget.NewHBox(h.DoneBtn, h.PhaseLabel)
	return h
}

func (h *FooterBar) done() {
	println("End Turn")
}
