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
	cx         int
	cy         int
	rivers     map[riverPoint]*river
}

type riverPoint struct {
	x int
	y int
}

type riverConnect struct {
	point riverPoint
	done  bool
}

type river struct {
	//adjacent []*riverConnect
	adjacent map[riverPoint]bool
}

func (m *MapEditorWidget) calcRiver() {
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
					adjacent: []*riverConnect{},
				}
			}
			i++
		}
	}

	// get the adjacent points
	for k, v := range m.rivers {
		for kk, _ := range m.rivers {
			dx := abs(k.x - kk.x)
			dy := abs(k.y - kk.y)
			if (dx == 1 && (dy == 1 || dy == 0)) ||
				(dy == 1 && (dx == 1 || dx == 0)) {
				// TODO - check here if the opposite exists
				v.adjacent = append(v.adjacent, &riverConnect{point: kk})
			}
		}
	}
}

// NewMapEditorWidget creates and returns a new map editor widget
func NewMapEditorWidget() *MapEditorWidget {
	e := &MapEditorWidget{
		cx: 1,
		cy: 1,
	}
	return e
}

func (m *MapEditorWidget) checkXY() bool {
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

func (m *MapEditorWidget) repaint() {
	if r, ok := widget.Renderer(m).(*mapEditorRender); ok {
		r.dirty = true
		r.Refresh()
	}
}

func (m *MapEditorWidget) SetMapSize(x, y int) {
	m.x = x
	m.y = y
	m.SetData(m.datastring)
	m.checkXY()
	m.repaint()
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
	m.calcRiver()
	m.repaint()
}

func (m *MapEditorWidget) setChar(b rune) {
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

func (m *MapEditorWidget) Key(event *fyne.KeyEvent) {
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

func (m *MapEditorWidget) Rune(r rune) {
	switch r {
	case 'h', 'H', ' ', 't', 'T', 'w', 'W', 'r':
		m.setChar(r)
	}
}

// Tapped is called when the user taps the map widget
func (m *MapEditorWidget) Tapped(event *fyne.PointEvent) {
	m.cx, m.cy = widget.Renderer(m).(*mapEditorRender).ConvertToGrid(event)
	m.repaint()
}

// TappedSecondary is called when the user right-taps the map widget
func (m *MapEditorWidget) TappedSecondary(event *fyne.PointEvent) {
}
