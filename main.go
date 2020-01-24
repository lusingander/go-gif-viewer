package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

func loadImage(img *image.GIFImage, v *imageView, b *navigateBar) {
	v.setImage(img)
	b.setImage(img)
	// TODO: remove old observer
	b.addObserver(func(n int) {
		v.Image.Image = img.Get(n)
		canvas.Refresh(v.Image)
	})
}

func run(args []string) error {
	a := app.New()
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	imageView := newImageView()
	navigateBar := newNavigateBar()
	panel := fyne.NewContainerWithLayout(layout.NewBorderLayout(
		navigateBar.CanvasObject, nil, nil, nil), navigateBar.CanvasObject, imageView.CanvasObject)
	img, err := image.LoadGIFImageFromPath(args[1])
	if err != nil {
		return err
	}
	loadImage(img, imageView, navigateBar)
	w.SetContent(panel)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
