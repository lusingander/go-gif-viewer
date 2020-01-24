package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

func run(args []string) error {
	a := app.New()
	w := a.NewWindow(appName)
	w.Resize(defaultWindowSize)
	img, err := image.LoadGIFImageFromPath(args[1])
	if err != nil {
		return err
	}
	imgArea := &canvas.Image{Image: img.Get(0), FillMode: canvas.ImageFillContain}
	viewArea := widget.NewScrollContainer(imgArea)
	viewArea.Resize(defaultWindowSize)
	navigateBar := newNavigateBar()
	navigateBar.setImage(img)
	navigateBar.addObserver(func(n int) {
		imgArea.Image = img.Get(n)
		canvas.Refresh(imgArea)
	})
	panel := fyne.NewContainerWithLayout(layout.NewBorderLayout(
		navigateBar.CanvasObject, nil, nil, nil), navigateBar.CanvasObject, viewArea)
	w.SetContent(panel)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
