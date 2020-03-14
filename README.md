go-gif-viewer
====

Simple animated GIF viewer with Go and [Fyne](https://fyne.io/)

<img src="./resource/document/screenshot.gif" width=400>

## Installation

`$ go get github.com/lusingander/go-gif-viewer`

## Usage

`$ go-gif-viewer sample.gif`

or

`$ go-gif-viewer` and select open file icon (<img src="./resource/icons/svg/open.svg" width=16>).

### Keybindings

|Key|Description|Icon|
|-|-|-|
|←|Go to previous frame|<img src="./resource/icons/svg/prev.svg" width=16>|
|→|Go to next frame|<img src="./resource/icons/svg/next.svg" width=16>|
|↑|Go to first frame|<img src="./resource/icons/svg/first.svg" width=16>|
|↓|Go to last frame|<img src="./resource/icons/svg/last.svg" width=16>|
|Space|Play / Pause|<img src="./resource/icons/svg/play.svg" width=16> / <img src="./resource/icons/svg/pause.svg" width=16>|
||||
|+|Zoom in|<img src="./resource/document/zoom-in.svg" width=16>|
|-|Zoom out|<img src="./resource/document/zoom-out.svg" width=16>|
|[|Decrease playback speed|-|
|]|Increase playback speed|-|
||||
|⌘O|Open image file|<img src="./resource/icons/svg/open.svg" width=16>|
|⌘W|Close image file|<img src="./resource/icons/svg/close.svg" width=16>|

----

Sample image: By <a href="//commons.wikimedia.org/wiki/User:Marvel" title="User:Marvel">Marvel</a> - Based upon a NASA image, see <a rel="nofollow" class="external autonumber" href="http://visibleearth.nasa.gov/view_rec.php?id=2433">[1]</a>., <a href="http://creativecommons.org/licenses/by-sa/3.0/" title="Creative Commons Attribution-Share Alike 3.0">CC BY-SA 3.0</a>, <a href="https://commons.wikimedia.org/w/index.php?curid=20654992">Link</a>
