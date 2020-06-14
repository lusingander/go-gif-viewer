package main

import (
	im "image"

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

	caches []im.Image
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
	return &imageView{
		Image:        image,
		CanvasObject: canvas,
		scale:        scaleDefault,
	}
}

func (v *imageView) zoomIn() {
	if v.scale < scaleMax {
		v.scale += scaleStep
		v.loadImageCaches()
	}
}

func (v *imageView) zoomOut() {
	if v.scale > scaleMin {
		v.scale -= scaleStep
		v.loadImageCaches()
	}
}

func (v *imageView) setImage(img *image.GIFImage) {
	v.GIFImage = img
	v.scale = scaleDefault
	v.loadImageCaches()
	v.refleshFrame(0)
}

func (v *imageView) loadImageCaches() {
	l := v.GIFImage.Length()
	caches := make([]im.Image, l)
	for i := 0; i < l; i++ {
		img := v.GIFImage.Get(i)
		w, h := v.scaledImageSize()
		caches[i] = resize.Resize(w, h, img, resize.NearestNeighbor)
	}
	v.caches = caches
}

func (v *imageView) clearImage() {
	v.Image.Image = nil
	v.reflesh()
}

func (v *imageView) refleshFrame(n int) {
	v.Image.Image = v.caches[n]
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
