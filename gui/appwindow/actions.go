package appwindow

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// ActionsPanel controls the actions
type ActionsPanel struct {
	app        *App
	Box        *fyne.Container
	actionsBox *widget.Box
	Scroll     *widget.ScrollContainer

	Header  *widget.Label
	Actions []*actionWidget
}

// CanvasObject gets the top level canvas object
func (a *ActionsPanel) CanvasObject() fyne.CanvasObject {
	return a.Scroll
}

// newActionsPanel is a pvt function to get a new ActionsPanel
func newActionsPanel(app *App) *ActionsPanel {
	a := &ActionsPanel{
		app:        app,
		Header:     widget.NewLabel("Actions for: " + strings.Join(app.Commanders, ", ")),
		actionsBox: widget.NewVBox(),
	}
	a.Box = fyne.NewContainerWithLayout(layout.NewBorderLayout(a.Header, nil, nil, nil),
		a.Header,
		a.actionsBox,
	)
	for _, command := range a.app.Commands {
		if command.Arrival.From > 0 {
			continue
		}
		a.actionsBox.Append(newActionWidget(a, command))
		for _, subCommand := range command.Subcommands {
			a.actionsBox.Append(newActionWidget(a, subCommand))
		}
	}
	a.Scroll = widget.NewScrollContainer(a.Box)
	a.Box.Show()
	return a
}
