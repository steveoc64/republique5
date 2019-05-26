package appwindow

import (
	"os"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type SideBar struct {
	Box *widget.Box
	app *App
}

func newSideBar(a *App) *SideBar {
	h := &SideBar{app: a}
	h.Box = widget.NewVBox(
		widget.NewButtonWithIcon("Briefing ", theme.FolderIcon(), a.showBriefing),
		widget.NewButtonWithIcon("Actions ", theme.ContentPasteIcon(), h.actions),
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

func (s *SideBar) table() {
	println("map", s.app.Briefing)
}
func (s *SideBar) orders() {
	println("orders", s.app.Briefing)
}
func (s *SideBar) units() {
	println("units", s.app.Briefing)
}
func (s *SideBar) formations() {
	println("formations", s.app.Briefing)
}
func (s *SideBar) advance() {
	s.app.PlayAudio("command")
}
func (s *SideBar) withdraw() {
	s.app.PlayAudio("infantry")
}
func (s *SideBar) surrender() {
	println("surrender", s.app.Briefing)
	os.Exit(1)
}
func (s *SideBar) actions() {
	s.app.PlayAudio("artillery")
}
