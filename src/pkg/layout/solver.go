package layout

// LayoutNode defines the structural interface required by the layout solver.
type LayoutNode interface {
	GetBox() *BoundingBox
	GetChildren() []LayoutNode
	GetText() string
	IsText() bool
}

// Solver calculates concrete pixel dimensions and positions for layout nodes.
type Solver interface {
	Solve(root LayoutNode, availableWidth, availableHeight float64) error
}
