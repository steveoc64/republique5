package mapeditor

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// Widget is a complete map editor as a widget
type Widget struct {
	size       fyne.Size
	position   fyne.Position
	hidden     bool
	data       []byte
	datastring string
	x          int
	y          int
	cx         int
	cy         int
	rivers     map[riverPoint]*river
}

type riverPoint struct {
	x int
	y int
}

type river struct {
	adjacent map[riverPoint]bool
}

func (m *Widget) calcRiver() {
	m.rivers = make(map[riverPoint]*river)

	abs := func(i int) int {
		if i < 0 {
			return -1 * i
		}
		return i
	}

	i := 0
	// create all the riverpoints
	for y := 0; y < m.y; y++ {
		for x := 0; x < m.x; x++ {
			if m.data[i] == 'r' {
				m.rivers[riverPoint{x, y}] = &river{
					adjacent: make(map[riverPoint]bool),
				}
			}
			i++
		}
	}

	// get the adjacent points
	for k, v := range m.rivers {
		for kk := range m.rivers {
			dx := abs(k.x - kk.x)
			dy := abs(k.y - kk.y)
			if (dx == 1 && (dy == 1 || dy == 0)) ||
				(dy == 1 && (dx == 1 || dx == 0)) {
				v.adjacent[riverPoint{kk.x, kk.y}] = false
			}
		}
	}
}

// NewMapEditorWidget creates and returns a new map editor widget
func NewMapEditorWidget() *Widget {
	e := &Widget{
		cx: 1,
		cy: 1,
	}
	return e
}

func (m *Widget) checkXY() bool {
	if m.cx < 1 {
		m.cx = 1
		return false
	}
	if m.cy < 1 {
		m.cy = 1
		return false
	}
	if m.cx > m.x {
		m.cx = m.x
		return false
	}
	if m.cy > m.y {
		m.cy = m.y
		return false
	}
	return true
}

func (m *Widget) repaint() {
	if r, ok := widget.Renderer(m).(*mapEditorRender); ok {
		r.dirty = true
		r.Refresh()
	}
}

// SetMapSize sets the map size by grid units
func (m *Widget) SetMapSize(x, y int) {
	m.x = x
	m.y = y
	m.SetData(m.datastring)
	m.checkXY()
	m.repaint()
}

// SetData sets the contents of the map
func (m *Widget) SetData(data string) {
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
	m.calcRiver()
	m.repaint()
}

func (m *Widget) setChar(b rune) {
	i := (m.cy-1)*m.x + (m.cx - 1)
	if i < 0 || i >= len(m.data) {
		println("error", m.cx, m.cy, m.x, m.y)
		return
	}
	switch b {
	case 'h', 'w', 't':
		if m.data[i] == byte(b) {
			b = b + 'A' - 'a'
		}
	}
	m.data[i] = byte(b)
	data := ""
	i = 0
	for y := 0; y < m.y; y++ {
		for x := 0; x < m.x; x++ {
			switch m.data[i] {
			case 0:
				m.data[i] = ' '
			}
			data = data + string(m.data[i])
			i++
		}
		data = data + "\n"
	}
	m.datastring = data

	m.calcRiver()
	m.repaint()
}

// Hide is called to hide the widget
func (m *Widget) Hide() {
	m.hidden = true
	for _, obj := range widget.Renderer(m).Objects() {
		obj.Show()
	}
}

// Size returns the size of the widget
func (m *Widget) Size() fyne.Size {
	return m.size
}

// MinSize returns the minimum size of the widget
func (m *Widget) MinSize() fyne.Size {
	return fyne.Size{
		Width:  800,
		Height: 600,
	}
}

// Move is called when the widget is to be moved
func (m *Widget) Move(p fyne.Position) {
	m.position = p
}

// Position retuns the current position of the widget
func (m *Widget) Position() fyne.Position {
	return m.position
}

// Resize is called when the widget needs to be resized
func (m *Widget) Resize(s fyne.Size) {
	m.size = s
	widget.Renderer(m).Layout(m.size)
}

// Show makes the widget visible
func (m *Widget) Show() {
	m.hidden = false
	for _, obj := range widget.Renderer(m).Objects() {
		obj.Show()
	}
}

// Visible returns true / false whether the widget is visible
func (m *Widget) Visible() bool {
	return !m.hidden
}

// CreateRenderer creates a renderer for the widget
func (m *Widget) CreateRenderer() fyne.WidgetRenderer {
	return newMapEditorRender(m)
}

// Key is called when a key is pressed
func (m *Widget) Key(event *fyne.KeyEvent) {
	switch event.Name {
	case "Right":
		m.cx++
		if m.cx > m.x {
			m.SetMapSize(m.x+1, m.y)
			return
		}
	case "Left":
		m.cx--
	case "Up":
		m.cy--
	case "Down":
		m.cy++
		if m.cy > m.y {
			m.SetMapSize(m.x, m.y+1)
			return
		}
	}
	if !m.checkXY() {
		return
	}
	m.repaint()
}

// Rune is called when a user types a char
func (m *Widget) Rune(r rune) {
	switch r {
	case 'h', 'H', ' ', 't', 'T', 'w', 'W', 'r':
		m.setChar(r)
	}
}

// Tapped is called when the user taps the map widget
func (m *Widget) Tapped(event *fyne.PointEvent) {
	m.cx, m.cy = widget.Renderer(m).(*mapEditorRender).ConvertToGrid(event)
	m.repaint()
}

// TappedSecondary is called when the user right-taps the map widget
func (m *Widget) TappedSecondary(event *fyne.PointEvent) {
}
