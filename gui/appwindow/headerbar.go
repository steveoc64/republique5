package appwindow

import (
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type HeaderBar struct {
	HBox    *widget.Box
	session *republique.Session

	GameName *widget.Label
	GameTime *widget.Label
}

func newHeaderBar(s *republique.Session) *HeaderBar {
	h := &HeaderBar{
		session:  s,
		GameName: widget.NewLabel(s.GameName),
		GameTime: widget.NewLabel(s.GameTime.Format(republique.DateFormat)),
	}
	h.HBox = widget.NewHBox(h.GameName, h.GameTime)
	return h
}

func (h *HeaderBar) Recalc() *HeaderBar {
	h.GameName.SetText(h.session.GameName)
	h.GameTime.SetText(h.session.GameTime.Format(republique.TimeFormat))
	return h
}
