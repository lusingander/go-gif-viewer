package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"github.com/lusingander/go-gif-viewer/image"
	fd "github.com/sqweek/dialog"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

type mainView struct {
	fyne.Window

	*fyne.Container

	*menuBar
	*imageView
	*navigateBar

	*player
}

func newMainView(w fyne.Window) *mainView {
	mainView := &mainView{
		Window: w,
	}
	menuBar := newMenuBar(
		mainView.openFileDialog,
		mainView.clearImage,
		mainView.info,
		mainView.credits,
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
	f, err := fd.File().Filter("GIF", "gif").Load()
	defer v.Window.RequestFocus() // dialog does not return focus...
	if err != nil {
		if err != fd.ErrCancelled {
			dialog.ShowError(err, v.Window)
		}
		return
	}
	err = v.loadImageFromPath(f)
	if err != nil {
		dialog.ShowError(err, v.Window)
		return
	}
}

func (v *mainView) info() {
	if v.imageView.GIFImage == nil {
		return
	}
	showInfoWindow(v.imageView.GIFImage)
}

func (v *mainView) credits() {
	CreditsWindow(fyne.CurrentApp()).Show()
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

func (v *mainView) addShortcuts() {
	// TODO: want to use ctrl on Windows...
	v.addSuperShotrcuts(fyne.KeyO, v.openFileDialog)
	v.addSuperShotrcuts(fyne.KeyW, v.clearImage)
}

func (v *mainView) addSuperShotrcuts(key fyne.KeyName, f func()) {
	v.Window.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: key, Modifier: desktop.SuperModifier},
		func(_ fyne.Shortcut) { f() },
	)
}

func run(args []string) error {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	v := newMainView(w)
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
	v.addShortcuts()
	if len(args) > 1 {
		err := v.loadImageFromPath(args[1])
		if err != nil {
			dialog.ShowError(err, w)
		}
	}
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
