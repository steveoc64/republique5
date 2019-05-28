package appwindow

import (
	rp "github.com/steveoc64/republique5/proto"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/hajimehoshi/oto"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

type App struct {
	app    fyne.App
	window fyne.Window
	conn   *grpc.ClientConn
	client rp.GameServiceClient

	ServerName string
	GameName   string
	GameTime   time.Time
	Briefing   string
	Commanders []string
	TeamName   string
	Token      string
	Expires    time.Time
	Phase      string

	layout          fyne.Layout
	container       *fyne.Container
	header          *HeaderBar
	sidebar         *SideBar
	footer          *FooterBar
	briefingPanel   *BriefingPanel
	actionsPanel    *ActionsPanel
	mapPanel        *MapPanel
	ordersPanel     *OrdersPanel
	unitsPanel      *UnitsPanel
	formationsPanel *FormationsPanel
	advancePanel    *AdvancePanel
	withdrawPanel   *WithdrawPanel
	surrenderPanel  *SurrenderPanel

	splashPanel *fyne.Container

	audioPort *oto.Player
}

func Show(app fyne.App, servername string, l *rp.LoginResponse, conn *grpc.ClientConn, client rp.GameServiceClient) {
	a := &App{
		app:        app,
		ServerName: servername,
		GameName:   l.GameName,
		GameTime:   time.Unix(l.GameTime.Seconds, 0),
		Briefing:   l.Briefing,
		Commanders: l.Commanders,
		TeamName:   l.TeamName,
		Token:      l.Token.Id,
		Expires:    time.Unix(l.Token.Expires.Seconds, 0),
		Phase:      "Pre Game Setup",
		conn:       conn,
		client:     client,
	}
	a.loadUI()
	a.window.CenterOnScreen()
	a.window.Show()
	a.PlayAudio("artillery")
	a.Phaser()
}

func (a *App) loadUI() {
	a.window = a.app.NewWindow("Republique 5.0")
	a.window.SetOnClosed(func() { os.Exit(0) })
	a.header = newHeaderBar(a)
	a.sidebar = newSideBar(a)
	a.footer = newFooterBar(a, a.endTurn)

	img := a.loadImage("banner")
	img.FillMode = canvas.ImageFillStretch
	a.splashPanel = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		img,
		widget.NewLabelWithStyle("Republique 5.0",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Italic: false}),
		widget.NewLabelWithStyle("Augmented Tabletop Miniatures",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: false, Italic: true}),
	)

	// Create the panels that go in the middle of the container, and then hide them
	a.briefingPanel = newBriefingPanel(a)
	a.actionsPanel = newActionsPanel(a)
	a.mapPanel = newMapPanel(a)
	a.ordersPanel = newOrdersPanel(a)
	a.unitsPanel = newUnitsPanel(a)
	a.formationsPanel = newFormationsPanel(a)
	a.advancePanel = newAdvancePanel(a)
	a.withdrawPanel = newWithdrawPanel(a)
	a.surrenderPanel = newSurrenderPanel(a)

	a.briefingPanel.Box.Hide()
	a.actionsPanel.Box.Hide()

	a.layout = layout.NewBorderLayout(a.header.Box, a.footer.Box, a.sidebar.Box, nil)
	a.container = fyne.NewContainerWithLayout(layout.NewBorderLayout(a.header.Box, a.footer.Box, a.sidebar.Box, nil),
		a.header.Box,
		a.sidebar.Box,
		a.footer.Box,
		a.splashPanel,
	)
	a.window.SetContent(a.container)

	a.window.Canvas().SetOnTypedRune(a.typedRune)
	a.window.Canvas().SetOnTypedKey(a.typedKey)
}

func (w *App) typedRune(r rune) {
}

func (w *App) typedKey(ev *fyne.KeyEvent) {
}

func (a *App) endTurn(done bool) {
	// TODO - transmit ImDone to the backend
}

func (a *App) loadImage(name string) *canvas.Image {
	dirname := filepath.Join(os.Getenv("HOME"), "republique")

	f := filepath.Join(dirname, name+".png")
	if _, err := os.Stat(f); err == nil {
		return canvas.NewImageFromFile(f)
	}
	f = filepath.Join(dirname, name+".jpg")
	if _, err := os.Stat(f); err == nil {
		return canvas.NewImageFromFile(f)
	}
	return nil
}

func (a *App) setPanel(p fyne.CanvasObject) {
	if len(a.container.Objects) == 4 {
		a.container.Objects[3].Hide()
	}
	a.container.Objects = append(a.container.Objects[:3], p)
	p.Show()
	a.container.Layout.Layout(a.container.Objects, a.container.Size())
	///a.container.Resize(a.container.Size())
}

func (a *App) showBriefing() {
	a.PlayAudio("command")
	a.setPanel(a.briefingPanel.Box)
}

func (a *App) showActions() {
	a.setPanel(a.actionsPanel.Box)
}

func (a *App) showMap() {
	a.setPanel(a.mapPanel.Box)
}

func (a *App) showOrders() {
	a.PlayAudio("command")
	a.setPanel(a.ordersPanel.Box)
}

func (a *App) showUnits() {
	a.PlayAudio("infantry")
	a.setPanel(a.unitsPanel.Box)
}

func (a *App) showFormations() {
	a.setPanel(a.formationsPanel.Box)
}

func (a *App) showAdvance() {
	a.setPanel(a.advancePanel.Box)
}

func (a *App) showWithdraw() {
	a.PlaySystemAudio("surrender")
	a.setPanel(a.withdrawPanel.Box)
}

func (a *App) showSurrender() {
	a.PlaySystemAudio("surrender")
	a.setPanel(a.surrenderPanel.Box)
}
