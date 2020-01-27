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
	*imageView
	*navigateBar
}

func newMainView() *mainView {
	imageView := newImageView()
	navigateBar := newNavigateBar()
	panel := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, navigateBar.CanvasObject, nil, nil),
		navigateBar.CanvasObject, imageView.CanvasObject,
	)
	return &mainView{
		Container:   panel,
		imageView:   imageView,
		navigateBar: navigateBar,
	}
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
}

func (v *mainView) clearImage() {
	v.imageView.clearImage()
	v.navigateBar.clearImage()
}

func (v *mainView) openFileDialog() {
	// TODO: https://github.com/sqweek/dialog/issues/24
	f, err := dialog.File().Filter("GIF", "gif").Load()
	if err != nil {
		return
	}
	v.loadImageFromPath(f)
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
	))
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
