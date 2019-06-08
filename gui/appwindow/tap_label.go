package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// TapLabel is a label that implements Tappable interface
type TapLabel struct {
	widget.Label
	OnTapped    func()
	OnSecondary func()
}

// Tapped handler for each tapLabel
func (l *TapLabel) Tapped(*fyne.PointEvent) {
	if l.OnTapped != nil {
		l.OnTapped()
	}
}

// TappedSecondary handler for each tapLabel
func (l *TapLabel) TappedSecondary(*fyne.PointEvent) {
	if l.OnSecondary != nil {
		l.OnSecondary()
	}
}

// CreateRenderer creates a renderer for the TapLabel
func (l *TapLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(&l.Label)
}

// NewTapLabel is a wrapper function to create a new tappable Label
func NewTapLabel(text string, alignment fyne.TextAlign, style fyne.TextStyle, tapped func(), secondary func()) *TapLabel {
	return &TapLabel{
		widget.Label{
			Text:      text,
			Alignment: alignment,
			TextStyle: style,
		},
		tapped,
		secondary,
	}
}
