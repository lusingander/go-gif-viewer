package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
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
	scroll := widget.NewScrollContainer(img.Get(0))
	scroll.Resize(defaultWindowSize)
	w.SetContent(scroll)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
