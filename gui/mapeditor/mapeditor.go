package mapeditor

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"strconv"
)

// MapEditor struct contains a mapeditor
type mapeditor struct {
	app fyne.App
	w   fyne.Window

	form *widget.Form
	x    *widget.Entry
	y    *widget.Entry
	data *widget.Entry
	m    *MapEditorWidget
}

// New creates a new map editor
func New(app fyne.App, x int, y int, data string) {
	m := &mapeditor{
		app:  app,
		w:    app.NewWindow("Map Editor"),
		form: widget.NewForm(),
		x:    widget.NewEntry(),
		y:    widget.NewEntry(),
		data: widget.NewMultiLineEntry(),
		m:    NewMapEditorWidget(),
	}
	m.form.OnSubmit = m.Submit
	m.w.SetContent(widget.NewVBox(m.form, m.m))

	m.form.AppendItem(&widget.FormItem{
		Text:   "X",
		Widget: m.x,
	})
	m.form.AppendItem(&widget.FormItem{
		Text:   "Y",
		Widget: m.y,
	})
	m.form.AppendItem(&widget.FormItem{
		Text:   "Map Data",
		Widget: m.data,
	})
	m.form.AppendItem(&widget.FormItem{
		Text:   "",
		Widget: widget.NewButton("Submit", m.Submit),
	})
	m.form.AppendItem(&widget.FormItem{
		Text:   "",
		Widget: widget.NewButton("Quit", m.Quit),
	})
	m.x.SetText(fmt.Sprintf("%d", x))
	m.y.SetText(fmt.Sprintf("%d", y))
	m.data.SetText(data)
	m.m.SetMapSize(x, y)
	m.m.SetData(data)
	m.w.Show()
	m.w.CenterOnScreen()
}

func (m *mapeditor) Submit() {
	x, err := strconv.Atoi(m.x.Text)
	if err != nil {
		return
	}
	y, err := strconv.Atoi(m.y.Text)
	if err != nil {
		return
	}
	m.m.SetMapSize(x, y)
	m.m.SetData(m.data.Text)
}

func (m *mapeditor) Quit() {
	x, _ := strconv.Atoi(m.x.Text)
	for i := 0; i < x; i++ {
		print("-")
	}
	print("\n")
	println(m.m.datastring)
	for i := 0; i < x; i++ {
		print("-")
	}
	print("\n")
	m.app.Quit()
}