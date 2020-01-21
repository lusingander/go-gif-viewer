package image

import (
	"image"
	"image/color"
	"image/gif"
	"os"

	"fyne.io/fyne"
)

type GIFImage struct {
	origin *gif.GIF
	images []image.Image
}

func (g *GIFImage) Size() fyne.Size {
	return getImageSize(g.origin)
}

func (g *GIFImage) Get(i int) image.Image {
	return g.images[i]
}

func (g *GIFImage) Length() int {
	return len(g.images)
}

func (g *GIFImage) DelayMilliSec() []int {
	delay := make([]int, g.Length())
	for i, d := range g.origin.Delay {
		delay[i] = d * 10
	}
	return delay
}

func LoadGIFImageFromPath(path string) (*GIFImage, error) {
	g, err := loadGIF(path)
	if err != nil {
		return nil, err
	}
	return newGIFImage(g)
}

func newGIFImage(g *gif.GIF) (*GIFImage, error) {
	rect := g.Image[0].Rect
	images := make([]image.Image, len(g.Image))
	images[0] = g.Image[0]
	for i := 1; i < len(g.Image); i++ {
		p := g.Image[i]
		img := image.NewRGBA(rect)
		for x := 0; x < rect.Dx(); x++ {
			for y := 0; y < rect.Dy(); y++ {
				if isInRect(x, y, p.Rect) && isOpaque(p.At(x, y)) {
					img.Set(x, y, p.At(x, y))
				} else {
					img.Set(x, y, images[i-1].At(x, y))
				}
			}
		}
		images[i] = img
	}
	return &GIFImage{g, images}, nil
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

func isOpaque(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a > 0
}

func isInRect(x, y int, r image.Rectangle) bool {
	return r.Min.X <= x && x < r.Max.X &&
		r.Min.Y <= y && y < r.Max.Y
}
