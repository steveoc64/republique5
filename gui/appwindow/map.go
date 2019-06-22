package appwindow

import (
	"context"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
)

const (
	ORDER_NONE = iota
	ORDER_MARRCH
	ORDER_DEFEND
	ORDER_ATTACK
)

// MapPanel is the UI for the map
type MapPanel struct {
	app       *App
	content   *fyne.Container
	mapWidget *MapWidget
	hbox      *widget.Box
	command   *rp.Command
	order     int
	unitDesc  *widget.Label
	marchBtn  *widget.Button
	defendBtn *widget.Button
	attackBtn *widget.Button
}

// CanvasObject returns the top level UI element for the map
func (m *MapPanel) CanvasObject() fyne.CanvasObject {
	return m.content
}

func (m *MapPanel) SetCommand(cmd *rp.Command) {
	m.command = cmd
	m.setOrder(ORDER_NONE)
	m.unitDesc.SetText(cmd.LongDescription())
}

func (m *MapPanel) unitInfo() {
	if m.command != nil {
		m.app.unitsPanel.ShowCommand(m.command)
		m.app.Tab(TAB_UNITS)
	}
}

func (m *MapPanel) setOrder(o int) {
	if m.order == o {
		m.order = ORDER_NONE
	} else {
		m.order = o
	}
	m.marchBtn.SetIcon(theme.RadioButtonIcon())
	m.defendBtn.SetIcon(theme.RadioButtonIcon())
	m.attackBtn.SetIcon(theme.RadioButtonIcon())
	switch m.order {
	case ORDER_MARRCH:
		m.marchBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_DEFEND:
		m.defendBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_ATTACK:
		m.attackBtn.SetIcon(theme.RadioButtonCheckedIcon())
	}
}
func (m *MapPanel) marchOrder() {
	m.setOrder(ORDER_MARRCH)
}

func (m *MapPanel) defendOrder() {
	m.setOrder(ORDER_DEFEND)
}

func (m *MapPanel) attackOrder() {
	m.setOrder(ORDER_ATTACK)
}

func (m *MapPanel) undoOrder() {
	m.setOrder(ORDER_NONE)
	// TODO - clear the current order for the unit
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func newMapPanel(app *App) *MapPanel {
	if err := app.GetMap(); err != nil {
		println("Failed to get map", err.Error())
		return nil
	}
	m := &MapPanel{
		app: app,
		unitDesc: widget.NewLabelWithStyle(
			"No Unit Selected",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Italic: true},
		),
	}

	m.marchBtn = widget.NewButtonWithIcon("Move", theme.RadioButtonIcon(), m.marchOrder)
	m.defendBtn = widget.NewButtonWithIcon("Defend", theme.RadioButtonIcon(), m.defendOrder)
	m.attackBtn = widget.NewButtonWithIcon("Attack", theme.RadioButtonIcon(), m.attackOrder)
	m.hbox = widget.NewHBox(
		widget.NewButtonWithIcon("Unit", theme.InfoIcon(), m.unitInfo),
		layout.NewSpacer(),
		m.unitDesc,
		layout.NewSpacer(),
		m.marchBtn,
		m.defendBtn,
		m.attackBtn,
		widget.NewButtonWithIcon("Undo", theme.CancelIcon(), m.undoOrder),
	)

	m.mapWidget = newMapWidget(app, app.MapData, m.unitDesc)
	m.content = fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, m.hbox, nil, nil),
		m.mapWidget,
		m.hbox,
	)

	return m
}

// GetMap fetches the map from the server
func (a *App) GetMap() error {
	mapData, err := a.gameServer.GetMap(context.Background(), &a.Token)
	if err != nil {
		return err
	}
	a.MapData = mapData
	return nil
}
