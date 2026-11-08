package dom

import (
	"testing"
)

func TestNodeHierarchyAndAttributes(t *testing.T) {
	root := NewElement("div", map[string]string{
		"id":    "root-container",
		"class": "container screen-root",
	})

	if root.TagName != "div" {
		t.Fatalf("expected tag name 'div', got '%s'", root.TagName)
	}

	if !root.HasClass("container") || !root.HasClass("screen-root") {
		t.Fatalf("expected classes 'container' and 'screen-root' to be present")
	}

	if root.HasClass("missing-class") {
		t.Fatalf("expected missing-class to be false")
	}

	child := NewElement("span", map[string]string{"id": "status-label"})
	textNode := NewText("Locked")
	child.AppendChild(textNode)
	root.AppendChild(child)

	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child in root, got %d", len(root.Children))
	}

	found := root.FindByID("status-label")
	if found == nil {
		t.Fatalf("expected to find child with ID 'status-label'")
	}

	if len(found.Children) != 1 || !found.Children[0].IsText() {
		t.Fatalf("expected text child under found node")
	}

	if found.Children[0].GetText() != "Locked" {
		t.Fatalf("expected text 'Locked', got '%s'", found.Children[0].GetText())
	}
}
