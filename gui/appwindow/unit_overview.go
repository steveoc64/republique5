package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

// UnitOverview holds the UI for veiwing units overview
type UnitOverview struct {
	app    *App
	panel  *UnitsPanel
	box    *widget.Box
	scroll *widget.ScrollContainer
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitOverview) CanvasObject() fyne.CanvasObject {
	return u.scroll
}

// newUnitOverview returns a new UnitsPanel, including the UI
func newUnitOverview(app *App, panel *UnitsPanel) *UnitOverview {
	u := &UnitOverview{
		app:   app,
		panel: panel,
		box:   widget.NewVBox(),
	}
	u.scroll = widget.NewScrollContainer(u.box)
	u.scroll.Resize(app.MinSize())
	u.build()
	u.CanvasObject().Show()
	return u
}


// build re-builds the overview to match the app data
func (u *UnitOverview) build() {
	u.box.Children = []fyne.CanvasObject{}
	for _, command := range u.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		u.box.Append(u.newCommanderButton(command, true))
		for _, unit := range command.Units {
			u.box.Append(u.newUnitLabel("   ", unit))
		}
		for _, subCommand := range command.Subcommands {
			u.box.Append(u.newCommanderButton(subCommand, false))
			for _, unit := range subCommand.Units {
				u.box.Append(u.newUnitLabel("      ", unit))
			}
		}
	}
	widget.Renderer(u.box).Layout(u.box.MinSize())
}

// commanderAction is the click handler for each commander button
func (u *UnitOverview) commanderAction(command *rp.Command) {
	println("clicked action for", command.LabelString())
	u.panel.ShowCommand(command)
}

// commanderButton returns a new commanderButton
func (u *UnitOverview) newCommanderButton(command *rp.Command, corps bool) *widget.Button {
	orderButton := theme.RadioButtonCheckedIcon()
	if command.GameState.CanOrder {
		orderButton = theme.RadioButtonIcon()
	}
	if corps {
		b := widget.NewButtonWithIcon(command.LabelString(), orderButton, func() {
			u.commanderAction(command)
		})
		b.Style = widget.PrimaryButton
		return b
	} else {
		b := widget.NewButtonWithIcon("  "+command.LabelString(), orderButton, func() {
			u.commanderAction(command)
		})
		return b
	}
}

// unitAction is the click handler for a unit button
func (u *UnitOverview) unitAction(unit *rp.Unit) {
	println("clicked action for unit", unit.Name)
	u.panel.ShowUnit(unit)
}

// newUnitLabel returns a new unitLabel
func (u *UnitOverview) newUnitLabel(spacer string, unit *rp.Unit) *TapLabel {
	st := fyne.TextStyle{Italic: unit.Arm == rp.Arm_CAVALRY, Bold: unit.Arm == rp.Arm_ARTILLERY, Monospace: unit.Arm == rp.Arm_INFANTRY}
	t := func() {
		u.unitAction(unit)
	}
	return NewTapLabel(spacer+unit.LabelString(), fyne.TextAlignLeading, st, t,t)
}

