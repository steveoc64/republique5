package appwindow

import (
	"fyne.io/fyne/theme"
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
	Commands   []*rp.Command
	TeamName   string
	Token      rp.TokenMessage
	Expires    time.Time
	Phase      string

	layout          fyne.Layout
	container       *fyne.Container
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
		Token:      rp.TokenMessage{Id: l.Token.Id},
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

	a.layout = layout.NewBorderLayout(a.header.Box, a.footer.Box, nil, nil)
	t := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Briefing", theme.FolderIcon(), a.briefingPanel.Box),
		widget.NewTabItemWithIcon("Actions", theme.ContentPasteIcon(), a.actionsPanel.Box),
		widget.NewTabItemWithIcon("Map", theme.ViewFullScreenIcon(), a.mapPanel.Box),
		widget.NewTabItemWithIcon("Orders", theme.DocumentCreateIcon(), a.ordersPanel.Box),
		widget.NewTabItemWithIcon("Units", theme.InfoIcon(), a.unitsPanel.Box),
		widget.NewTabItemWithIcon("Formation", theme.ContentCopyIcon(), a.formationsPanel.Box),
		widget.NewTabItemWithIcon("Advance", theme.MailSendIcon(), a.advancePanel.Box),
		widget.NewTabItemWithIcon("Withdraw", theme.MailReplyIcon(), a.withdrawPanel.Box),
		widget.NewTabItemWithIcon("Surrender", theme.CancelIcon(), a.surrenderPanel.Box),
	)
	t.SetTabLocation(widget.TabLocationLeading)
	a.container = fyne.NewContainerWithLayout(layout.NewBorderLayout(a.header.Box, a.footer.Box, nil, nil),
		a.header.Box,
		a.footer.Box,
		t,
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
