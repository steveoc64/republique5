package appwindow

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

var CommandFieldNames = []string{
	"ID",
	"Nationality",
	"Arm",
	"Rank",
	"Name",
	"Commander",
	"Command Rating",
	"Ability",
	"Grade",
	"Drill",
	"Notes",
	"Reserve",
}

// UnitOverview holds the UI for veiwing units overview
type UnitCommand struct {
	app    *App
	panel  *UnitsPanel
	form   *widget.Form
	scroll *widget.ScrollContainer
	fields map[string]*widget.Entry
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitCommand) CanvasObject() fyne.CanvasObject {
	return u.scroll
}

// newUnitCommand return a new UnitCommand
func newUnitCommand(app *App, panel *UnitsPanel) *UnitCommand {
	u := &UnitCommand{
		app:   app,
		panel: panel,
		form:  widget.NewForm(),
	}
	u.scroll = widget.NewScrollContainer(u.form)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}

func (u *UnitCommand) newItem(label string) *widget.FormItem {
	e := widget.NewEntry()
	e.ReadOnly = true
	u.fields[label] = e
	return &widget.FormItem{
		Text:   label,
		Widget: e,
	}
}

// build re-builds the overview to match the app data
func (u *UnitCommand) build() {
	u.fields = make(map[string]*widget.Entry)
	for _, v := range CommandFieldNames {
		u.form.AppendItem(u.newItem(v))
	}
	u.form.Show()
}

func (u *UnitCommand) setField(name, value string) {
	if e, ok := u.fields[name]; ok {
		e.SetText(value)
	}
}

func (u *UnitCommand) Populate(command *rp.Command) {
	println("populating", command.Name)
	u.setField("ID", fmt.Sprintf("%04d", command.Id))
	u.setField("Nationality", upString(command.Nationality.String()))
	u.setField("Arm", upString(command.Arm.String()))
	u.setField("Rank", upString(command.Rank.String()))
	u.setField("Name", command.Name)
	u.setField("Commander", command.CommanderName)
	u.setField("Command Rating", upString(command.CommandRating.String()))
	u.setField("Ability", fmt.Sprintf("%v", command.CommanderBonus))
	u.setField("Grade", upString(command.Grade.String()))
	u.setField("Drill", upString(command.Drill.String()))
	u.setField("Notes", command.Notes)
	u.setField("Reserve", fmt.Sprintf("%v", command.Reserve))
}
