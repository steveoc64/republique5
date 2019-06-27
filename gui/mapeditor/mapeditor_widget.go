package mapeditor

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// MapEditorWidget is a complete map editor
type MapEditorWidget struct {
	size       fyne.Size
	position   fyne.Position
	hidden     bool
	data       []byte
	datastring string
	x          int
	y          int
}

// NewMapEditorWidget creates and returns a new map editor widget
func NewMapEditorWidget() *MapEditorWidget {
	e := &MapEditorWidget{}
	return e
}

func (m *MapEditorWidget) SetMapSize(x, y int) {
	m.x = x
	m.y = y
	m.SetData(m.datastring)
	if r, ok := widget.Renderer(m).(*mapEditorRender); ok {
		r.dirty = true
		r.Refresh()
	}
}

func (m *MapEditorWidget) SetData(data string) {
	m.datastring = data
	m.data = make([]byte, m.x*m.y)
	lines := strings.Split(data, "\n")
	ii := 0
	for k, v := range lines {
		// break if too many lines
		if k >= m.y {
			break
		}
		// pad it out
		for i := 0; i < m.x; i++ {
			switch {
			case i < len(v):
				m.data[ii] = v[i]
			default:
				m.data[ii] = ' '
			}
			ii++
		}
	}
	if r, ok := widget.Renderer(m).(*mapEditorRender); ok {
		r.dirty = true
		r.Refresh()
	}
}

func (m *MapEditorWidget) Hide() {
	m.hidden = true
	for _, obj := range widget.Renderer(m).Objects() {
		obj.Show()
	}
}

func (m *MapEditorWidget) Size() fyne.Size {
	return m.size
}

func (m *MapEditorWidget) MinSize() fyne.Size {
	return fyne.Size{
		Width:  800,
		Height: 600,
	}
}

func (m *MapEditorWidget) Move(p fyne.Position) {
	m.position = p
}

func (m *MapEditorWidget) Position() fyne.Position {
	return m.position
}

func (m *MapEditorWidget) Resize(s fyne.Size) {
	m.size = s
	widget.Renderer(m).Layout(m.size)
}

func (m *MapEditorWidget) Show() {
	m.hidden = false
	for _, obj := range widget.Renderer(m).Objects() {
		obj.Show()
	}
}

func (m *MapEditorWidget) Visible() bool {
	return !m.hidden
}

func (m *MapEditorWidget) CreateRenderer() fyne.WidgetRenderer {
	return newMapEditorRender(m)
}
