package appwindow

import (
	"context"

	"fyne.io/fyne/theme"
	rp "github.com/steveoc64/republique5/proto"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// UnitsPanel holds the UI for veiwing units
type UnitsPanel struct {
	app            *App
	Tabs           *widget.TabContainer
	OverviewBox    *widget.Box
	OverviewScroll *widget.ScrollContainer
	CommandBox     *widget.Box
	CommandScroll  *widget.ScrollContainer
	UnitBox        *widget.Box
	UnitScroll     *widget.ScrollContainer
}

// CanvasObject returns the top level widget in the UnitsPanel
func (u *UnitsPanel) CanvasObject() fyne.CanvasObject {
	return u.Tabs
}

// newUnitsPanel returns a new UnitsPanel, including the UI
func newUnitsPanel(app *App) *UnitsPanel {
	u := &UnitsPanel{
		app:         app,
		OverviewBox: widget.NewVBox(),
		CommandBox:  widget.NewVBox(),
		UnitBox:     widget.NewVBox(),
	}
	u.OverviewScroll = widget.NewScrollContainer(u.OverviewBox)
	u.OverviewScroll.Resize(app.MinSize())
	u.CommandScroll = widget.NewScrollContainer(u.CommandBox)
	u.CommandScroll.Resize(app.MinSize())
	u.UnitScroll = widget.NewScrollContainer(u.UnitBox)
	u.UnitScroll.Resize(app.MinSize())

	u.Tabs = widget.NewTabContainer(
		widget.NewTabItem("Overview", u.OverviewScroll),
		widget.NewTabItem("Command", u.CommandBox),
		widget.NewTabItem("Unit", u.UnitBox),
	)

	app.GetUnits()
	u.BuildOverview()
	return u
}

// BuildOverview creates the label and buttons to make up the overview
func (u *UnitsPanel) BuildOverview() {
	u.OverviewBox.Children = []fyne.CanvasObject{}
	for _, command := range u.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		u.OverviewBox.Append(u.newCommanderButton(command, true))
		for _, unit := range command.Units {
			u.OverviewBox.Append(u.newUnitLabel("   ", unit))
		}
		for _, subCommand := range command.Subcommands {
			u.OverviewBox.Append(u.newCommanderButton(subCommand, false))
			for _, unit := range subCommand.Units {
				u.OverviewBox.Append(u.newUnitLabel("      ", unit))
			}
		}
	}
	u.OverviewBox.Show()
}

// commanderAction is the click handler for each commander button
func (u *UnitsPanel) commanderAction(command *rp.Command) {
	println("clicked action for", command.LabelString())
	u.Tabs.SelectTabIndex(1)
}

// commanderButton returns a new commanderButton
func (u *UnitsPanel) newCommanderButton(command *rp.Command, corps bool) *widget.Button {
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

type UnitLabel struct {
	widget.Label
	OnTapped func()
}

func (ul *UnitLabel) Hidex() {
	println("getting a hide call")
}

// Tapped handler for each unitlabel
func (ul *UnitLabel) Tapped(*fyne.PointEvent) {
	println("Clicked on a unit label")
	if ul.OnTapped != nil {
		ul.OnTapped()
	}
}

func (ul *UnitLabel) TappedSecondary(*fyne.PointEvent) {
	println("tapped secondary")
}

func (ul *UnitLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(&ul.Label)
}

func newUnitLabel(text string, alignment fyne.TextAlign, style fyne.TextStyle, tapped func()) *UnitLabel {
	return &UnitLabel{
		widget.Label{
			Text:      text,
			Alignment: alignment,
			TextStyle: style,
		},
		tapped,
	}
}

// unitAction is the click handler for a unit button
func (u *UnitsPanel) unitAction(unit *rp.Unit) {
	println("clicked action for unit", unit.Name)
	u.Tabs.SelectTabIndex(2)
}

// newUnitLabel returns a new unitLabel
func (u *UnitsPanel) newUnitLabel(spacer string, unit *rp.Unit) *UnitLabel {
	st := fyne.TextStyle{Italic: unit.Arm == rp.Arm_CAVALRY, Bold: unit.Arm == rp.Arm_ARTILLERY, Monospace: unit.Arm == rp.Arm_INFANTRY}
	return newUnitLabel(spacer+unit.LabelString(), fyne.TextAlignLeading, st, func() {
		u.unitAction(unit)
	})
}

// GetUnits fetches the units from the server
func (a *App) GetUnits() error {
	u, err := a.gameServer.GetUnits(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.Commands = u.Commands
	return nil
}
