package appwindow

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// CommandFieldNames is a slice of field names for the UnitCommand panel
var CommandFieldNames = []string{
	"Unit ID",
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

// UnitCommand holds the UI for veiwing units overview
type UnitCommand struct {
	app    *App
	panel  *UnitsPanel
	box    *widget.Box
	form   *widget.Form
	scroll *widget.ScrollContainer
	fields map[string]*canvas.Text
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

// newItem creates a new form item
func (u *UnitCommand) newItem(label string, rgba color.RGBA, style fyne.TextStyle, fontSize int) *widget.FormItem {
	t := canvas.NewText(label, rgba)
	t.TextSize = fontSize
	t.TextStyle = style
	u.fields[label] = t
	return &widget.FormItem{
		Text:   label,
		Widget: t,
	}
}

// build re-builds the overview to match the app data
func (u *UnitCommand) build() {
	u.fields = make(map[string]*canvas.Text)
	s := fyne.TextStyle{}
	green := color.RGBA{40, 240, 180, 1}
	blue := color.RGBA{40, 180, 240, 1}
	red := color.RGBA{200, 0, 0, 1}
	u.form.AppendItem(u.newItem("Unit ID", green, fyne.TextStyle{Bold: true}, 48))
	u.form.AppendItem(u.newItem("Grid", blue, s, 18))
	u.form.AppendItem(u.newItem("Type", blue, s, 18))
	u.form.AppendItem(u.newItem("Name", blue, s, 18))
	u.form.AppendItem(u.newItem("Notes", blue, s, 18))
	u.form.AppendItem(u.newItem("Commander", blue, s, 18))
	u.form.AppendItem(u.newItem("Drill", blue, s, 18))
	u.form.AppendItem(u.newItem("Reserve", blue, s, 18))
	u.form.AppendItem(u.newItem("Can Order", blue, s, 18))
	u.form.AppendItem(u.newItem("Can Move", blue, s, 18))
	u.form.AppendItem(u.newItem("Can Rally", blue, s, 18))
	u.form.AppendItem(u.newItem("Panic State", red, s, 18))
	u.form.Show()
}

// setField sets the text of the given field, by name
func (u *UnitCommand) setField(name, value string) {
	if t, ok := u.fields[name]; ok {
		t.Text = value
	}
}

// Populate refreshes the UnitCommand from the given command data
func (u *UnitCommand) Populate(command *rp.Command) {
	img := u.box.Children[0].(*canvas.Image)
	lbl := u.box.Children[1].(*widget.Label)
	switch command.GetGameState().GetFormation() {
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
	u.setField("Unit ID", fmt.Sprintf("%d", command.Id))
	u.setField("Grid", fmt.Sprintf("%d,%d - %s",
		command.GetGameState().GetGrid().GetX(),
		command.GetGameState().GetGrid().GetY(),
		upString(command.GetGameState().GetPosition().String())))
	u.setField("Type", fmt.Sprintf("%s %s %s (%s)",
		upString(command.GetNationality().String()),
		upString(command.GetArm().String()),
		upString(command.GetRank().String()),
		upString(command.GetGrade().String())))
	u.setField("Name", command.GetName())
	u.setField("Notes", command.GetNotes())
	u.setField("Commander", fmt.Sprintf("%s (+%d)",
		command.GetCommanderName(),
		command.GetCommanderBonus()))
	u.setField("Drill", fmt.Sprintf("%s Drill - %s Command",
		upString(command.GetDrill().String()),
		upString(command.GetCommandRating().String())))
	u.setField("Reserve", fmt.Sprintf("%v", command.GetReserve()))
	u.setField("Can Order", fmt.Sprintf("%v", command.GetGameState().GetCanOrder()))
	u.setField("Can Move", fmt.Sprintf("%v", command.GetGameState().GetCanMove()))
	u.setField("Can Rally", fmt.Sprintf("%v", command.GetGameState().GetCanRally()))
	u.setField("Panic State", fmt.Sprintf("%v", command.GetGameState().GetPanicState()))
}
