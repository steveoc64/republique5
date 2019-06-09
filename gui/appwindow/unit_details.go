package appwindow

import (
	"fmt"
	"fyne.io/fyne/layout"
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitDetails holds the UI for veiwing unit details
type UnitDetails struct {
	app              *App
	panel            *UnitsPanel
	formationImg     *canvas.Image
	formationLabel   *widget.Label
	box              *widget.Box
	form             *widget.Form
	scroll           *widget.ScrollContainer
	fields           map[string]*canvas.Text
	unit             *rp.Unit
	hasPrev, hasNext bool
	prevBtn, nextBtn *TapIcon
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitDetails) CanvasObject() fyne.CanvasObject {
	return u.scroll
}

// newUnitCommand return a new UnitCommand
func newUnitDetails(app *App, panel *UnitsPanel) *UnitDetails {
	formationImg := canvas.NewImageFromResource(resourceLineJpg)
	formationImg.FillMode = canvas.ImageFillOriginal
	formationLabel := widget.NewLabelWithStyle("Formation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	hbox := widget.NewHBox()
	u := &UnitDetails{
		app:            app,
		panel:          panel,
		form:           widget.NewForm(),
		formationImg:   formationImg,
		formationLabel: formationLabel,
		box: widget.NewVBox(
			//formationImg,
			formationLabel,
			hbox,
		),
	}
	mkbtn := func(res fyne.Resource, f func()) *TapIcon {
		b := NewTapIcon(res, f, f)
		return b
	}
	u.prevBtn = mkbtn(resourcePrevSvg, u.prevUnit)
	u.nextBtn = mkbtn(resourceNextSvg, u.nextUnit)
	hbox.Append(layout.NewSpacer())
	hbox.Append(u.prevBtn)
	hbox.Append(mkbtn(resourceParentSvg, u.parent))
	hbox.Append(u.nextBtn)
	hbox.Append(layout.NewSpacer())
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
	if strings.HasPrefix(label, "_") {
		label = ""
	}
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
	u.form.AppendItem(u.newItem("_UnitID", green, fyne.TextStyle{Bold: true}, 48))
	u.form.AppendItem(u.newItem("Name", blue, s, 18))
	u.form.AppendItem(u.newItem("_Grid", blue, s, 18))
	u.form.AppendItem(u.newItem("_Type", blue, s, 18))
	u.form.AppendItem(u.newItem("Notes", blue, s, 18))
	u.form.AppendItem(u.newItem("Strength", blue, s, 18))
	u.form.AppendItem(u.newItem("Skirmishers", blue, s, 18))
	u.form.AppendItem(u.newItem("Bn Guns", blue, s, 18))
	u.form.AppendItem(u.newItem("Drill", blue, s, 18))
	u.form.AppendItem(u.newItem("Reserve", blue, s, 18))
}

// setField sets the contents of the field given by name
func (u *UnitDetails) setField(name, value string) {
	if t, ok := u.fields[name]; ok {
		t.Text = value
		canvas.Refresh(t)
	}
}

// Populate refreshes the UnitDetail fields from the given unit data
func (u *UnitDetails) Populate(unit *rp.Unit) {
	u.unit = unit
	switch unit.GameState.Formation {
	case rp.Formation_LINE:
		u.formationImg.Resource = resourceLineJpg
		u.formationLabel.SetText("Line Formation")
	case rp.Formation_ATTACK_COLUMN, rp.Formation_CLOSED_COLUMN:
		u.formationImg.Resource = resourceAttackcolumnJpg
		u.formationLabel.SetText("Attack Columns")
	case rp.Formation_MARCH_COLUMN:
		u.formationImg.Resource = resourceMarchcolumnJpg
		u.formationLabel.SetText("March Column")
	case rp.Formation_SUPPORTING_LINES:
		u.formationImg.Resource = resourceSupportingJpg
		u.formationLabel.SetText("Supporting Lines")
	case rp.Formation_DEBANDE:
		u.formationImg.Resource = resourceDebandeJpg
		u.formationLabel.SetText("Debande")
	case rp.Formation_ECHELON:
		u.formationImg.Resource = resourceEchelonJpg
		u.formationLabel.SetText("Echelon")

	}
	//canvas.Refresh(u.formationImg)
	u.setField("_UnitID", fmt.Sprintf("%d", unit.Id))
	u.setField("_Grid", fmt.Sprintf("%d,%d",
		unit.GameState.GetGrid().GetX(),
		unit.GameState.GetGrid().GetY()))
	u.setField("_Type", fmt.Sprintf("%s %s - %s %s",
		upString(unit.Nationality.String()),
		upString(unit.Arm.String()),
		upString(unit.Grade.String()),
		upString(unit.UnitType.String()),
	))
	u.setField("Name", unit.Name)
	u.setField("Notes", unit.Notes)
	u.setField("Strength", unit.GetStrengthLabel(false))
	u.setField("Skirmishers", unit.GetSKLabel())
	u.setField("Bn Guns", fmt.Sprintf("%v", unit.BnGuns))
	u.setField("Drill", upString(unit.Drill.String()))
	u.setField("Reserve", fmt.Sprintf("%v", unit.CommandReserve))

	// calc if has next prev
	c := u.app.GetUnitCommander(u.unit.GetId())
	if c == nil {
		u.hasNext = false
		u.hasPrev = false
	} else {
		for i, v := range c.Units {
			if v.Id == u.unit.Id {
				u.hasPrev = i != 0
				u.hasNext = i < (len(c.Units) - 1)
				break
			}
		}
	}
	if u.hasPrev {
		u.prevBtn.Enable()
	} else {
		u.prevBtn.Disable()
	}
	if u.hasNext {
		u.nextBtn.Enable()
	} else {
		u.nextBtn.Disable()
	}
}

func (u *UnitDetails) nextUnit() {
	c := u.app.GetUnitCommander(u.unit.GetId())
	target := -1
	if c != nil {
		for i, unit := range c.Units {
			if i == target {
				u.Populate(unit)
				return
			}
			if unit.Id == u.unit.Id {
				target = i + 1
			}
		}
	}
}

func (u *UnitDetails) prevUnit() {
	c := u.app.GetUnitCommander(u.unit.GetId())
	target := -1
	if c != nil {
		for i, unit := range c.Units {
			if unit.Id == u.unit.Id {
				target = i - 1
				break
			}
		}
		if target != -1 {
			u.Populate(c.Units[target])
		}
	}
}

func (u *UnitDetails) parent() {
	c := u.app.GetUnitCommander(u.unit.GetId())
	if c != nil {
		u.panel.ShowCommand(c)
	}
}
