package dom

import (
	"strings"

	"io.tylr.chrono/pkg/layout"
	"io.tylr.chrono/pkg/style"
)

// NodeType identifies the semantic type of a DOM node.
type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
	CommentNode
	DocumentNode
)

// Node represents a single element or text entity in the Chrono DOM tree.
type Node struct {
	Type       NodeType
	TagName    string
	Attributes map[string]string
	Children   []*Node
	Parent     *Node
	Text       string

	// Layout and visual attributes resolved during pipeline execution
	Style *style.ComputedStyle
	Box   *layout.BoundingBox
}

// NewElement creates a new element node with initialized attributes and children.
func NewElement(tagName string, attrs map[string]string) *Node {
	normalizedAttrs := make(map[string]string)
	for k, v := range attrs {
		normalizedAttrs[strings.ToLower(k)] = v
	}

	return &Node{
		Type:       ElementNode,
		TagName:    strings.ToLower(tagName),
		Attributes: normalizedAttrs,
		Children:   make([]*Node, 0),
		Style:      style.NewComputedStyle(),
		Box:        &layout.BoundingBox{},
	}
}

// NewText creates a new text node with the given content.
func NewText(content string) *Node {
	return &Node{
		Type:  TextNode,
		Text:  content,
		Style: style.NewComputedStyle(),
		Box:   &layout.BoundingBox{},
	}
}

// AppendChild appends a child node to this node and updates the parent pointer.
func (n *Node) AppendChild(child *Node) {
	if child == nil {
		return
	}
	child.Parent = n
	n.Children = append(n.Children, child)
}

// GetAttribute returns the value of an attribute and whether it was present.
func (n *Node) GetAttribute(key string) (string, bool) {
	if n.Attributes == nil {
		return "", false
	}
	val, ok := n.Attributes[strings.ToLower(key)]
	return val, ok
}

// SetAttribute sets an attribute value on the node.
func (n *Node) SetAttribute(key, value string) {
	if n.Attributes == nil {
		n.Attributes = make(map[string]string)
	}
	n.Attributes[strings.ToLower(key)] = value
}

// HasAttribute checks if an attribute key exists on the node.
func (n *Node) HasAttribute(key string) bool {
	if n.Attributes == nil {
		return false
	}
	_, ok := n.Attributes[strings.ToLower(key)]
	return ok
}

// HasClass checks if the node has a specific class name.
func (n *Node) HasClass(className string) bool {
	classAttr, ok := n.GetAttribute("class")
	if !ok {
		return false
	}
	for _, c := range strings.Fields(classAttr) {
		if c == className {
			return true
		}
	}
	return false
}

// ID returns the element's id attribute, if present.
func (n *Node) ID() string {
	val, _ := n.GetAttribute("id")
	return val
}

// FindByID searches the subtree for a node with the specified ID.
func (n *Node) FindByID(id string) *Node {
	if n.ID() == id {
		return n
	}
	for _, child := range n.Children {
		if found := child.FindByID(id); found != nil {
			return found
		}
	}
	return nil
}

// LayoutNode interface implementation

// GetBox returns the node's bounding box.
func (n *Node) GetBox() *layout.BoundingBox {
	return n.Box
}

// GetChildren returns children as slice of LayoutNode.
func (n *Node) GetChildren() []layout.LayoutNode {
	nodes := make([]layout.LayoutNode, len(n.Children))
	for i, c := range n.Children {
		nodes[i] = c
	}
	return nodes
}

// GetText returns the node's text content.
func (n *Node) GetText() string {
	return n.Text
}

// IsText returns true if this is a text node.
func (n *Node) IsText() bool {
	return n.Type == TextNode
}

// GetComputedStyle returns the node's computed style.
func (n *Node) GetComputedStyle() *style.ComputedStyle {
	return n.Style
}

// GetTagName returns the tag name of the node.
func (n *Node) GetTagName() string {
	return n.TagName
}
