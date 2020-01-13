package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

type navigateBar struct {
	bar         fyne.CanvasObject
	countLabel  *widget.Label
	countSlider *widget.Slider

	current, total int
	totalDigit     int

	observers []func(int)

	// TODO: fix
	delayMilliSec int
	stopPlay      chan bool
}

func (b *navigateBar) addObserver(f func(int)) {
	b.observers = append(b.observers, f)
}

func (b *navigateBar) start() {
	if b.stopPlay != nil {
		return
	}
	b.stopPlay = make(chan bool)
	go func() {
		t := time.NewTicker(time.Duration(b.delayMilliSec) * time.Millisecond)
		for {
			select {
			case <-t.C:
				b.next()
			case <-b.stopPlay:
				t.Stop()
				return
			}
		}
	}()
}

func (b *navigateBar) stop() {
	if b.stopPlay != nil {
		b.stopPlay <- true
		b.stopPlay = nil
	}
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
	for _, o := range b.observers {
		o(b.current - 1)
	}
}

func (b *navigateBar) createCountText() string {
	// 桁数によって伸び縮みしておりスライダーも伸び縮みしてしまう？
	return fmt.Sprintf("%*d/%*d",
		b.totalDigit, b.current, b.totalDigit, b.total)
}

func newNavigateBar(img *image.GIFImage) *navigateBar {
	n := img.Length()
	bar := &navigateBar{
		current:       1,
		total:         n,
		observers:     make([]func(int), 0),
		delayMilliSec: img.DelayMilliSec()[0], // TODO: fix
	}
	// TODO: use button with icon
	start := widget.NewButton("Start", bar.start)
	stop := widget.NewButton("Stop", bar.stop)
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), bar.prev)
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), bar.next)
	slider := widget.NewSlider(0, float64(n-1))
	slider.OnChanged = func(f float64) { bar.change(int(f) + 1) }
	bar.countSlider = slider
	bar.totalDigit = len(strconv.Itoa(n))
	count := widget.NewLabel(bar.createCountText())
	bar.countLabel = count
	buttons := widget.NewHBox(start, stop, prev, next)
	bar.bar = fyne.NewContainerWithLayout(layout.NewBorderLayout(
		nil, nil, buttons, count), buttons, count, slider)
	return bar
}
