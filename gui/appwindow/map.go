package appwindow

import (
	"context"

	"github.com/davecgh/go-spew/spew"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	rp "github.com/steveoc64/republique5/proto"
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
	order        rp.Order
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
	if cmd == nil {
		m.setOrder(rp.Order_RESTAGE)
		m.mapWidget.grid.Select(0)
		widget.Renderer(m.mapWidget).Refresh()
		m.hbox2.Hide()
		m.unitDesc.SetText("")
		return
	}
	spew.Dump("Set command", cmd.GameState.GetOrders().String(), "path", cmd.GetGameState().GetObjective())
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
	m.setOrder(cmd.GetGameState().GetOrders())
}

func (m *MapPanel) unitInfo() {
	if m.command != nil {
		m.app.unitsPanel.ShowCommand(m.command)
		m.app.Tab(TAB_UNITS)
	}
}

func (m *MapPanel) gotoOrders() {
	m.app.Tab(TAB_ORDERS)
}

func (m *MapPanel) setOrder(o rp.Order) {
	if m.order == o {
		m.order = rp.Order_RESTAGE
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
	case rp.Order_MARCH:
		m.marchBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_DEFEND:
		m.defendBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_ENGAGE:
		m.engageBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_ATTACK:
		m.attackBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_CHARGE:
		m.chargeBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_FIRE:
		m.fireBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_RALLY:
		m.rallyBtn.SetIcon(theme.RadioButtonCheckedIcon())
	case rp.Order_PURSUIT:
		m.pursuitBtn.SetIcon(theme.RadioButtonCheckedIcon())
	}
	if m.command != nil {
		m.command.SetOrder(m.order)
	}
}
func (m *MapPanel) marchOrder() {
	m.setOrder(rp.Order_MARCH)
}

func (m *MapPanel) defendOrder() {
	m.setOrder(rp.Order_DEFEND)
}

func (m *MapPanel) engageOrder() {
	m.setOrder(rp.Order_ENGAGE)
}

func (m *MapPanel) attackOrder() {
	m.setOrder(rp.Order_ATTACK)
}

func (m *MapPanel) chargeOrder() {
	m.setOrder(rp.Order_CHARGE)
}

func (m *MapPanel) fireOrder() {
	m.setOrder(rp.Order_FIRE)
}

func (m *MapPanel) rallyOrder() {
	m.setOrder(rp.Order_RALLY)
}

func (m *MapPanel) pursuitOrder() {
	m.setOrder(rp.Order_PURSUIT)
}

func (m *MapPanel) clearOrder() {
	m.setOrder(rp.Order_RESTAGE)
	m.SetCommand(nil)
}

func (m *MapPanel) doneOrder() {
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
		widget.NewButtonWithIcon("Orders", theme.InfoIcon(), m.gotoOrders),
		layout.NewSpacer(),
		m.unitDesc,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Clear", theme.CancelIcon(), m.clearOrder),
		widget.NewButtonWithIcon("Done", theme.CheckButtonCheckedIcon(), m.doneOrder),
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

// Tap tells the map controller that a tap occurred at a spot
func (m *MapPanel) Tap(x, y int32) {
	if m.app.MapData.Side == rp.MapSide_TOP {
		x = m.app.MapData.X - x + 1
		y = m.app.MapData.Y - y + 1
	}
	println("got a map tap at", x, y)
	if x < 1 || x > m.app.MapData.X || y < 1 || y > m.app.MapData.Y {
		// out of bounds
		return
	}

	if m.command == nil {
		// no command set
		return
	}

	unitX := m.command.GameState.Grid.GetX()
	unitY := m.command.GameState.Grid.GetY()
	/*
		if m.app.MapData.Side == rp.MapSide_TOP {
			unitX = m.app.MapData.X - unitX + 1
			unitY = m.app.MapData.Y - unitY + 1
		}

	*/

	dx := (unitX - x)
	dy := (unitY - y)
	distance := dx*dx + dy*dy
	path := []*rp.Grid{}
	switch m.order {
	case rp.Order_RESTAGE:
		// nothing to do
		return
	case rp.Order_MARCH:
		// march to location if not too far
		maxd := int32(9)
		switch m.command.Arm {
		case rp.Arm_ARTILLERY:
			maxd = 8
		case rp.Arm_CAVALRY:
			maxd = 11
		}
		if distance > maxd {
			println("too far", distance)
			return
		}
		path = m.command.AddToObjective(x, y)
	case rp.Order_DEFEND:
		// defend at current location
		return
	case rp.Order_ENGAGE:
		// can engage out to 2 grids
		if distance > 4 {
			println("too far", distance)
			return
		}
		path = m.command.AddToObjective(x, y)
	case rp.Order_ATTACK:
		// can attack out to 1 grid
		if distance > 2 {
			println("too far", distance)
			return
		}
		path = m.command.SetObjective(x, y)
	case rp.Order_CHARGE:
		// check the range
		if distance > 9 {
			println("too far", distance)
			return
		}
		path = m.command.SetObjective(x, y)
	case rp.Order_PURSUIT:
		// check the range
		if distance > 12 {
			println("too far", distance)
			return
		}
		path = m.command.AddToObjective(x, y)
	case rp.Order_FIRE:
		// can fire out to adj grid
		if distance > 4 {
			println("too far", distance)
			return
		}
		path = m.command.SetObjective(x, y)
	}
	println("path redundant TODO rm me and the use of the path var", path)
	// orders have changed, so rebuild the orders panel
	// TODO - use databinding later when its available
	m.app.ordersPanel.build()
	widget.Renderer(m.mapWidget).Refresh()
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
