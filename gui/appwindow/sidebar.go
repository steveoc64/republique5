package appwindow

import (
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
	"os"
)

type SideBar struct {
	VBox    *widget.Box
	session *republique.Session
}

func newSideBar(s *republique.Session) *SideBar {
	h := &SideBar{
		session: s,
	}
	h.VBox = widget.NewVBox(
		widget.NewButtonWithIcon("Briefing ", theme.FolderOpenIcon(), h.briefing),
		widget.NewButtonWithIcon("Map      ", theme.HomeIcon(), h.table),
		widget.NewButtonWithIcon("Orders   ", theme.DocumentCreateIcon(), h.orders),
		widget.NewButtonWithIcon("Units    ", theme.ViewRefreshIcon(), h.units),
		widget.NewButtonWithIcon("Withdraw ", theme.NavigateBackIcon(), h.withdraw),
		widget.NewButtonWithIcon("Surrender", theme.WarningIcon(), h.surrender),
	)
	return h
}

func (s *SideBar) briefing() {
	println("briefing", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) table() {
	println("map", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) orders() {
	println("orders", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) units() {
	println("units", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) withdraw() {
	println("withdraw", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) surrender() {
	println("surrender", s.session.LoginDetails.GetBriefing())
	os.Exit(1)
}
