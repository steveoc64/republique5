package mainwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/steveoc64/republique5/republique"
)

type mainwindow struct {
	session *republique.Session
	window  fyne.Window

	header  *widget.Box
	sidebar *widget.Box
	footer  *widget.Box
}

func Show(s *republique.Session, app fyne.App) {
	w := &mainwindow{session: s}
	w.loadUI(app)
	w.window.Show()
	println("Session token", s.LoginDetails.Token.Id, "expires", s.LoginDetails.Token.Expires.String())
}

func (w *mainwindow) loadUI(app fyne.App) {
	w.window = app.NewWindow("Republique 5.0")
	w.header = widget.NewHBox()
	w.sidebar = widget.NewVBox()
	w.footer = widget.NewHBox()
	w.window.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(w.header, w.footer, w.sidebar, nil)))

	w.window.Canvas().SetOnTypedRune(w.typedRune)
	w.window.Canvas().SetOnTypedKey(w.typedKey)
}

func (w *mainwindow) typedRune(r rune) {
}

func (w *mainwindow) typedKey(ev *fyne.KeyEvent) {
}
