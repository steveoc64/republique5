package appwindow

import (
	"context"
	"strings"

	"fyne.io/fyne/theme"
	rp "github.com/steveoc64/republique5/proto"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type UnitsPanel struct {
	app *App
	Box *widget.Box

	Header *widget.Label
}

func newUnitsPanel(app *App) *UnitsPanel {
	h := &UnitsPanel{
		app:    app,
		Header: widget.NewLabel("Units for: " + strings.Join(app.Commanders, ", ")),
	}
	h.Box = widget.NewVBox()
	app.GetUnits()
	h.BuildUnits()
	return h
}

func (u *UnitsPanel) commanderAction(command *rp.Command) {
	println("clicked action for", command.LabelString())
}

func (u *UnitsPanel) commanderButton(command *rp.Command, corps bool) *widget.Button {
	if corps {
		b := widget.NewButtonWithIcon(command.LabelString(), theme.RadioButtonCheckedIcon(), func() {
			u.commanderAction(command)
		})
		b.Style = widget.PrimaryButton
		return b
	} else {
		b := widget.NewButtonWithIcon("  "+command.LabelString(), theme.RadioButtonIcon(), func() {
			u.commanderAction(command)
		})
		return b
	}
	return nil
}

func (u *UnitsPanel) unitLabel(spacer string, unit *rp.Unit) *widget.Label {
	st := fyne.TextStyle{Italic: unit.Arm == rp.Arm_CAVALRY, Bold: unit.Arm == rp.Arm_ARTILLERY, Monospace: unit.Arm == rp.Arm_INFANTRY}
	w := widget.NewLabelWithStyle(spacer+unit.LabelString(), fyne.TextAlignLeading, st)
	return w
}

func (u *UnitsPanel) BuildUnits() {
	u.Box.Children = []fyne.CanvasObject{}
	u.Box.Append(u.Header)
	for _, command := range u.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		u.Box.Append(u.commanderButton(command, true))
		for _, unit := range command.Units {
			u.Box.Append(u.unitLabel("   ", unit))
		}
		for _, subCommand := range command.Subcommands {
			u.Box.Append(u.commanderButton(subCommand, false))
			for _, unit := range subCommand.Units {
				u.Box.Append(u.unitLabel("      ", unit))
			}
		}
	}
	u.Box.Show()
}

func (a *App) GetUnits() error {
	u, err := a.client.GetUnits(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.Commands = u.Commands
	return nil
}
