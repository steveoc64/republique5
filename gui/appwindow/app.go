package appwindow

import (
	// ensure that jpep and png engines are imported
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/theme"
	rp "github.com/steveoc64/republique5/proto"

	"google.golang.org/grpc"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/hajimehoshi/oto"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

const (
	TAB_BRIEFING = iota
	TAB_UNITS
	TAB_ORDERS
	TAB_ACTIONS
	TAB_MAP
	TAB_FORMATION
	TAB_ADVANCE
	TAB_WITHDRAW
	TAB_SURRENDER
)

// App is the top level object to contain the application data
type App struct {
	// game state
	ServerName string
	GameName   string
	GameTime   time.Time
	Briefing   string
	Commanders []string
	Commands   []*rp.Command
	TeamName   string
	Token      rp.TokenMessage
	Expires    time.Time
	Phase      string
	MapData    *rp.MapData

	// comms and RPC stuff
	conn       *grpc.ClientConn
	gameServer rp.GameServiceClient

	// fyne layout and widgets
	app         fyne.App
	window      fyne.Window
	img         *canvas.Image
	layout      fyne.Layout
	container   *fyne.Container
	tabs        *widget.TabContainer
	isDarkTheme bool

	// panels in the tab container
	header          *HeaderBar
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

	// audio bits
	audioPort *oto.Player
}

// Show creates the application UI
func Show(app fyne.App, servername string, l *rp.LoginResponse, conn *grpc.ClientConn, gameServer rp.GameServiceClient) {
	a := &App{
		app:         app,
		ServerName:  servername,
		GameName:    l.GameName,
		GameTime:    time.Unix(l.GameTime.Seconds, 0),
		Briefing:    l.Briefing,
		Commanders:  l.Commanders,
		TeamName:    l.TeamName,
		Token:       rp.TokenMessage{Id: l.Token.Id},
		Expires:     time.Unix(l.Token.Expires.Seconds, 0),
		Phase:       "Pre Game Setup",
		conn:        conn,
		gameServer:  gameServer,
		isDarkTheme: true,
	}
	a.loadUI()
	a.window.CenterOnScreen()
	go func() {
		time.Sleep(100 * time.Millisecond)
		a.window.Show()
	}()
	a.PlayAudio("artillery")
	a.Phaser()
}

func (a *App) Tab(tab int) {
	println("setting tab to", tab)
	a.tabs.SelectTabIndex(tab)
}

// loadUI generates the UI elements
func (a *App) loadUI() {
	a.GetUnits()
	a.window = a.app.NewWindow("Republique 5.0")
	a.header = newHeaderBar(a)
	a.footer = newFooterBar(a, a.endTurn)
	a.img = canvas.NewImageFromResource(resourceBannerJpg)

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

	a.layout = layout.NewBorderLayout(a.header.CanvasObject(), a.footer.CanvasObject(), nil, nil)
	t := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Briefing", theme.FolderOpenIcon(), a.briefingPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Units", theme.InfoIcon(), a.unitsPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Orders", theme.DocumentCreateIcon(), a.ordersPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Actions", theme.ContentPasteIcon(), a.actionsPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Map", theme.ViewFullScreenIcon(), a.mapPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Formation", theme.ContentCopyIcon(), a.formationsPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Advance", theme.MailSendIcon(), a.advancePanel.CanvasObject()),
		widget.NewTabItemWithIcon("Withdraw", theme.MailReplyIcon(), a.withdrawPanel.CanvasObject()),
		widget.NewTabItemWithIcon("Surrender", theme.CancelIcon(), a.surrenderPanel.CanvasObject()),
	)
	t.SetTabLocation(widget.TabLocationLeading)
	a.tabs = t
	a.container = fyne.NewContainerWithLayout(layout.NewBorderLayout(a.header.Box, a.footer.Box, nil, nil),
		a.header.Box,
		a.footer.Box,
		t,
	)
	a.window.SetContent(a.container)
	a.window.SetOnClosed(func() { os.Exit(0) })
	a.window.Canvas().SetOnTypedRune(a.typedRune)
	a.window.Canvas().SetOnTypedKey(a.typedKey)
}

// MinSize returns the minimum size for the panel objects - helper function for
// panels to get a consistent sizing hint
func (a *App) MinSize() fyne.Size {
	return fyne.NewSize(1280, 900)
}

func (a *App) typedRune(r rune) {
}

func (a *App) typedKey(ev *fyne.KeyEvent) {
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

// ToggleTheme changes the theme between darkMode and lightMode and back again
func (a *App) ToggleTheme() {
	a.isDarkTheme = !a.isDarkTheme
	if a.isDarkTheme {
		a.app.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.app.Settings().SetTheme(theme.LightTheme())
	}
}

// GetUnitCommander gets the commander for the given unit
func (a *App) GetUnitCommander(unitID int32) *rp.Command {
	for _, command := range a.Commands {
		for _, u := range command.Units {
			if u.Id == unitID {
				return command
			}
		}
		for _, subCommand := range command.Subcommands {
			if subCommand.Id == unitID {
				return command
			}
			for _, u := range subCommand.Units {
				if u.Id == unitID {
					return subCommand
				}
			}
		}
	}
	return nil
}
