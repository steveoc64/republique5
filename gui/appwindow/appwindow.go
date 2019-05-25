package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type appwindow struct {
	session *republique.Session
	window  fyne.Window

	header  *HeaderBar
	sidebar *SideBar
	footer  *widget.Box
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
	w.footer = widget.NewHBox()
	w.window.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(w.header.HBox, w.footer, w.sidebar.VBox, nil),
		w.header.HBox,
		w.sidebar.VBox))

	w.window.Canvas().SetOnTypedRune(w.typedRune)
	w.window.Canvas().SetOnTypedKey(w.typedKey)
}

func (w *appwindow) typedRune(r rune) {
}

func (w *appwindow) typedKey(ev *fyne.KeyEvent) {
}
