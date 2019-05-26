package appwindow

import (
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
		widget.NewButtonWithIcon("Actions ", theme.ContentPasteIcon(), a.showActions),
		widget.NewButtonWithIcon("Map      ", theme.ViewFullScreenIcon(), a.showMap),
		widget.NewButtonWithIcon("Orders   ", theme.DocumentCreateIcon(), a.showOrders),
		widget.NewButtonWithIcon("Units    ", theme.InfoIcon(), a.showUnits),
		widget.NewButtonWithIcon("Formation", theme.ContentCopyIcon(), a.showFormations),
		widget.NewButtonWithIcon("Advance", theme.MailSendIcon(), a.showAdvance),
		widget.NewButtonWithIcon("Withdraw ", theme.MailReplyIcon(), a.showWithdraw),
		widget.NewButtonWithIcon("Surrender", theme.CancelIcon(), a.showSurrender),
	)
	return h
}
