package appwindow

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

// TapRect is a Rect that implements Tappable interface
type TapRect struct {
	widget.Box
	rect        *canvas.Rectangle
	OnTapped    func()
	OnSecondary func()
}

// Tapped handler for each TapIcon
func (l *TapRect) Tapped(*fyne.PointEvent) {
	if l.OnTapped != nil {
		l.OnTapped()
	}
}

// TappedSecondary handler for each TapIcon
func (l *TapRect) TappedSecondary(*fyne.PointEvent) {
	if l.OnSecondary != nil {
		l.OnSecondary()
	}
}

// CreateRenderer creates a renderer for the TapIcon
func (l *TapRect) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(&l.Box)
}

// NewTapRect is a wrapper function to create a new tappable rectangle
func NewTapRect(c color.Color, tapped func(), secondary func()) *TapRect {
	t := &TapRect{
		Box:         widget.Box{},
		rect:        canvas.NewRectangle(c),
		OnTapped:    tapped,
		OnSecondary: secondary,
	}
	t.Box.Append(t.rect)
	return t
}
