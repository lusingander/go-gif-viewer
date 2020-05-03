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

func thumbnails(gif *image.GIFImage) []fyne.CanvasObject {
	ts := make([]fyne.CanvasObject, 0)
	size := uint(thumbnailListDefaultWidth)
	for i := 0; i < gif.Length(); i++ {
		resized := resize.Resize(size, size, gif.Get(i), resize.Bilinear)
		img := &canvas.Image{
			Image:    resized,
			FillMode: canvas.ImageFillOriginal,
		}
		ts = append(ts, thumbnailContainer(img, i))
	}
	return ts
}

func newThumbnailListWindow(gif *image.GIFImage) fyne.Window {
	w := fyne.CurrentApp().NewWindow(thumbnailListWindowName)
	w.SetContent(
		widget.NewVScrollContainer(
			fyne.NewContainerWithLayout(
				layout.NewVBoxLayout(),
				thumbnails(gif)...,
			),
		),
	)
	w.Resize(thumbnailListDefaultSize)
	return w
}
