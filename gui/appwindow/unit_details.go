package appwindow

import (
	"fmt"
	"fyne.io/fyne/theme"
	"image/color"
	"math/rand"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

var UnitFieldNames = []string{
	"ID",
	"Grid",
	"Type",
	"Name",
	"Notes",
	"Drill",
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
			//widget.NewLabelWithStyle("Formation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
			canvas.NewText("Formation", theme.PrimaryColor()),
		),
	}
	u.box.Append(u.form)
	u.scroll = widget.NewScrollContainer(u.box)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}

func (u *UnitDetails) newItem(label string) *widget.FormItem {
	//e := widget.NewEntry()
	//e.ReadOnly = true
	e := canvas.NewText("", color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 1})
	e.TextSize = 10 + rand.Intn(20)

	u.fields[label] = e
	return &widget.FormItem{
		Text:   label,
		Widget: e,
	}
}

// build re-builds the overview to match the app data
func (u *UnitDetails) build() {
	u.fields = make(map[string]*canvas.Text)
	for _, v := range UnitFieldNames {
		u.form.AppendItem(u.newItem(v))
	}
	u.form.Show()
}

func (u *UnitDetails) setField(name, value string) {
	if e, ok := u.fields[name]; ok {
		//e.SetText(value)
		e.Text = value
	}
}
func (u *UnitDetails) Populate(unit *rp.Unit) {
	println("populating unit", unit.Name)
	img := u.box.Children[0].(*canvas.Image)
	lbl := u.box.Children[1].(*canvas.Text)
	switch unit.GameState.Formation {
	case rp.Formation_LINE:
		img.Resource = resourceLineJpg
		lbl.Text = "Line Formation"
	case rp.Formation_ATTACK_COLUMN, rp.Formation_CLOSED_COLUMN:
		img.Resource = resourceAttackcolumnJpg
		lbl.Text = "Attack Columns"
	case rp.Formation_MARCH_COLUMN:
		img.Resource = resourceMarchcolumnJpg
		lbl.Text = "March Column"
	case rp.Formation_SUPPORTING_LINES:
		img.Resource = resourceSupportingJpg
		lbl.Text = "Supporting Lines"
	case rp.Formation_DEBANDE:
		img.Resource = resourceDebandeJpg
		lbl.Text = "Debande"
	case rp.Formation_ECHELON:
		img.Resource = resourceEchelonJpg
		lbl.Text = "Echelon"

	}
	img.FillMode = canvas.ImageFillOriginal
	canvas.Refresh(img)
	u.setField("ID", fmt.Sprintf("%d", unit.Id))
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
	u.setField("Drill", upString(unit.Drill.String()))
}
