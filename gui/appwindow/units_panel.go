package appwindow

import (
	"context"
	"time"

	"github.com/steveoc64/memdebug"

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
	dark     bool
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitsPanel) CanvasObject() fyne.CanvasObject {
	return u.Tabs
}

// newUnitsPanel returns a new UnitsPanel, including the UI
func newUnitsPanel(app *App) *UnitsPanel {
	u := &UnitsPanel{
		app:  app,
		dark: app.isDarkTheme,
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

// ShowCommand navigates to the commander details and populates the commander
func (u *UnitsPanel) ShowCommand(command *rp.Command) {
	u.Tabs.SelectTabIndex(1)
	u.Command.Populate(command)
}

// ShowUnit navigates to the unit details and populates the unit
func (u *UnitsPanel) ShowUnit(unit *rp.Unit) {
	u.Tabs.SelectTabIndex(2)
	u.Details.Populate(unit)
}

func (u *UnitsPanel) darkTheme() {
	u.dark = true
	u.Details.darkTheme()
	u.Command.darkTheme()
}

func (u *UnitsPanel) lightTheme() {
	u.dark = false
	u.Details.lightTheme()
	u.Command.lightTheme()
}

// GetUnits fetches the units from the server
func (a *App) GetUnits() error {
	memdebug.Print(time.Now(), "fetching units")
	u, err := a.gameServer.GetUnits(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.Commands = u.Commands
	a.store.CommanderMap.Load(u.Commands)

	u, err = a.gameServer.GetEnemy(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.Enemy = u.Commands

	return nil
}
