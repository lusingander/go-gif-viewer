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

var (
	playIcon  fyne.Resource = theme.NewThemedResource(resourcePlaySvg, nil)
	pauseIcon fyne.Resource = theme.NewThemedResource(resourcePauseSvg, nil)

	prevIcon  fyne.Resource = theme.NewThemedResource(resourcePrevSvg, nil)
	nextIcon  fyne.Resource = theme.NewThemedResource(resourceNextSvg, nil)
	firstIcon fyne.Resource = theme.NewThemedResource(resourceFirstSvg, nil)
	lastIcon  fyne.Resource = theme.NewThemedResource(resourceLastSvg, nil)
)

type playButton struct {
	*widget.Button

	playing     bool
	play, pause func()
}

func newPlayButton(play, pause func()) *playButton {
	return &playButton{
		Button:  widget.NewButtonWithIcon("", playIcon, play),
		playing: false,
		play:    play,
		pause:   pause,
	}
}

func (b *playButton) click() {
	if b.playing {
		b.OnTapped = b.play
		b.SetIcon(playIcon)
	} else {
		b.OnTapped = b.pause
		b.SetIcon(pauseIcon)
	}
	b.playing = !b.playing
}

type navigateBar struct {
	fyne.CanvasObject

	countLabel  *widget.Label
	countSlider *widget.Slider
	*playButton

	current, total int
	totalDigit     int

	observers []func(int)

	canPlay bool

	// TODO: fix
	delayMilliSec int
	stopPlay      chan bool
}

func (b *navigateBar) addObserver(f func(int)) {
	b.observers = append(b.observers, f)
}

func (b *navigateBar) clearObservers() {
	b.observers = make([]func(int), 0)
}

func (b *navigateBar) start() {
	if b.stopPlay != nil || !b.canPlay {
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
	b.playButton.click()
}

func (b *navigateBar) stop() {
	if b.stopPlay != nil {
		b.stopPlay <- true
		b.stopPlay = nil
		b.playButton.click()
	}
}

func (b *navigateBar) next() {
	if !b.canPlay {
		return
	}
	if b.current == b.total {
		b.first()
	} else if b.current < b.total {
		b.current++
		b.update()
	}
}

func (b *navigateBar) prev() {
	if !b.canPlay {
		return
	}
	if b.current == 1 {
		b.last()
	} else if b.current > 1 {
		b.current--
		b.update()
	}
}

func (b *navigateBar) first() {
	if !b.canPlay {
		return
	}
	if b.current > 1 {
		b.current = 1
		b.update()
	}
}

func (b *navigateBar) last() {
	if !b.canPlay {
		return
	}
	if b.current < b.total {
		b.current = b.total
		b.update()
	}
}

func (b *navigateBar) change(n int) {
	if !b.canPlay {
		return
	}
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
	if !b.canPlay {
		return ""
	}
	// 桁数によって伸び縮みしておりスライダーも伸び縮みしてしまう？
	return fmt.Sprintf("%*d/%*d",
		b.totalDigit, b.current, b.totalDigit, b.total)
}

func (b *navigateBar) setImage(img *image.GIFImage) {
	b.clearObservers()
	n := img.Length()
	b.current = 1
	b.total = n
	b.countSlider.Max = float64(n - 1)
	b.totalDigit = len(strconv.Itoa(n))
	b.delayMilliSec = img.DelayMilliSec()[0] // TODO: fix
	b.canPlay = true
	b.update()
}

func newNavigateBar() *navigateBar {
	bar := &navigateBar{
		observers: make([]func(int), 0),
		canPlay:   false,
	}
	slider := createSliderBar(bar)
	buttons := createButtons(bar)
	bar.CanvasObject = widget.NewVBox(slider, buttons)
	return bar
}

func createSliderBar(bar *navigateBar) fyne.CanvasObject {
	bar.countSlider = widget.NewSlider(0, 1)
	bar.countSlider.OnChanged = func(f float64) { bar.change(int(f) + 1) }
	bar.countLabel = widget.NewLabel(bar.createCountText())
	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, bar.countLabel),
		bar.countLabel, bar.countSlider,
	)
}

func createButtons(bar *navigateBar) fyne.CanvasObject {
	bar.playButton = newPlayButton(bar.start, bar.stop)
	first := widget.NewButtonWithIcon("", firstIcon, bar.first)
	prev := widget.NewButtonWithIcon("", prevIcon, bar.prev)
	next := widget.NewButtonWithIcon("", nextIcon, bar.next)
	last := widget.NewButtonWithIcon("", lastIcon, bar.last)
	return widget.NewHBox(
		layout.NewSpacer(), first, prev, bar.playButton.Button, next, last, layout.NewSpacer(),
	)
}
