package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"github.com/lusingander/go-gif-viewer/image"
	"github.com/sqweek/dialog"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

func loadImage(img *image.GIFImage, v *imageView, b *navigateBar) {
	v.setImage(img)
	b.setImage(img)
	b.addObserver(v.refleshFrame)
}

func run(args []string) error {
	a := app.New()
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	imageView := newImageView()
	navigateBar := newNavigateBar()
	panel := fyne.NewContainerWithLayout(layout.NewBorderLayout(
		navigateBar.CanvasObject, nil, nil, nil), navigateBar.CanvasObject, imageView.CanvasObject)
	w.SetContent(panel)
	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", func() {
				// TODO: refactoring
				f, err := dialog.File().Filter("GIF", "gif").Load()
				if err != nil {
					return
				}
				img, err := image.LoadGIFImageFromPath(f)
				if err != nil {
					return
				}
				loadImage(img, imageView, navigateBar)
			}),
		),
	))
	if len(args) > 1 {
		// TODO: refactoring
		img, err := image.LoadGIFImageFromPath(args[1])
		if err != nil {
			return err
		}
		loadImage(img, imageView, navigateBar)
	}
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
