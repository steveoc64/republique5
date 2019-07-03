package appwindow

import (
	"fmt"
	"strings"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// OrdersPanel is the UI for placing commander orders
type OrdersPanel struct {
	app *App
	Box *fyne.Container

	header *widget.Label
	items  *widget.Box
}

// CanvasObject returns the top level UI element for the orders
func (o *OrdersPanel) CanvasObject() fyne.CanvasObject {
	return o.Box
}

func newOrdersPanel(app *App) *OrdersPanel {
	h := &OrdersPanel{
		app:    app,
		header: widget.NewLabel("Orders for: " + strings.Join(app.Commanders, ", ")),
		items:  widget.NewVBox(),
	}
	h.Box = fyne.NewContainerWithLayout(layout.NewBorderLayout(h.header, nil, nil, nil),
		h.header,
		h.items,
	)
	h.build()
	h.Box.Show()

	return h
}

func (o *OrdersPanel) build() {
	o.items.Children = []fyne.CanvasObject{}
	for _, command := range o.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		o.items.Append(o.newCommanderButton(command))
		for _, subCommand := range command.Subcommands {
			o.items.Append(o.newCommanderButton(subCommand))
		}
	}
	widget.Renderer(o.items).Layout(o.items.MinSize())
}

func (o *OrdersPanel) newCommanderButton(command *rp.Command) *widget.Box {
	println("adding commander", command.Name, command.Rank.String())
	box := widget.NewVBox()
	orderButton := theme.CheckButtonIcon()
	if command.GameState.GetHas().GetOrder() {
		orderButton = theme.CheckButtonCheckedIcon()
	}

	// TODO - if corps, then make it a primary button
	btn := widget.NewButtonWithIcon("  "+command.LabelString(), orderButton, func() {
		o.commanderAction(command)
	})
	if len(command.Subcommands) > 0 {
		btn.Style = widget.PrimaryButton
	}
	box.Append(btn)

	// add the objective paths
	paths := []string{}
	for _, v := range command.GetGameState().GetObjective() {
		paths = append(paths, fmt.Sprintf("%d,%d", v.X, v.Y))
	}
	orderName := upString(command.GetGameState().GetOrders().String())
	box.Append(widget.NewLabel(orderName + ": " + strings.Join(paths, " -> ")))

	return box
}

func (o *OrdersPanel) commanderAction(command *rp.Command) {
	o.app.mapPanel.mapWidget.Select(command.Id)
	o.app.Tab(TAB_MAP)
}
