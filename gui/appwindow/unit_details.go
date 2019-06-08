package appwindow

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitFieldNames is a slice of the field names on the units display
var UnitFieldNames = []string{
	"Unit ID",
	"Grid",
	"Type",
	"Name",
	"Notes",
	"Strength",
	"Skirmishers",
	"Bn Guns",
	"Drill",
	"Reserve",
}

// UnitDetails holds the UI for veiwing unit details
type UnitDetails struct {
	app    *App
	panel  *UnitsPanel
	box    *widget.Box
	form   *widget.Form
	scroll *widget.ScrollContainer
	fields map[string]*canvas.Text
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
		form:  widget.NewForm(),
		box: widget.NewVBox(
			canvas.NewImageFromResource(resourceLineJpg),
			widget.NewLabelWithStyle("Formation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	}
	u.box.Append(u.form)
	u.scroll = widget.NewScrollContainer(u.box)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}

// newItem creates a new form item with the given label
func (u *UnitDetails) newItem(label string, rgba color.RGBA, style fyne.TextStyle, fontSize int) *widget.FormItem {
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
func (u *UnitDetails) build() {
	u.fields = make(map[string]*canvas.Text)
	s := fyne.TextStyle{}
	green := color.RGBA{140, 240, 180, 1}
	blue := color.RGBA{140, 180, 240, 1}
	u.form.AppendItem(u.newItem("Unit ID", green, fyne.TextStyle{Bold: true}, 48))
	u.form.AppendItem(u.newItem("Grid", blue, s, 18))
	u.form.AppendItem(u.newItem("Type", blue, s, 18))
	u.form.AppendItem(u.newItem("Name", blue, s, 18))
	u.form.AppendItem(u.newItem("Notes", blue, s, 18))
	u.form.AppendItem(u.newItem("Strength", blue, s, 18))
	u.form.AppendItem(u.newItem("Skirmishers", blue, s, 18))
	u.form.AppendItem(u.newItem("Bn Guns", blue, s, 18))
	u.form.AppendItem(u.newItem("Drill", blue, s, 18))
	u.form.AppendItem(u.newItem("Reserve", blue, s, 18))
	u.form.Show()
}

// setField sets the contents of the field given by name
func (u *UnitDetails) setField(name, value string) {
	if t, ok := u.fields[name]; ok {
		t.Text = value
	}
}

// Populate refreshes the UnitDetail fields from the given unit data
func (u *UnitDetails) Populate(unit *rp.Unit) {
	println("populating unit", unit.Name)
	img := u.box.Children[0].(*canvas.Image)
	lbl := u.box.Children[1].(*widget.Label)
	switch unit.GameState.Formation {
	case rp.Formation_LINE:
		img.Resource = resourceLineJpg
		lbl.SetText("Line Formation")
	case rp.Formation_ATTACK_COLUMN, rp.Formation_CLOSED_COLUMN:
		img.Resource = resourceAttackcolumnJpg
		lbl.SetText("Attack Columns")
	case rp.Formation_MARCH_COLUMN:
		img.Resource = resourceMarchcolumnJpg
		lbl.SetText("March Column")
	case rp.Formation_SUPPORTING_LINES:
		img.Resource = resourceSupportingJpg
		lbl.SetText("Supporting Lines")
	case rp.Formation_DEBANDE:
		img.Resource = resourceDebandeJpg
		lbl.SetText("Debande")
	case rp.Formation_ECHELON:
		img.Resource = resourceEchelonJpg
		lbl.SetText("Echelon")

	}
	img.FillMode = canvas.ImageFillOriginal
	canvas.Refresh(img)
	u.setField("Unit ID", fmt.Sprintf("%d", unit.Id))
	u.setField("Grid", fmt.Sprintf("%d,%d",
		unit.GameState.GetGrid().GetX(),
		unit.GameState.GetGrid().GetY(),
	))
	u.setField("Type", fmt.Sprintf("%s %s - %s %s",
		upString(unit.Nationality.String()),
		upString(unit.Arm.String()),
		upString(unit.Grade.String()),
		upString(unit.UnitType.String()),
	))
	u.setField("Name", unit.Name)
	u.setField("Notes", unit.Notes)
	u.setField("Strength", unit.GetStrengthLabel())
	u.setField("Skirmishers", unit.GetSKLabel())
	u.setField("Bn Guns", fmt.Sprintf("%v", unit.BnGuns))
	u.setField("Drill", upString(unit.Drill.String()))
	u.setField("Reserve", fmt.Sprintf("%v", unit.CommandReserve))
}
