package layout

import (
	"testing"
)

func TestBoundingBoxContains(t *testing.T) {
	box := BoundingBox{X: 10, Y: 20, Width: 100, Height: 50}

	if !box.Contains(15, 25) {
		t.Errorf("expected (15, 25) to be inside box")
	}
	if !box.Contains(10, 20) {
		t.Errorf("expected top-left boundary to be inside box")
	}
	if !box.Contains(110, 70) {
		t.Errorf("expected bottom-right boundary to be inside box")
	}
	if box.Contains(5, 25) || box.Contains(15, 75) {
		t.Errorf("expected out-of-bounds coordinates to return false")
	}
}

func TestBoundingBoxInnerBox(t *testing.T) {
	box := BoundingBox{X: 0, Y: 0, Width: 100, Height: 100}
	insets := EdgeInsets{Top: 10, Right: 15, Bottom: 10, Left: 15}

	inner := box.InnerBox(insets)

	if inner.X != 15 || inner.Y != 10 {
		t.Errorf("expected inner origin (15, 10), got (%v, %v)", inner.X, inner.Y)
	}
	if inner.Width != 70 || inner.Height != 80 {
		t.Errorf("expected inner size (70, 80), got (%v, %v)", inner.Width, inner.Height)
	}
}
