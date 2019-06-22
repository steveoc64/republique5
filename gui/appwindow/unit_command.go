package appwindow

import (
	"fmt"
	"fyne.io/fyne/theme"
	"image/color"
	"strings"

	"fyne.io/fyne/layout"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitCommand holds the UI for veiwing units overview
type UnitCommand struct {
	app            *App
	panel          *UnitsPanel
	box            *widget.Box
	form           *widget.Form
	scroll         *widget.ScrollContainer
	fields         map[string]*canvas.Text
	formationImg   *canvas.Image
	formationLabel *widget.Label
	unitList       *widget.Box
	command        *rp.Command
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitCommand) CanvasObject() fyne.CanvasObject {
	return u.scroll
}

// newUnitCommand return a new UnitCommand
func newUnitCommand(app *App, panel *UnitsPanel) *UnitCommand {
	u := &UnitCommand{
		app:            app,
		panel:          panel,
		form:           widget.NewForm(),
		box:            widget.NewVBox(),
		formationImg:   canvas.NewImageFromResource(resourceCmdLineJpg),
		formationLabel: widget.NewLabelWithStyle("Formation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		unitList:       widget.NewVBox(),
	}
	u.formationImg.FillMode = canvas.ImageFillOriginal
	//u.box.Append(u.formationImg)
	u.box.Append(u.formationLabel)
	u.box.Append(fyne.NewContainerWithLayout(layout.NewGridLayout(2), u.form, u.unitList))
	u.scroll = widget.NewScrollContainer(u.box)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}

func (u *UnitCommand) gotoMap() {
	println("goto map", u.command.Name, u.command.Id)
	u.app.mapPanel.mapWidget.Select(u.command.Id)
	u.app.Tab(TAB_MAP)
}

// newItem creates a new form item
func (u *UnitCommand) newItem(label string, rgba color.RGBA, style fyne.TextStyle, fontSize int) *widget.FormItem {
	t := canvas.NewText(label, rgba)
	t.TextSize = fontSize
	t.TextStyle = style
	u.fields[label] = t
	if strings.HasPrefix(label, "_") {
		label = ""
	}
	return &widget.FormItem{
		Text:   label,
		Widget: t,
	}
}

// build re-builds the overview to match the app data
func (u *UnitCommand) build() {
	u.fields = make(map[string]*canvas.Text)
	s := fyne.TextStyle{}
	u.form.AppendItem(u.newItem("_UnitID", command_green, fyne.TextStyle{Bold: true}, 48))
	u.form.AppendItem(u.newItem("Commander", command_blue, s, 18))
	u.form.AppendItem(u.newItem("_Grid", command_blue, s, 18))
	u.form.AppendItem(u.newItem("_Type", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Name", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Notes", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Strength", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Drill", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Reserve", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Can Order", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Can Move", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Can Rally", command_blue, s, 18))
	u.form.AppendItem(u.newItem("Panic State", command_red, s, 18))
	u.form.Append("", widget.NewButtonWithIcon("View on Map", theme.ViewFullScreenIcon(), u.gotoMap))
}

// setField sets the text of the given field, by name
func (u *UnitCommand) setField(name, value string) {
	if t, ok := u.fields[name]; ok {
		t.Text = value
		canvas.Refresh(t)
	}
}

// Populate refreshes the UnitCommand from the given command data
func (u *UnitCommand) Populate(command *rp.Command) {
	switch command.GetGameState().GetFormation() {
	case rp.Formation_LINE:
		u.formationImg.Resource = resourceCmdLineJpg
		u.formationLabel.SetText("Formed by Lines of Brigades")
	case rp.Formation_DOUBLE_LINE:
		u.formationImg.Resource = resourceCmdDoubleJpg
		u.formationLabel.SetText("Formed by Double Lines of Brigades")
	case rp.Formation_COLUMN:
		u.formationImg.Resource = resourceCmdColJpg
		u.formationLabel.SetText("Formed by Columns of Brigades")
	case rp.Formation_MARCH_COLUMN:
		u.formationImg.Resource = resourceCmdMcolJpg
		u.formationLabel.SetText("In March Column")
	}
	//canvas.Refresh(u.formationImg)
	u.setField("_UnitID", fmt.Sprintf("%d", command.Id))
	u.setField("_Grid", fmt.Sprintf("%d,%d - %s",
		command.GetGameState().GetGrid().GetX(),
		command.GetGameState().GetGrid().GetY(),
		upString(command.GetGameState().GetPosition().String())))
	u.setField("_Type", fmt.Sprintf("%s %s %s (%s)",
		upString(command.GetNationality().String()),
		upString(command.GetArm().String()),
		upString(command.GetRank().String()),
		upString(command.GetGrade().String())))
	u.setField("Name", command.GetName())
	u.setField("Notes", command.GetNotes())
	u.setField("Commander", fmt.Sprintf("%s (+%d)",
		command.GetCommanderName(),
		command.GetCommanderBonus()))
	u.setField("Strength", command.GetCommandStrengthLabel())
	u.setField("Drill", fmt.Sprintf("%s Drill - %s Command",
		upString(command.GetDrill().String()),
		upString(command.GetCommandRating().String())))
	u.setField("Reserve", fmt.Sprintf("%v", command.GetReserve()))
	u.setField("Can Order", fmt.Sprintf("%v", command.GetGameState().GetCan().GetOrder()))
	u.setField("Can Move", fmt.Sprintf("%v", command.GetGameState().GetCan().GetMove()))
	u.setField("Can Rally", fmt.Sprintf("%v", command.GetGameState().GetCan().GetRally()))
	u.setField("Panic State", fmt.Sprintf("%v", command.GetGameState().GetPanicState()))

	u.command = command
	u.populateUnits()
}

func (u *UnitCommand) populateUnits() {
	u.unitList.Children = []fyne.CanvasObject{}
	for _, unit := range u.command.GetUnits() {
		u.unitList.Append(u.newUnitLabel(unit))
	}
	widget.Refresh(u.unitList)
}

func (u *UnitCommand) newUnitLabel(unit *rp.Unit) *TapLabel {
	st := fyne.TextStyle{Italic: unit.Arm == rp.Arm_CAVALRY, Bold: unit.Arm == rp.Arm_ARTILLERY, Monospace: unit.Arm == rp.Arm_INFANTRY}
	t := func() {
		u.panel.ShowUnit(unit)
	}
	lbl := NewTapLabel(unit.ShortLabel(), fyne.TextAlignLeading, st, t, t)
	return lbl
}
