package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	appName = "GIF Viewer"
)

var defaultWindowSize = fyne.NewSize(400, 400)

func createNavigateBar(n int) fyne.CanvasObject {
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {})
	slider := widget.NewSlider(0, float64(n-1))
	d := len(strconv.Itoa(n))
	count := widget.NewLabel(fmt.Sprintf("%*d/%*d", d, 1, d, n))
	buttons := widget.NewHBox(prev, next)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(
		nil, nil, buttons, count), buttons, count, slider)
}

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
	navigateBar := createNavigateBar(img.Length())
	panel := fyne.NewContainerWithLayout(layout.NewBorderLayout(
		navigateBar, nil, nil, nil), navigateBar, viewArea)
	w.SetContent(panel)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
