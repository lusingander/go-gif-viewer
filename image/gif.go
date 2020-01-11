package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type GIFImage struct {
	origin *gif.GIF
	images []*canvas.Image
}

func (g *GIFImage) Size() fyne.Size {
	return getImageSize(g.origin)
}

func (g *GIFImage) Get(i int) *canvas.Image {
	return g.images[i]
}

func (g *GIFImage) Length() int {
	return len(g.images)
}

func LoadGIFImageFromPath(path string) (*GIFImage, error) {
	g, err := loadGIF(path)
	if err != nil {
		return nil, err
	}
	return newGIFImage(g)
}

func newGIFImage(g *gif.GIF) (*GIFImage, error) {
	size := getImageSize(g)
	images := make([]*canvas.Image, 0, len(g.Image))
	for i, base := range g.Image {
		image, err := newImage(base, fmt.Sprintf("res_%d", i))
		if err != nil {
			return nil, err
		}
		image.SetMinSize(size)
		images = append(images, image)
	}
	return &GIFImage{g, images}, nil
}

func newImage(src image.Image, name string) (*canvas.Image, error) {
	buf := new(bytes.Buffer)
	if err := gif.Encode(buf, src, nil); err != nil {
		return nil, err
	}
	res := fyne.NewStaticResource(name, buf.Bytes())
	return canvas.NewImageFromResource(res), nil
}

func loadGIF(path string) (*gif.GIF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return gif.DecodeAll(f)
}

func getImageSize(g *gif.GIF) fyne.Size {
	s := g.Image[0].Rect.Size()
	return fyne.NewSize(s.X, s.Y)
}
