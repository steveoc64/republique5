package appwindow

import (
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
	"os"
)

type SideBar struct {
	Box     *widget.Box
	session *republique.Session
}

func newSideBar(s *republique.Session) *SideBar {
	h := &SideBar{
		session: s,
	}
	h.Box = widget.NewVBox(
		widget.NewButtonWithIcon("Briefing ", theme.FolderIcon(), h.briefing),
		widget.NewButtonWithIcon("Map      ", theme.ViewFullScreenIcon(), h.table),
		widget.NewButtonWithIcon("Orders   ", theme.DocumentCreateIcon(), h.orders),
		widget.NewButtonWithIcon("Units    ", theme.InfoIcon(), h.units),
		widget.NewButtonWithIcon("Formation", theme.ContentCopyIcon(), h.formations),
		widget.NewButtonWithIcon("Advance", theme.MailSendIcon(), h.advance),
		widget.NewButtonWithIcon("Withdraw ", theme.MailReplyIcon(), h.withdraw),
		widget.NewButtonWithIcon("Surrender", theme.CancelIcon(), h.surrender),
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
func (s *SideBar) formations() {
	println("formations", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) advance() {
	println("general advance", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) withdraw() {
	println("withdraw", s.session.LoginDetails.GetBriefing())
}
func (s *SideBar) surrender() {
	println("surrender", s.session.LoginDetails.GetBriefing())
	os.Exit(1)
}
