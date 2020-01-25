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

	*image.GIFImage
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
	v.GIFImage = img
	v.refleshFrame(0)
}

func (v *imageView) refleshFrame(n int) {
	v.Image.Image = v.GIFImage.Get(n)
	v.reflesh()
}

func (v *imageView) reflesh() {
	canvas.Refresh(v.Image)
}
