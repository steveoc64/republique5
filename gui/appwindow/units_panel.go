package appwindow

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitsPanel holds the UI for veiwing units
type UnitsPanel struct {
	app      *App
	Tabs     *widget.TabContainer
	Overview *UnitOverview
	Command  *UnitCommand
	Details  *UnitDetails
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitsPanel) CanvasObject() fyne.CanvasObject {
	return u.Tabs
}

// newUnitsPanel returns a new UnitsPanel, including the UI
func newUnitsPanel(app *App) *UnitsPanel {
	app.GetUnits()

	u := &UnitsPanel{
		app: app,
	}
	u.Overview = newUnitOverview(app, u)
	u.Command = newUnitCommand(app, u)
	u.Details = newUnitDetails(app, u)

	u.Tabs = widget.NewTabContainer(
		widget.NewTabItem("Overview", u.Overview.CanvasObject()),
		widget.NewTabItem("Command", u.Command.CanvasObject()),
		widget.NewTabItem("Unit", u.Details.CanvasObject()),
	)
	u.CanvasObject().Show()

	return u
}

func (u *UnitsPanel) ShowCommand(command *rp.Command) {
	u.Command.Populate(command)
	u.Tabs.SelectTabIndex(1)
}

func (u *UnitsPanel) ShowUnit(unit *rp.Unit) {
	u.Details.Populate(unit)
	u.Tabs.SelectTabIndex(2)
}

// GetUnits fetches the units from the server
func (a *App) GetUnits() error {
	u, err := a.gameServer.GetUnits(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.Commands = u.Commands
	return nil
}