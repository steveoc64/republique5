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

	header    *widget.Label
	ordersBox *widget.Box
	Scroll    *widget.ScrollContainer
}

// CanvasObject returns the top level UI element for the orders
func (o *OrdersPanel) CanvasObject() fyne.CanvasObject {
	return o.Scroll
}

func newOrdersPanel(app *App) *OrdersPanel {
	o := &OrdersPanel{
		app:       app,
		header:    widget.NewLabel("Orders for: " + strings.Join(app.Commanders, ", ")),
		ordersBox: widget.NewVBox(),
	}
	o.Box = fyne.NewContainerWithLayout(layout.NewBorderLayout(o.header, nil, nil, nil),
		o.header,
		o.ordersBox,
	)

	for _, command := range o.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		o.ordersBox.Append(newCommanderOrders(o, command))
		for _, subCommand := range command.Subcommands {
			o.ordersBox.Append(newCommanderOrders(o, subCommand))
		}
	}
	o.Scroll = widget.NewScrollContainer(o.Box)
	o.Scroll.Show()
	return o
}
