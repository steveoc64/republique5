package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type WithdrawPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newWithdrawPanel(app *App) *WithdrawPanel {
	h := &WithdrawPanel{
		app:    app,
		Header: widget.NewLabelWithStyle("Withdraw", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		Notes:  widget.NewLabel("Really withdraw ?"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
