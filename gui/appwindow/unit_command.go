package appwindow

import (
	"fmt"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

var CommandFieldNames = []string{
	"ID",
	"Grid",
	"Type",
	"Name",
	"Notes",
	"Commander",
	"Drill",
	"Reserve",
	"Can Order",
	"Can Move",
	"Can Rally",
	"Panic State",
}

// UnitOverview holds the UI for veiwing units overview
type UnitCommand struct {
	app    *App
	panel  *UnitsPanel
	box    *widget.Box
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
		box:   widget.NewVBox(),
	}
	u.box.Append(canvas.NewImageFromResource(resourceCmdLineJpg))
	u.box.Append(widget.NewLabelWithStyle("Formation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	u.box.Append(u.form)
	u.scroll = widget.NewScrollContainer(u.box)
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
	img := u.box.Children[0].(*canvas.Image)
	lbl := u.box.Children[1].(*widget.Label)
	switch command.GameState.Formation {
	case rp.Formation_LINE:
		img.Resource = resourceCmdLineJpg
		lbl.SetText("Formed by Lines of Brigades")
	case rp.Formation_DOUBLE_LINE:
		img.Resource = resourceCmdDoubleJpg
		lbl.SetText("Formed by Double Lines of Brigades")
	case rp.Formation_COLUMN:
		img.Resource = resourceCmdColJpg
		lbl.SetText("Formed by Columns of Brigades")
	case rp.Formation_MARCH_COLUMN:
		img.Resource = resourceCmdMcolJpg
		lbl.SetText("In March Column")
	}
	img.FillMode = canvas.ImageFillOriginal
	img.Show()
	u.setField("ID", fmt.Sprintf("%d", command.Id))
	u.setField("Grid", fmt.Sprintf("%d,%d - %s",
		command.GameState.GetGrid().GetX(),
		command.GameState.GetGrid().GetY(),
		upString(command.GameState.Position.String())))
	u.setField("Type", fmt.Sprintf("%s %s %s (%s)",
		upString(command.Nationality.String()),
		upString(command.Arm.String()),
		upString(command.Rank.String()),
		upString(command.Grade.String())))
	u.setField("Name", command.Name)
	u.setField("Notes", command.Notes)
	u.setField("Commander", fmt.Sprintf("%s (+%d)",
		command.CommanderName,
		command.CommanderBonus))
	u.setField("Drill", fmt.Sprintf("%s Drill - %s Command",
		upString(command.Drill.String()),
		upString(command.CommandRating.String())))
	u.setField("Reserve", fmt.Sprintf("%v", command.Reserve))
	u.setField("Can Order", fmt.Sprintf("%v", command.GameState.CanOrder))
	u.setField("Can Move", fmt.Sprintf("%v", command.GameState.CanMove))
	u.setField("Can Rally", fmt.Sprintf("%v", command.GameState.CanRally))
	u.setField("Panic State", fmt.Sprintf("%v", command.GameState.PanicState))
}
