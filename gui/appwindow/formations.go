package appwindow

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

// FormationsPanel is the UI for the formations
type FormationsPanel struct {
	app    *App
	Box    *widget.Box
	Scroll *widget.ScrollContainer
}

// CanvasObject returns the top level UI element for the formations panel
func (f *FormationsPanel) CanvasObject() fyne.CanvasObject {
	return f.Scroll
}

// newFormationsPanel creates a new formations panel and returns it
func newFormationsPanel(app *App) *FormationsPanel {
	h := &FormationsPanel{app: app}
	header := widget.NewLabelWithStyle("Division Level Formations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	notes := widget.NewLabelWithStyle("Each block shows the position of each brigade relative to each other", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	img := canvas.NewImageFromResource(resourceFormationsJpg)
	img.FillMode = canvas.ImageFillOriginal
	n0 := widget.NewLabelWithStyle("Each Brigade block above is organised in one of the following formations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	n1 := widget.NewLabel("Line Formation - Maximum Firepower - Not very Resilient")
	n1.Alignment = fyne.TextAlignCenter
	i1 := canvas.NewImageFromResource(resourceLineJpg)
	i1.FillMode = canvas.ImageFillOriginal
	n2 := widget.NewLabel("Supporting Lines - Adequate Firepower - Good Resilience")
	n2.Alignment = fyne.TextAlignCenter
	i2 := canvas.NewImageFromResource(resourceSupportingJpg)
	i2.FillMode = canvas.ImageFillOriginal
	n3 := widget.NewLabel("Attack Column - Maximum Bayonets and Shock Value - Good Cavalry Defence")
	n3.Alignment = fyne.TextAlignCenter
	i3 := canvas.NewImageFromResource(resourceAttackcolumnJpg)
	i3.FillMode = canvas.ImageFillOriginal
	n4 := widget.NewLabel("March Column - Fastest - Poor Defence - Slow to Deploy to other formations")
	n4.Alignment = fyne.TextAlignCenter
	i4 := canvas.NewImageFromResource(resourceMarchcolumnJpg)
	i4.FillMode = canvas.ImageFillOriginal
	n5 := widget.NewLabel("Echelon - Good Firepower - Good Flank Defence")
	n5.Alignment = fyne.TextAlignCenter
	i5 := canvas.NewImageFromResource(resourceEchelonJpg)
	i5.FillMode = canvas.ImageFillOriginal
	n6 := widget.NewLabel("Debande - All Bases Deployed as Skirmishers")
	n6.Alignment = fyne.TextAlignCenter
	i6 := canvas.NewImageFromResource(resourceDebandeJpg)
	i6.FillMode = canvas.ImageFillOriginal
	h.Box = widget.NewVBox(header, notes, img,
		n0,
		n1, i1,
		n2, i2,
		n3, i3,
		n4, i4,
		n5, i5,
		n6, i6)
	h.Scroll = widget.NewScrollContainer(h.Box)
	widget.Renderer(h.Box).Layout(h.Box.MinSize().Union(h.Box.Size()))
	h.Scroll.Resize(app.MinSize())
	return h
}
