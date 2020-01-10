package main

import (
	"bytes"
	"image/gif"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

const filePath = "./sample.gif"

func loadGIF(path string) (*gif.GIF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return gif.DecodeAll(f)
}

func loadGIFImageFromPath(path string) (*canvas.Image, error) {
	g, err := loadGIF(path)
	if err != nil {
		return nil, err
	}
	_ = g
	buf := new(bytes.Buffer)
	err = gif.Encode(buf, g.Image[0], nil)
	if err != nil {
		return nil, err
	}
	res := fyne.NewStaticResource("res", buf.Bytes())
	return canvas.NewImageFromResource(res), nil
}

var defaultWindowSize = fyne.NewSize(400, 400)

func run(args []string) error {
	a := app.New()
	w := a.NewWindow("GIF Viewer")
	w.Resize(defaultWindowSize)
	img, err := loadGIFImageFromPath(args[1])
	if err != nil {
		return err
	}
	img.SetMinSize(fyne.NewSize(600, 600))
	scroll := widget.NewScrollContainer(img)
	scroll.Resize(defaultWindowSize)
	w.SetContent(scroll)
	w.ShowAndRun()
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
