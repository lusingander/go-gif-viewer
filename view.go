package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

type imageView struct {
	*canvas.Image
	fyne.CanvasObject
}

func newImageView() *imageView {
	image := &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	canvas := widget.NewScrollContainer(image)
	canvas.Resize(defaultWindowSize)
	return &imageView{
		Image:        image,
		CanvasObject: canvas,
	}
}

func (v *imageView) setImage(img *image.GIFImage) {
	v.Image.Image = img.Get(0)
}
