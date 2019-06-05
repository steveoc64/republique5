package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitOverview holds the UI for veiwing units overview
type UnitDetails struct {
	app    *App
	panel  *UnitsPanel
	box    *fyne.Container
	scroll *widget.ScrollContainer
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitDetails) CanvasObject() fyne.CanvasObject {
	return u.scroll
}

// newUnitCommand return a new UnitCommand
func newUnitDetails(app *App, panel *UnitsPanel) *UnitDetails {
	u := &UnitDetails{
		app:   app,
		panel: panel,
		box:   fyne.NewContainerWithLayout(layout.NewGridLayout(2)),
	}
	u.scroll = widget.NewScrollContainer(u.box)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}

func (u *UnitDetails) newItem(label string) *widget.Label {
	return widget.NewLabel(label)
}

// build re-builds the overview to match the app data
func (u *UnitDetails) build() {
	items := []fyne.CanvasObject{}
	items = append(items,
		u.newItem("Name"),
		u.newItem("Commander"),
		u.newItem("Command Rating"),
		u.newItem("Rank"),
		u.newItem("Ability"),
		u.newItem("Arm"),
		u.newItem("Nationality"),
		u.newItem("Grade"),
		u.newItem("Drill"),
		u.newItem("Notes"),
		u.newItem("Reserve"),
		u.newItem("ID"))
	// TODO - add gamestate items in here

	u.box.Objects = items
}

func (u *UnitDetails) Populate(unit *rp.Unit) {
	println("populating unit", unit.Name)
}
