package login

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/republique5/gui/appwindow"
	rp "github.com/steveoc64/republique5/republique/proto"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"google.golang.org/grpc"
)

type login struct {
	app         fyne.App
	servername  string
	accessCodes [3][4]int
	codeStrings [3]string
	mode        int
	i           int

	server    *widget.Entry
	descr     *widget.Label
	failed    *widget.Label
	code      [4]*widget.Label
	buttons   map[string]*widget.Button
	window    fyne.Window
	functions map[string]func()
	onLogin   func()
}

func Show(app fyne.App, servername string) {
	c := newLogin()
	c.app = app
	c.servername = servername
	c.loadUI()
}

func (c *login) paintCode() {
	for k, v := range c.accessCodes {
		c.codeStrings[k] = fmt.Sprintf("%d%d%d%d", v[0], v[1], v[2], v[3])
	}
	for i := 0; i < 4; i++ {
		if i < c.i {
			if c.i >= 4 {
				c.code[i].SetText("*")
			} else {
				c.code[i].SetText(fmt.Sprintf("%d", c.accessCodes[c.mode][i]))
			}
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
	c.failed.Hide()
	c.accessCodes[c.mode][c.i] = d
	c.i++
	if c.i >= 4 {
		c.paintCode()
		time.Sleep(time.Millisecond * 200)
		c.setMode(c.mode + 1)
	}
	c.paintCode()
}

func (c *login) addIconButton(text string, icon fyne.Resource, action func()) *widget.Button {
	button := widget.NewButtonWithIcon(text, icon, action)
	c.buttons[text] = button

	return button
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
	case fyne.KeyBackspace, fyne.KeyDelete, fyne.KeyEscape:
		if c.i == 0 {
			c.setMode(c.mode - 1)
			return
		}
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

func (c *login) del() {
	c.typedKey(&fyne.KeyEvent{Name: fyne.KeyDelete})
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
		if err := c.login(); err != nil {
			c.setMode(0)
			c.failed.Show()
			return
		}
		return
	}
	c.mode = m
	c.paintCode()
}

func (c *login) login() error {
	println("Connecting to server", c.server.Text, "AccessCodes", c.codeStrings[0], c.codeStrings[1], c.codeStrings[2])
	serverAddr := c.server.Text
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := rp.NewGameServiceClient(conn)
	rsp, err := client.Login(context.Background(), &rp.LoginMessage{
		AccessCode: c.codeStrings[0],
		TeamCode:   c.codeStrings[1],
		PlayerCode: c.codeStrings[2],
	})
	if err != nil {
		return err
	}
	appwindow.Show(c.app, c.servername, rsp)
	c.window.Hide()
	return nil
}

func (c *login) loadUI() {
	c.server = widget.NewEntry()
	c.server.SetText(c.servername)
	c.descr = widget.NewLabel("Access Code")
	c.descr.Alignment = fyne.TextAlignCenter
	c.failed = widget.NewLabel("Failed - Try Again")
	c.failed.Alignment = fyne.TextAlignCenter
	c.failed.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	c.mode = 0

	for i := 0; i < 4; i++ {
		c.code[i] = widget.NewLabel("_")
		c.code[i].Alignment = fyne.TextAlignCenter
		c.code[i].TextStyle.Monospace = true
	}

	c.window = c.app.NewWindow("Login")
	c.window.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		c.server,
		c.failed,
		c.descr,
		fyne.NewContainerWithLayout(layout.NewGridLayout(6),
			widget.NewLabel(" "),
			c.code[0], c.code[1], c.code[2], c.code[3],
			widget.NewLabel(" "),
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.addIconButton("Del", theme.CancelIcon(), c.del),
			c.digitButton(0),
			c.addIconButton("OK", theme.ConfirmIcon(), c.ok),
		),
	))

	c.window.Canvas().SetOnTypedRune(c.typedRune)
	c.window.Canvas().SetOnTypedKey(c.typedKey)
	c.window.Show()
	c.window.CenterOnScreen()
	c.failed.Hide()
}

func newLogin() *login {
	c := &login{}
	c.functions = make(map[string]func())
	c.buttons = make(map[string]*widget.Button)

	return c
}
