package appwindow

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/steveoc64/memdebug"

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
	t1 := time.Now()
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
	memdebug.Print(t1, "build orders")
}

func (o *OrdersPanel) newCommanderButton(command *rp.Command) *widget.Box {
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
	orderName := upString(command.GetGameState().GetOrders().String())
	box.Append(widget.NewLabelWithStyle(orderName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	x := command.GetGameState().GetGrid().GetX()
	y := command.GetGameState().GetGrid().GetY()
	for k, v := range command.GetGameState().GetObjective() {
		if k > 0 {
			distance := math.Sqrt(float64((v.X-x)*(v.X-x) + (v.Y-y)*(v.Y-y)))
			speed := 1.0
			switch {
			case command.Rank == rp.Rank_CORPS, command.Rank == rp.Rank_ARMY:
				speed = 2.0
			case command.Arm == rp.Arm_CAVALRY:
				speed = 1.5
			}
			fromGrid := o.app.mapPanel.mapWidget.grid.Value(int32(x-1), int32(y-1))
			switch fromGrid {
			case 't', 'w':
				speed *= 0.8
			case 'T', 'W', 'h', 'r':
				speed *= 0.6
			case 'H':
				speed *= 0.5
			}
			println("fromgrid", fromGrid, speed)
			toGrid := o.app.mapPanel.mapWidget.grid.Value(int32(v.X-1), int32(v.Y-1))
			switch toGrid {
			case 't', 'w':
				speed *= 0.7
			case 'T', 'W', 'h':
				speed *= 0.5
			case 'H', 'r':
				speed *= 0.4
			}
			going := "at a good march"
			switch {
			case speed >= 1.5:
				going = "with great speed"
			case speed <= 0.4:
				going = "very slow"
			case speed <= 0.5:
				going = "harsh terrain"
			case speed <= 0.6:
				going = "slow going"
			case speed <= 0.7:
				going = "with some delays"
			case speed <= 0.8:
				going = "with minor delays"
			}
			println("toGrid", toGrid, speed)
			elapsed := ((distance / (speed * 3.0)) * 60.0) // 20mins to the mile in good order
			turns := int(elapsed) / 20
			path := fmt.Sprintf("-> %d,%d  (%0.1f miles %s, about %d turns)", v.X, v.Y, distance, going, turns)
			box.Append(widget.NewLabelWithStyle(path, fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
			x = v.X
			y = v.Y
		}
	}

	return box
}

func (o *OrdersPanel) commanderAction(command *rp.Command) {
	o.app.mapPanel.mapWidget.Select(command.Id)
	o.app.Tab(TAB_MAP)
}
