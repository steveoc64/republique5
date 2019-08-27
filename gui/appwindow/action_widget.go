package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/gui/store"
	rp "github.com/steveoc64/republique5/proto"
)

type actionWidget struct {
	widget.Box
	panel   *ActionsPanel
	command *rp.Command
	btn     *widget.Button
}

func newActionWidget(panel *ActionsPanel, command *rp.Command) *actionWidget {
	vbox := widget.NewVBox()
	a := &actionWidget{
		Box:     *vbox,
		panel:   panel,
		command: command,
	}
	a.btn = widget.NewButtonWithIcon(command.LabelString(), theme.CheckButtonIcon(), func() {
		a.commanderAction()
	})
	switch command.Rank {
	case rp.Rank_ARMY, rp.Rank_CORPS:
		a.btn.HideShadow = false
		a.btn.Style = widget.PrimaryButton
	}
	a.Append(a.btn)
	a.panel.app.store.CommanderMap.AddCommandListener(command, a.Listen)
	return a
}

func (a *actionWidget) commanderAction() {
	a.panel.app.mapPanel.mapWidget.Select(a.command.Id)
	a.panel.app.Tab(TabMap)
}

// Listen is the subscription callback when the data the actionWidget is interested in changes
func (a *actionWidget) Listen(data fyne.DataItem) {
	cmd, ok := data.(*store.Commander)
	if ok && a != nil {
		a.command = cmd.Data
		a.Show()
	}
}

// Show builds / refreshes the actionWidget
func (a *actionWidget) Show() {
	allDone := false

	// zap the contents
	a.Children = a.Children[:1]

	// Does it need an order ?
	if a.command.GetGameState().GetOrders() == rp.Order_NO_ORDERS {
		a.addItem("Currently has No Orders", false, rp.ActionType_COMMENT)
	}

	has := a.command.GetGameState().GetHas()
	if has == nil {
		has = &rp.UnitAction{}
	}
	can := a.command.GetGameState().GetCan()
	if can == nil {
		can = &rp.UnitAction{}
	}
	must := a.command.GetGameState().GetMust()

	// TODO - go through the list of actions attached to this commander, and add widgets for them all

	// tick if all done
	orderButton := theme.CheckButtonIcon()
	if allDone {
		orderButton = theme.CheckButtonCheckedIcon()
	}
	a.btn.SetIcon(orderButton)

	// does it need figure changes on the table ?
	figs := ""
	switch a.command.Arm {
	case rp.Arm_ARTILLERY:
		figs = a.getArtilleryFigures()
	case rp.Arm_CAVALRY:
		figs = a.getCavalryFigures()
	case rp.Arm_INFANTRY:
		figs = a.getInfantryFigures()
	}
	if figs != "" {
		a.addItem(figs, true)
	}

	// draw it all up
	widget.Refresh(a)
	a.Box.Show()
}

func (a *actionWidget) addItem(s string, done bool, t rp.ActionType) {
	w := widget.NewLabelWithStyle(s, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	a.Append(w)
	w.Show()
}

func (a *actionWidget) getArtilleryFigures() string {
	switch a.command.GetGameState().GetOrders() {
	case rp.Order_MARCH, rp.Order_RESTAGE:
		for _, v := range a.command.GetUnits() {
			if v.Arm == rp.Arm_ARTILLERY {
				if v.GetGameState().GunsDeployed {
					return "Limber Up Gun models ready for movement"
				}
			}
		}
	case rp.Order_FIRE, rp.Order_DEFEND:
		for _, v := range a.command.GetUnits() {
			if v.Arm == rp.Arm_ARTILLERY {
				if !v.GetGameState().GunsDeployed {
					return "UnLimber and Deploy Gun models"
				}
			}
		}
	}
	return ""
}

func (a *actionWidget) getCavalryFigures() string {
	return ""
}

func (a *actionWidget) getInfantryFigures() string {
	switch a.command.GetGameState().GetOrders() {
	case rp.Order_RESTAGE, rp.Order_NO_ORDERS, rp.Order_RALLY:
		return ""
	case rp.Order_MARCH:
		if a.command.GetRank() != rp.Rank_CORPS && a.command.GetRank() != rp.Rank_ARMY {
			switch a.command.GetGameState().GetFormation() {
			case rp.Formation_MARCH_COLUMN, rp.Formation_COLUMN:
				// can march in these orders
				return ""
			default:
				// anything else needs to go into march column
				return "Arrange brigades into march column"
			}
		}
	case rp.Order_ENGAGE, rp.Order_DEFEND, rp.Order_ATTACK:
		switch a.command.GetGameState().GetFormation() {
		case rp.Formation_RESERVE, rp.Formation_MARCH_COLUMN:
			return "Arrange brigades into line of battle formation"
		}
	}
	return ""
}
