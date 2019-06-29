package appwindow

type riverPoint struct {
	x int32
	y int32
}

type river struct {
	adjacent []riverPoint
}

func (m *MapWidget) calcRiver() {
	m.rivers = make(map[riverPoint]*river)

	abs := func(i int32) int32 {
		if i < 0 {
			return -1 * i
		}
		return i
	}

	i := 0
	// create all the riverpoints
	for y := int32(0); y < m.grid.y; y++ {
		for x := int32(0); x < m.grid.x; x++ {
			if m.grid.value[i] == 'r' {
				m.rivers[riverPoint{x, y}] = &river{
					adjacent: []riverPoint{},
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
				v.adjacent = append(v.adjacent, kk)
			}
		}
	}
}
