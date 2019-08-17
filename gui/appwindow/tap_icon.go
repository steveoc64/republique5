package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

// TapIcon is an Icon that implements Tappable interface
type TapIcon struct {
	widget.Box
	img         *canvas.Image
	OnTapped    func()
	OnSecondary func()
	enabled     bool
}

// Tapped handler for each TapIcon
func (l *TapIcon) Tapped(*fyne.PointEvent) {
	if l.enabled && l.OnTapped != nil {
		l.OnTapped()
	}
}

// TappedSecondary handler for each TapIcon
func (l *TapIcon) TappedSecondary(*fyne.PointEvent) {
	if l.enabled && l.OnSecondary != nil {
		l.OnSecondary()
	}
}

// CreateRenderer creates a renderer for the TapIcon
func (l *TapIcon) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(&l.Box)
}

// Disable makes it disabled
func (l *TapIcon) Disable() {
	l.enabled = false
	l.img.Translucency = 1.0
}

// Enable makes it enabled
func (l *TapIcon) Enable() {
	l.enabled = true
	l.img.Translucency = 0.0
}

// NewTapIcon is a wrapper function to create a new tappable icon
func NewTapIcon(res fyne.Resource, tapped func(), secondary func()) *TapIcon {
	t := &TapIcon{
		Box:         widget.Box{},
		img:         canvas.NewImageFromResource(res),
		OnTapped:    tapped,
		OnSecondary: secondary,
		enabled:     true,
	}
	t.img.SetMinSize(fyne.Size{Width: 64, Height: 64})
	t.Box.Append(t.img)
	return t
}
