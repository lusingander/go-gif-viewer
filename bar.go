package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

var (
	openIcon  fyne.Resource = theme.NewThemedResource(resourceOpenSvg, nil)
	closeIcon fyne.Resource = theme.NewThemedResource(resourceCloseSvg, nil)

	infoIcon    fyne.Resource = theme.NewThemedResource(theme.InfoIcon(), nil)
	creditsIcon fyne.Resource = theme.NewThemedResource(resourceCreditsSvg, nil) // TODO: better icon...
	aboutIcon   fyne.Resource = theme.NewThemedResource(theme.QuestionIcon(), nil)

	zoomInIcon  fyne.Resource = theme.NewThemedResource(theme.ZoomInIcon(), nil)
	zoomOutIcon fyne.Resource = theme.NewThemedResource(theme.ZoomOutIcon(), nil)

	playIcon  fyne.Resource = theme.NewThemedResource(theme.MediaPlayIcon(), nil)
	pauseIcon fyne.Resource = theme.NewThemedResource(theme.MediaPauseIcon(), nil)

	prevIcon  fyne.Resource = theme.NewThemedResource(theme.MediaFastRewindIcon(), nil)
	nextIcon  fyne.Resource = theme.NewThemedResource(theme.MediaFastForwardIcon(), nil)
	firstIcon fyne.Resource = theme.NewThemedResource(theme.MediaSkipPreviousIcon(), nil)
	lastIcon  fyne.Resource = theme.NewThemedResource(theme.MediaSkipNextIcon(), nil)
)

var speeds = []string{
	"0.25x",
	"0.5x",
	"1.0x",
	"2.0x",
}

func defaultSpeed() string {
	return speeds[2] // 1.0x
}

type menuBar struct {
	fyne.CanvasObject

	*widget.Toolbar
	*widget.Select

	*player
}

func createSpeedSelect() *widget.Select {
	// Note: when using NewSelect, defaultPlaceHolder cannot be set, so the width of Select is specified
	// TODO: do not show shadow / move to toolbar
	sel := &widget.Select{
		BaseWidget:  widget.BaseWidget{},
		Selected:    defaultSpeed(),
		Options:     speeds,
		PlaceHolder: ".", // this value determines the width!
	}
	return sel
}

func newMenuBar(open, close, info, credits, about, zoomIn, zoomOut func()) *menuBar {
	toolBar := widget.NewToolbar(
		widget.NewToolbarAction(openIcon, open),
		widget.NewToolbarAction(closeIcon, close),
		widget.NewToolbarAction(infoIcon, info),
		widget.NewToolbarAction(creditsIcon, credits),
		widget.NewToolbarAction(aboutIcon, about),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(zoomInIcon, zoomIn),
		widget.NewToolbarAction(zoomOutIcon, zoomOut),
	)
	sppedSelect := createSpeedSelect()
	return &menuBar{
		CanvasObject: fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, nil, nil, sppedSelect),
			sppedSelect, toolBar,
		),
		Toolbar: toolBar,
		Select:  sppedSelect,
	}
}

func (b *menuBar) setPlayer(p *player) {
	b.player = p
	b.Select.OnChanged = func(s string) {
		// TODO: fix
		if b.player != nil {
			b.player.setSpeed(parseSpeed(s))
		}
	}
	b.Select.SetSelected(defaultSpeed())
}

func (b *menuBar) currentSpeed() float64 {
	return parseSpeed(b.Select.Selected)
}

func parseSpeed(s string) float64 {
	f, _ := strconv.ParseFloat(s[:len(s)-1], 64)
	return f
}

func (b *menuBar) increaseSpeed() {
	// TODO: fix
	if b.player == nil {
		return
	}
	for i := 0; i < len(speeds)-1; i++ {
		if speeds[i] == b.Selected {
			b.SetSelected(speeds[i+1])
			b.player.setSpeed(b.currentSpeed())
			return
		}
	}
}

func (b *menuBar) decreaseSpeed() {
	// TODO: fix
	if b.player == nil {
		return
	}
	for i := 1; i < len(speeds); i++ {
		if speeds[i] == b.Selected {
			b.SetSelected(speeds[i-1])
			b.player.setSpeed(b.currentSpeed())
			return
		}
	}
}

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

	*widget.Label
	*widget.Slider
	*playButton
	*player

	current, total int
	totalDigit     int

	observers []func(int)

	canPlay bool
}

func (b *navigateBar) addObserver(f func(int)) {
	b.observers = append(b.observers, f)
}

func (b *navigateBar) clearObservers() {
	b.observers = make([]func(int), 0)
}

func (b *navigateBar) pressPlayButton() {
	if b.player == nil {
		return
	}
	if b.player.playing() {
		b.stop()
	} else {
		b.start()
	}
}

func (b *navigateBar) start() {
	if !b.canPlay || b.player.playing() {
		return
	}
	b.playButton.click()
	b.player.play()
}

func (b *navigateBar) stop() {
	if !b.player.playing() {
		return
	}
	b.player.pause()
	b.playButton.click()
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
	b.SetText(b.createCountText())
	// Note: Slider doesn't have proper method...
	b.Slider.Value = float64(b.current - 1)
	b.Slider.Refresh()
	for _, o := range b.observers {
		o(b.current - 1)
	}
}

func (b *navigateBar) createCountText() string {
	if !b.canPlay {
		return ""
	}
	return fmt.Sprintf("%*d/%*d",
		b.totalDigit, b.current, b.totalDigit, b.total)
}

func (b *navigateBar) setImage(img *image.GIFImage) {
	b.clearObservers()
	n := img.Length()
	b.current = 1
	b.total = n
	b.Slider.Max = float64(n - 1)
	b.totalDigit = len(strconv.Itoa(n))
	b.canPlay = true
	b.update()
}

func (b *navigateBar) clearImage() {
	b.clearObservers()
	b.current = 1
	b.Slider.Max = 1
	b.canPlay = false
	b.update()
}

func (b *navigateBar) setPlayer(p *player) {
	b.player = p
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
	bar.Slider = widget.NewSlider(0, 1)
	bar.Slider.OnChanged = func(f float64) { bar.change(int(f) + 1) }
	bar.Label = widget.NewLabel(bar.createCountText())
	bar.Label.TextStyle.Monospace = true
	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, bar.Label),
		bar.Label, bar.Slider,
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
