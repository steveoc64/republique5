package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"github.com/steveoc64/republique5/republique"
)

type appwindow struct {
	session *republique.Session
	window  fyne.Window

	header  *HeaderBar
	sidebar *SideBar
	footer  *FooterBar
}

func Show(s *republique.Session, app fyne.App) {
	println("Session token", s.LoginDetails.Token.Id, "expires", s.LoginDetails.Token.Expires.String())
	w := &appwindow{session: s}
	w.loadUI(app)
	w.window.Show()
}

func (w *appwindow) loadUI(app fyne.App) {
	w.window = app.NewWindow("Republique 5.0")
	w.header = newHeaderBar(w.session)
	w.sidebar = newSideBar(w.session)
	w.footer = newFooterBar(w.session)
	w.window.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(w.header.Box, w.footer.Box, w.sidebar.Box, nil),
		w.header.Box,
		w.sidebar.Box,
		w.footer.Box,
	))

	w.window.Canvas().SetOnTypedRune(w.typedRune)
	w.window.Canvas().SetOnTypedKey(w.typedKey)
}

func (w *appwindow) typedRune(r rune) {
}

func (w *appwindow) typedKey(ev *fyne.KeyEvent) {
}
