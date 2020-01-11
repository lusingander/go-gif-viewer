package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type navigateBar struct {
	bar         fyne.CanvasObject
	countLabel  *widget.Label
	countSlider *widget.Slider

	current, total int
	totalDigit     int
}

func (b *navigateBar) next() {
	if b.current < b.total {
		b.current++
		b.update()
	}
}

func (b *navigateBar) prev() {
	if b.current > 1 {
		b.current--
		b.update()
	}
}

func (b *navigateBar) change(n int) {
	if 1 <= n && n <= b.total {
		b.current = n
		b.update()
	}
}

func (b *navigateBar) update() {
	b.countLabel.SetText(b.createCountText())
	// Note: Slider doesn't have proper method...
	b.countSlider.Value = float64(b.current - 1)
	b.countSlider.Refresh()
}

func (b *navigateBar) createCountText() string {
	// 桁数によって伸び縮みしておりスライダーも伸び縮みしてしまう？
	return fmt.Sprintf("%*d/%*d",
		b.totalDigit, b.current, b.totalDigit, b.total)
}

func newNavigateBar(n int) *navigateBar {
	bar := &navigateBar{
		current: 1,
		total:   n,
	}
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), bar.prev)
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), bar.next)
	slider := widget.NewSlider(0, float64(n-1))
	slider.OnChanged = func(f float64) { bar.change(int(f) + 1) }
	bar.countSlider = slider
	bar.totalDigit = len(strconv.Itoa(n))
	count := widget.NewLabel(bar.createCountText())
	bar.countLabel = count
	buttons := widget.NewHBox(prev, next)
	bar.bar = fyne.NewContainerWithLayout(layout.NewBorderLayout(
		nil, nil, buttons, count), buttons, count, slider)
	return bar
}