package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
	"github.com/nfnt/resize"
)

type imageView struct {
	*canvas.Image
	fyne.CanvasObject

	*image.GIFImage

	scale float64
}

func newImageView() *imageView {
	image := &canvas.Image{
		FillMode: canvas.ImageFillOriginal,
	}
	imageBox := widget.NewVBox(
		layout.NewSpacer(),
		widget.NewHBox(
			layout.NewSpacer(),
			image,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	canvas := widget.NewScrollContainer(imageBox)
	canvas.Resize(defaultWindowSize)
	return &imageView{
		Image:        image,
		CanvasObject: canvas,
		scale:        1.0,
	}
}

func (v *imageView) zoomIn() {
	if v.scale < 2.0 {
		v.scale += 0.1
	}
}

func (v *imageView) zoomOut() {
	if v.scale > 0.2 {
		v.scale -= 0.1
	}
}

func (v *imageView) setImage(img *image.GIFImage) {
	v.GIFImage = img
	v.scale = 1.0
	v.refleshFrame(0)
}

func (v *imageView) clearImage() {
	v.Image.Image = nil
	v.reflesh()
}

func (v *imageView) refleshFrame(n int) {
	img := v.GIFImage.Get(n)
	w, h := v.scaledImageSize()
	v.Image.Image = resize.Resize(w, h, img, resize.Bilinear)
	v.reflesh()
}

func (v *imageView) reflesh() {
	canvas.Refresh(v.Image)
}

func (v *imageView) scaledImageSize() (uint, uint) {
	size := v.GIFImage.Size()
	w := float64(size.Width) * v.scale
	h := float64(size.Height) * v.scale
	return uint(w), uint(h)
}
