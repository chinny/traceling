package lighten

import (
	"image"
	"image/color"
	"testing"
)

func TestToTraceableRemapsBlackAndKeepsWhite(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	src.Set(0, 0, color.RGBA{0, 0, 0, 255})       // black line
	src.Set(1, 0, color.RGBA{255, 255, 255, 255}) // white paper

	got := ToTraceable(src, 210)

	if v := got.GrayAt(0, 0).Y; v != 210 {
		t.Errorf("black pixel = %d, want 210", v)
	}
	if v := got.GrayAt(1, 0).Y; v != 255 {
		t.Errorf("white pixel = %d, want 255", v)
	}
}

func TestToTraceableIsMonotonic(t *testing.T) {
	lut := buildLUT(210)
	for v := 1; v < 256; v++ {
		if lut[v] < lut[v-1] {
			t.Fatalf("lut not monotonic at %d: %d < %d", v, lut[v], lut[v-1])
		}
	}
}

func TestToTraceableGenericPath(t *testing.T) {
	// Non-RGBA source exercises the fallback path.
	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.NRGBA{0, 0, 0, 255})

	if v := ToTraceable(src, 180).GrayAt(0, 0).Y; v != 180 {
		t.Errorf("black pixel via generic path = %d, want 180", v)
	}
}
