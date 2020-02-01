package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
	"github.com/nfnt/resize"
)

const (
	scaleDefault = 1.0
	scaleMax     = 2.0 - 0.01
	scaleMin     = 0.1 + 0.01
	scaleStep    = 0.1
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
		scale:        scaleDefault,
	}
}

func (v *imageView) zoomIn() {
	if v.scale < scaleMax {
		v.scale += scaleStep
	}
}

func (v *imageView) zoomOut() {
	if v.scale > scaleMin {
		v.scale -= scaleStep
	}
}

func (v *imageView) setImage(img *image.GIFImage) {
	v.GIFImage = img
	v.scale = scaleDefault
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
	w, h := v.GIFImage.Size()
	fw := float64(w) * v.scale
	fh := float64(h) * v.scale
	return uint(fw), uint(fh)
}
