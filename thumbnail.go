package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
	"github.com/nfnt/resize"
)

const (
	thumbnailListWindowName = "List"
)

var (
	thumbnailListDefaultWidth  = 150
	thumbnailListDefaultHeight = 400
	thumbnailListDefaultSize   = fyne.NewSize(thumbnailListDefaultWidth, thumbnailListDefaultHeight)
)

func thumbnailContainer(img *canvas.Image, i int) fyne.CanvasObject {
	label := fmt.Sprintf("%d", i+1)
	return widget.NewGroup(label, img)
}

func addThumbnails(container *fyne.Container, gif *image.GIFImage) {
	size := uint(thumbnailListDefaultWidth)
	for i := 0; i < gif.Length(); i++ {
		resized := resize.Resize(size, size, gif.Get(i), resize.Bilinear)
		img := &canvas.Image{
			Image:    resized,
			FillMode: canvas.ImageFillOriginal,
		}
		container.AddObject(thumbnailContainer(img, i))
	}
}

func newThumbnailListWindow(gif *image.GIFImage) fyne.Window {
	w := fyne.CurrentApp().NewWindow(thumbnailListWindowName)
	l := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	w.SetContent(widget.NewVScrollContainer(l))
	w.Resize(thumbnailListDefaultSize)
	go addThumbnails(l, gif)
	return w
}
