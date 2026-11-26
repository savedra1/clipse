package app

import (
	"image"
	"testing"

	"github.com/savedra1/clipse/config"
)

func TestSmartResizeSmallWindow(t *testing.T) {
	config.ClipseConfig.ImageDisplay.ScaleX = 9
	config.ClipseConfig.ImageDisplay.ScaleY = 9
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// A preview area with no height must not panic.
	smartResize(img, 146, 0)

	// A tall box and a wide image: the width is the limit.
	got := smartResize(img, 15, 10).Bounds()
	if got.Dx() > 135 || got.Dy() > 90 {
		t.Errorf("resized to %dx%d, want it within 135x90", got.Dx(), got.Dy())
	}
}
