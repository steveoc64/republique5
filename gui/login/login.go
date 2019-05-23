package login

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type login struct {
	accessCodes [3][4]int
	mode        int
	i           int

	descr     *widget.Label
	code      [4]*widget.Label
	buttons   map[string]*widget.Button
	window    fyne.Window
	functions map[string]func()
}

func (c *login) GetWindow() fyne.Window {
	return c.window
}

func (c *login) paintCode() {
	for i := 0; i < 4; i++ {
		if i < c.i {
			c.code[i].SetText(fmt.Sprintf("%d", c.accessCodes[c.mode][i]))
		} else {
			c.code[i].SetText("_")
		}
	}
}

func (c *login) clear() {
	c.i = 0
	c.paintCode()
}

func (c *login) digit(d int) {
	c.accessCodes[c.mode][c.i] = d
	c.i++
	if c.i >= 4 {
		c.paintCode()
		time.Sleep(time.Millisecond * 600)
		c.setMode(c.mode + 1)
	}
	c.paintCode()
}

func (c *login) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button

	return button
}

func (c *login) digitButton(number int) *widget.Button {
	str := fmt.Sprintf("%d", number)
	action := func() {
		c.digit(number)
	}
	c.functions[str] = action

	return c.addButton(str, action)
}

func (c *login) typedRune(r rune) {
	action := c.functions[string(r)]
	if action != nil {
		action()
	}
}

func (c *login) typedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn, fyne.KeyEnter, fyne.KeyTab:
		c.setMode(c.mode + 1)
	case fyne.KeyUp:
		c.setMode(c.mode - 1)
	case fyne.KeyBackspace, fyne.KeyDelete:
		c.i--
		if c.i < 0 {
			c.i = 0
		}
		c.paintCode()
	}
}

func (c *login) ok() {
	c.setMode(c.mode + 1)
}

func (c *login) setMode(m int) {
	if m < 0 {
		m = 0
	}
	if m > 3 {
		m = 2
	}
	c.i = 0
	switch m {
	case 0:
		c.descr.SetText("Access Code")
	case 1:
		c.descr.SetText("Team Code")
	case 2:
		c.descr.SetText("Player Code")
	case 3:
		// all done ! call the server and validate
		os.Exit(1)
		return
	}
	c.mode = m
	c.paintCode()
}

func (c *login) loadUI(app fyne.App) {
	c.descr = widget.NewLabel("Access Code")
	c.descr.Alignment = fyne.TextAlignCenter
	c.mode = 0
	for i := 0; i < 4; i++ {
		c.code[i] = widget.NewLabel("_")
		c.code[i].Alignment = fyne.TextAlignCenter
		c.code[i].TextStyle.Monospace = true
	}

	c.window = app.NewWindow("Login")
	c.window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		c.descr,
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			c.code[0], c.code[1], c.code[2], c.code[3],
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.digitButton(0),
			c.addButton("OK", c.ok),
		),
	))

	c.window.Canvas().SetOnTypedRune(c.typedRune)
	c.window.Canvas().SetOnTypedKey(c.typedKey)
	c.window.Show()
}

func newLogin() *login {
	c := &login{}
	c.functions = make(map[string]func())
	c.buttons = make(map[string]*widget.Button)

	return c
}

// Show loads a calculator example window for the specified app context
func Show(app fyne.App) fyne.Window {
	c := newLogin()
	c.loadUI(app)
	return c.GetWindow()
}
