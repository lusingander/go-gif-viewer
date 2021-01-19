package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	appName = "GIF Viewer"
	version = "0.2.2"
)

var defaultWindowSize = fyne.NewSize(400, 400)

type mainView struct {
	fyne.Window

	*fyne.Container

	*menuBar
	*imageView
	*navigateBar

	*player

	isThumbnailListWindowOpening bool
	isInfoWindowOpening          bool
	isCreditsWindowOpening       bool
	isAboutWindowOpening         bool
}

func newMainView(w fyne.Window) *mainView {
	mainView := &mainView{
		Window: w,
	}
	menuBar := newMenuBar(
		mainView.openFileDialog,
		mainView.clearImage,
		mainView.thumbnailList,
		mainView.info,
		mainView.credits,
		mainView.about,
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
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, v.Window)
			return
		}
		v.withLoadingDialog(func() {
			err = v.loadImageFromPath(reader.URI().String()[7:]) // `file://`
			if err != nil {
				dialog.ShowError(err, v.Window)
				return
			}
		})
	}, v.Window)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".gif"}))
	fd.Show()
}

func (v *mainView) withLoadingDialog(f func()) {
	d := dialog.NewProgressInfinite("Loading", "Now loading...", v.Window)
	go func() {
		f()
		d.Hide()
	}()
	d.Show()
}

func (v *mainView) thumbnailList() {
	if v.imageView.GIFImage == nil {
		return
	}
	wf := func() fyne.Window { return newThumbnailListWindow(v.imageView.GIFImage) }
	v.openWindow(wf, &v.isThumbnailListWindowOpening)
}

func (v *mainView) info() {
	if v.imageView.GIFImage == nil {
		return
	}
	wf := func() fyne.Window { return newInfoWindow(v.imageView.GIFImage) }
	v.openWindow(wf, &v.isInfoWindowOpening)
}

func (v *mainView) credits() {
	wf := func() fyne.Window { return CreditsWindow(fyne.CurrentApp(), fyne.NewSize(800, 400)) }
	v.openWindow(wf, &v.isCreditsWindowOpening)
}

func (v *mainView) about() {
	v.openWindow(newAboutWindow, &v.isAboutWindowOpening)
}

func (v *mainView) openWindow(windowOpen func() fyne.Window, isOpening *bool) {
	if *isOpening {
		return
	}
	w := windowOpen()
	w.SetOnClosed(func() { *isOpening = false })
	w.Show()
	*isOpening = true
}

func (v *mainView) zoomIn() {
	if v.imageView.GIFImage == nil {
		return
	}
	v.imageView.zoomIn()
	v.navigateBar.update()
}

func (v *mainView) zoomOut() {
	if v.imageView.GIFImage == nil {
		return
	}
	v.imageView.zoomOut()
	v.navigateBar.update()
}

func (v *mainView) handleKeys(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyLeft:
		v.prev()
	case fyne.KeyRight:
		v.next()
	case fyne.KeyUp:
		v.first()
	case fyne.KeyDown:
		v.last()
	case fyne.KeySpace:
		v.pressPlayButton()
	}
}

func (v *mainView) handleRune(r rune) {
	switch r {
	case '+':
		v.zoomIn()
	case '-':
		v.zoomOut()
	case '[':
		v.decreaseSpeed()
	case ']':
		v.increaseSpeed()
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
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	v := newMainView(w)
	w.SetContent(v.Container)
	w.Canvas().SetOnTypedKey(v.handleKeys)
	w.Canvas().SetOnTypedRune(v.handleRune)
	v.addShortcuts()
	if len(args) > 1 {
		v.withLoadingDialog(func() {
			err := v.loadImageFromPath(args[1])
			if err != nil {
				dialog.ShowError(err, w)
			}
		})
	}
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
