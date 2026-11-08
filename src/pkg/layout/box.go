package layout

// BoundingBox defines a 2D rectangular region on screen.
type BoundingBox struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Contains returns true if the coordinate is within the bounding box.
func (b BoundingBox) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

// EdgeInsets specifies insets for margins, padding, or borders.
type EdgeInsets struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// InnerBox computes the inner bounding box after subtracting insets.
func (b BoundingBox) InnerBox(insets EdgeInsets) BoundingBox {
	return BoundingBox{
		X:      b.X + insets.Left,
		Y:      b.Y + insets.Top,
		Width:  max(0, b.Width-(insets.Left+insets.Right)),
		Height: max(0, b.Height-(insets.Top+insets.Bottom)),
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
