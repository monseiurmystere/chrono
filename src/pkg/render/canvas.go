package render

import (
	"io.tylr.chrono/pkg/layout"
	"io.tylr.chrono/pkg/style"
)

// Canvas represents a 2D pixel buffer target that can be blitted to Wayland shm surfaces.
type Canvas interface {
	Width() int
	Height() int
	Clear(color style.Color)
	DrawRect(box layout.BoundingBox, fill style.Color, border style.Color, borderWidth float64, borderRadius float64)
	DrawText(text string, x, y float64, fontFamily string, size float64, weight int, color style.Color)
	Buffer() []byte
}

// RenderableNode defines the node interface required by the painter.
type RenderableNode interface {
	GetBox() *layout.BoundingBox
	GetComputedStyle() *style.ComputedStyle
	GetChildren() []RenderableNode
	GetText() string
	IsText() bool
	TagName() string
}

// Painter renders a tree of renderable nodes to a target canvas.
type Painter interface {
	Paint(root RenderableNode, canvas Canvas) error
}
