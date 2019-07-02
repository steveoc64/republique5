package appwindow

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// OrdersPanel is the UI for placing commander orders
type OrdersPanel struct {
	app *App
	Box *fyne.Container

	Header *widget.Label
	Notes  *widget.Label
}

// CanvasObject returns the top level UI element for the orders
func (o *OrdersPanel) CanvasObject() fyne.CanvasObject {
	return o.Box
}

func newOrdersPanel(app *App) *OrdersPanel {
	h := &OrdersPanel{
		app:    app,
		Header: widget.NewLabel("Orders for: " + strings.Join(app.Commanders, ", ")),
		Notes:  widget.NewLabel("No orders yet"),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewBorderLayout(h.Header, nil, nil, nil),
		h.Header,
		h.Notes,
	)
	h.Box.Show()

	return h
}
