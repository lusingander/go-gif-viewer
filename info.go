package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	infoWindowName = "Info"
)

var isInfoWindowOpening bool

func showInfoWindow(gif *image.GIFImage) {
	if isInfoWindowOpening {
		return
	}
	w := fyne.CurrentApp().NewWindow(infoWindowName)
	w.SetContent(
		widget.NewVBox(
			label("File name", gif.FileName()),
			label("File size", formatFileSize(gif)),
			label("Image size", formatImageSize(gif)),
			label("Frame count", gif.Length()),
		),
	)
	w.SetOnClosed(func() { isInfoWindowOpening = false })
	w.Show()
	isInfoWindowOpening = true
}

func formatFileSize(gif *image.GIFImage) string {
	return fmt.Sprintf("%d bytes", gif.FileSizeByte())
}

func formatImageSize(gif *image.GIFImage) string {
	w, h := gif.Size()
	return fmt.Sprintf("%d x %d", w, h)
}

func label(k string, v interface{}) *widget.Label {
	return widget.NewLabel(fmt.Sprintf("%s: %v", k, v))
}
