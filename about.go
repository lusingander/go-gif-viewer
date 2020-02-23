package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

const (
	aboutWindowName = "About"
	repositoryURL   = "https://github.com/lusingander/go-gif-viewer"
)

func newAboutWindow() fyne.Window {
	w := fyne.CurrentApp().NewWindow(aboutWindowName)
	w.SetContent(
		widget.NewVBox(
			widget.NewHBox(
				widget.NewLabelWithStyle("go-gif-viewer", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("(Version %s)", version)),
			),
			widget.NewHyperlink(repositoryURL, parseRepositoryURL()),
		),
	)
	return w
}

func parseRepositoryURL() *url.URL {
	u, _ := url.Parse(repositoryURL)
	return u
}
