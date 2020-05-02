package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lusingander/go-gif-viewer/image"
)

const (
	infoWindowName = "Info"
)

func newInfoWindow(gif *image.GIFImage) fyne.Window {
	w := fyne.CurrentApp().NewWindow(infoWindowName)
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			fyne.NewContainerWithLayout(
				layout.NewVBoxLayout(),
				keyLabel("File name"),
				keyLabel("File size"),
				keyLabel("Image size"),
				keyLabel("Frame count"),
			),
			fyne.NewContainerWithLayout(
				layout.NewVBoxLayout(),
				valueLabel(gif.FileName()),
				valueLabel(formatFileSize(gif)),
				valueLabel(formatImageSize(gif)),
				valueLabel(gif.Length()),
			),
		),
	)
	return w
}

func formatFileSize(gif *image.GIFImage) string {
	return fmt.Sprintf("%d bytes", gif.FileSizeByte())
}

func formatImageSize(gif *image.GIFImage) string {
	w, h := gif.Size()
	return fmt.Sprintf("%d x %d", w, h)
}

func keyLabel(l string) *widget.Label {
	return widget.NewLabelWithStyle(l+":", fyne.TextAlignTrailing, fyne.TextStyle{})
}

func valueLabel(v interface{}) *widget.Label {
	return widget.NewLabelWithStyle(fmt.Sprintf("%v", v), fyne.TextAlignLeading, fyne.TextStyle{})
}
