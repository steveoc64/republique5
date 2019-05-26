package appwindow

import (
	"fyne.io/fyne/widget"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	rp "github.com/steveoc64/republique5/republique/proto"
)

type App struct {
	app    fyne.App
	window fyne.Window

	ServerName string
	GameName   string
	GameTime   time.Time
	Briefing   string
	Commanders []string
	TeamName   string
	Token      string
	Expires    time.Time
	Phase      string

	layout        fyne.Layout
	container     *fyne.Container
	header        *HeaderBar
	sidebar       *SideBar
	footer        *FooterBar
	briefingPanel *BriefingPanel
	splashPanel   *fyne.Container
}

func Show(app fyne.App, servername string, l *rp.LoginResponse) {
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
	}
	a.loadUI()
	a.window.CenterOnScreen()
	a.window.Show()
}

func (a *App) loadUI() {
	a.window = a.app.NewWindow("Republique 5.0")
	a.header = newHeaderBar(a)
	a.sidebar = newSideBar(a)
	a.footer = newFooterBar(a, a.endTurn)

	a.splashPanel = fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		widget.NewLabelWithStyle(
			"\n\n\n\nRepublique 5.0\n\n\nAugmented Reality Tabletop Miniatures",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Italic: true}))

	// Create the panels that go in the middle of the container, and then hide them
	a.briefingPanel = newBriefingPanel(a)
	a.briefingPanel.Box.Hide()

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
	println("end turn", done)
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
	a.setPanel(a.briefingPanel.Box)
}
