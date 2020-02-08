package image

import (
	"testing"
)

func TestLoadGif(t *testing.T) {
	_, err := LoadGIFImageFromPath("../resource/sample/sample.gif")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}
}
