package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"strings"
)

type OrdersPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

func newOrdersPanel(app *App) *OrdersPanel {
	h := &OrdersPanel{
		app:    app,
		Header: widget.NewLabel("Orders for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel("No orders yet"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		h.Header,
		h.Notes,
	)
	h.Box.Show()
	return h
}
