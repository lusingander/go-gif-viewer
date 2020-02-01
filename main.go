package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"github.com/lusingander/go-gif-viewer/image"
	"github.com/sqweek/dialog"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

type mainView struct {
	*fyne.Container

	*menuBar
	*imageView
	*navigateBar

	*player
}

func newMainView() *mainView {
	mainView := &mainView{}
	menuBar := newMenuBar(
		mainView.openFileDialog,
		mainView.clearImage,
		mainView.zoomIn,
		mainView.zoomOut,
	)
	imageView := newImageView()
	navigateBar := newNavigateBar()
	panel := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(menuBar.CanvasObject, navigateBar.CanvasObject, nil, nil),
		menuBar.CanvasObject, navigateBar.CanvasObject, imageView.CanvasObject,
	)
	mainView.Container = panel
	mainView.menuBar = menuBar
	mainView.imageView = imageView
	mainView.navigateBar = navigateBar
	return mainView
}

func (v *mainView) loadImageFromPath(path string) error {
	img, err := image.LoadGIFImageFromPath(path)
	if err != nil {
		return err
	}
	v.loadImage(img)
	return nil
}

func (v *mainView) loadImage(img *image.GIFImage) {
	v.imageView.setImage(img)
	v.navigateBar.setImage(img)
	v.navigateBar.addObserver(v.refleshFrame)

	v.player = newPlayer(img, v.menuBar.currentSpeed())
	v.player.addObserver(v.navigateBar.next)
	v.menuBar.setPlayer(v.player)
	v.navigateBar.setPlayer(v.player)
}

func (v *mainView) clearImage() {
	v.imageView.clearImage()
	v.navigateBar.clearImage()

	v.player = nil
	v.menuBar.setPlayer(v.player)
	v.navigateBar.setPlayer(v.player)
}

func (v *mainView) openFileDialog() {
	// TODO: https://github.com/sqweek/dialog/issues/24
	f, err := dialog.File().Filter("GIF", "gif").Load()
	if err != nil {
		return
	}
	v.loadImageFromPath(f)
}

func (v *mainView) zoomIn() {
	v.imageView.zoomIn()
	v.navigateBar.update()
}

func (v *mainView) zoomOut() {
	v.imageView.zoomOut()
	v.navigateBar.update()
}

func (v *mainView) handleKeys(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyLeft:
		v.navigateBar.prev()
	case fyne.KeyRight:
		v.navigateBar.next()
	case fyne.KeyUp:
		v.navigateBar.first()
	case fyne.KeyDown:
		v.navigateBar.last()
	case fyne.KeySpace:
		v.navigateBar.pressPlayButton()
	}
}

func (v *mainView) handleRune(r rune) {
	switch r {
	case '+':
		v.zoomIn()
	case '-':
		v.zoomOut()
	}
}

func run(args []string) error {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	v := newMainView()
	w.SetContent(v.Container)
	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", v.openFileDialog),
			fyne.NewMenuItem("Close", v.clearImage),
		),
		fyne.NewMenu("View",
			fyne.NewMenuItem("Zoom In", v.zoomIn),
			fyne.NewMenuItem("Zoom Out", v.zoomOut),
		),
	))
	w.Canvas().SetOnTypedKey(v.handleKeys)
	w.Canvas().SetOnTypedRune(v.handleRune)
	if len(args) > 1 {
		v.loadImageFromPath(args[1])
	}
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
