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
	ORDER_ENGAGE
	ORDER_ATTACK
	ORDER_FIRE
	ORDER_CHARGE
	ORDER_RALLY
	ORDER_PURSUIT
)

// MapPanel is the UI for the map
type MapPanel struct {
	app          *App
	content      *fyne.Container
	mapWidget    *MapWidget
	vbox         *widget.Box
	hbox1        *widget.Box
	hbox2        *widget.Box
	borderLayout fyne.Layout
	command      *rp.Command
	order        int
	unitDesc     *widget.Label
	marchBtn     *widget.Button
	defendBtn    *widget.Button
	attackBtn    *widget.Button
	engageBtn    *widget.Button
	fireBtn      *widget.Button
	chargeBtn    *widget.Button
	rallyBtn     *widget.Button
	pursuitBtn   *widget.Button
}

// CanvasObject returns the top level UI element for the map
func (m *MapPanel) CanvasObject() fyne.CanvasObject {
	return m.content
}

func (m *MapPanel) SetCommand(cmd *rp.Command) {
	m.command = cmd
	m.setOrder(ORDER_NONE)
	if cmd == nil {
		m.mapWidget.grid.Select(0)
		widget.Renderer(m.mapWidget).Refresh()
		m.hbox2.Hide()
		m.unitDesc.SetText("")
		return
	}
	m.hbox2.Show()
	m.unitDesc.SetText(cmd.LongDescription())
	switch cmd.GetRank() {
	case rp.Rank_CORPS, rp.Rank_ARMY:
		m.marchBtn.Show()
		m.defendBtn.Hide()
		m.attackBtn.Hide()
		m.engageBtn.Hide()
		m.fireBtn.Hide()
		m.chargeBtn.Hide()
		m.rallyBtn.Show()
		m.pursuitBtn.Hide()
	default:
		switch cmd.GetArm() {
		case rp.Arm_INFANTRY:
			m.marchBtn.Show()
			m.defendBtn.Show()
			m.attackBtn.Show()
			m.engageBtn.Show()
			m.fireBtn.Hide()
			m.chargeBtn.Hide()
			m.rallyBtn.Hide()
			m.pursuitBtn.Hide()
		case rp.Arm_CAVALRY:
			m.marchBtn.Show()
			m.defendBtn.Show()
			m.attackBtn.Hide()
			m.engageBtn.Hide()
			m.fireBtn.Hide()
			m.chargeBtn.Show()
			m.rallyBtn.Hide()
			m.pursuitBtn.Show()
		case rp.Arm_ARTILLERY:
			m.marchBtn.Show()
			m.defendBtn.Show()
			m.attackBtn.Hide()
			m.engageBtn.Hide()
			m.fireBtn.Show()
			m.chargeBtn.Hide()
			m.rallyBtn.Hide()
			m.pursuitBtn.Hide()
		default:
			m.marchBtn.Hide()
			m.defendBtn.Hide()
			m.attackBtn.Hide()
			m.engageBtn.Hide()
			m.fireBtn.Hide()
			m.chargeBtn.Hide()
			m.rallyBtn.Hide()
			m.pursuitBtn.Hide()
		}
	}
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
	m.engageBtn.SetIcon(theme.RadioButtonIcon())
	m.fireBtn.SetIcon(theme.RadioButtonIcon())
	m.chargeBtn.SetIcon(theme.RadioButtonIcon())
	m.rallyBtn.SetIcon(theme.RadioButtonIcon())
	m.pursuitBtn.SetIcon(theme.RadioButtonIcon())
	switch m.order {
	case ORDER_MARRCH:
		m.marchBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_DEFEND:
		m.defendBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_ENGAGE:
		m.engageBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_ATTACK:
		m.attackBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_CHARGE:
		m.chargeBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_FIRE:
		m.fireBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_RALLY:
		m.rallyBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case ORDER_PURSUIT:
		m.pursuitBtn.SetIcon(theme.RadioButtonCheckedIcon())
	}
}
func (m *MapPanel) marchOrder() {
	m.setOrder(ORDER_MARRCH)
}

func (m *MapPanel) defendOrder() {
	m.setOrder(ORDER_DEFEND)
}

func (m *MapPanel) engageOrder() {
	m.setOrder(ORDER_ENGAGE)
}

func (m *MapPanel) attackOrder() {
	m.setOrder(ORDER_ATTACK)
}

func (m *MapPanel) chargeOrder() {
	m.setOrder(ORDER_CHARGE)
}

func (m *MapPanel) fireOrder() {
	m.setOrder(ORDER_FIRE)
}

func (m *MapPanel) rallyOrder() {
	m.setOrder(ORDER_RALLY)
}

func (m *MapPanel) pursuitOrder() {
	m.setOrder(ORDER_PURSUIT)
}

func (m *MapPanel) undoOrder() {
	m.setOrder(ORDER_NONE)
	m.SetCommand(nil)
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
	m.engageBtn = widget.NewButtonWithIcon("Engage", theme.RadioButtonIcon(), m.engageOrder)
	m.attackBtn = widget.NewButtonWithIcon("Attack", theme.RadioButtonIcon(), m.attackOrder)
	m.fireBtn = widget.NewButtonWithIcon("Fire", theme.RadioButtonIcon(), m.fireOrder)
	m.chargeBtn = widget.NewButtonWithIcon("Charge", theme.RadioButtonIcon(), m.chargeOrder)
	m.rallyBtn = widget.NewButtonWithIcon("Rally", theme.RadioButtonIcon(), m.rallyOrder)
	m.pursuitBtn = widget.NewButtonWithIcon("Pursuit", theme.RadioButtonIcon(), m.pursuitOrder)
	m.hbox1 = widget.NewHBox(
		widget.NewButtonWithIcon("Unit", theme.InfoIcon(), m.unitInfo),
		layout.NewSpacer(),
		m.unitDesc,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Undo", theme.CancelIcon(), m.undoOrder),
	)
	m.hbox2 = widget.NewHBox(
		layout.NewSpacer(),
		m.marchBtn,
		m.defendBtn,
		layout.NewSpacer(),
		m.engageBtn,
		m.attackBtn,
		m.chargeBtn,
		m.fireBtn,
		layout.NewSpacer(),
		m.rallyBtn,
		m.pursuitBtn,
		layout.NewSpacer(),
	)
	m.vbox = widget.NewVBox(m.hbox1, m.hbox2)

	m.mapWidget = newMapWidget(app, app.MapData, m.unitDesc)
	m.borderLayout = layout.NewBorderLayout(nil, m.vbox, nil, nil)
	m.content = fyne.NewContainerWithLayout(m.borderLayout,
		m.mapWidget,
		m.vbox,
	)
	m.SetCommand(nil)
	m.hbox2.Hide()

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
