package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/gui/store"
	rp "github.com/steveoc64/republique5/proto"
)

type commanderOrders struct {
	widget.Box
	panel   *OrdersPanel
	command *rp.Command
	btn     *widget.Button
	order   *widget.Label
	store   *store.Store
}

func newCommanderOrders(panel *OrdersPanel, command *rp.Command, store *store.Store) *commanderOrders {
	vbox := widget.NewVBox()
	o := &commanderOrders{
		Box:     *vbox,
		panel:   panel,
		command: command,
		store:   store,
	}

	orderButton := theme.CheckButtonIcon()
	if command.GameState.GetHas().GetOrder() {
		orderButton = theme.CheckButtonCheckedIcon()
	}
	o.btn = widget.NewButtonWithIcon("  "+command.LabelString(), orderButton, func() {
		o.commanderAction()
	})
	switch command.Rank {
	case rp.Rank_ARMY, rp.Rank_CORPS:
		o.btn.HideShadow = false
		o.btn.Style = widget.PrimaryButton
	}
	o.Append(o.btn)
	orderName := upString(command.GameState.GetOrders().String())
	o.order = widget.NewLabelWithStyle(orderName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	o.Append(o.order)
	store.CommanderMap.AddCommandListener(command, o.Listen)

	return o
}

func (o *commanderOrders) Listen(data fyne.DataItem) {
	cmd, ok := data.(*store.Commander)
	if ok && o != nil {
		o.command = cmd.Data
		o.Show()
	}
}

func (o *commanderOrders) Show() {
	// do the command
	orderButton := theme.CheckButtonIcon()
	if o.command.GameState.GetHas().GetOrder() {
		orderButton = theme.CheckButtonCheckedIcon()
	}
	o.btn.SetIcon(orderButton)

	// do the label
	orderName := upString(o.command.GetGameState().GetOrders().String())
	o.order.SetText(orderName)

	// zap the contents
	o.Children = o.Children[:2]

	// add new contents
	waypoints := o.panel.app.MapData.GetWaypoints(o.command)
	for _, v := range waypoints {
		w := widget.NewLabelWithStyle(v.Path, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
		o.Append(w)
	}
	// paint it all
	widget.Refresh(o)
	o.Box.Show()
}

func (o *commanderOrders) commanderAction() {
	o.panel.app.mapPanel.mapWidget.Select(o.command.Id)
	o.panel.app.Tab(TabMap)
}
