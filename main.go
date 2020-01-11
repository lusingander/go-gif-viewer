package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
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
	viewArea := widget.NewScrollContainer(img.Get(0))
	viewArea.Resize(defaultWindowSize)
	navigateBar := newNavigateBar(img.Length())
	panel := fyne.NewContainerWithLayout(layout.NewBorderLayout(
		navigateBar.bar, nil, nil, nil), navigateBar.bar, viewArea)
	w.SetContent(panel)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
