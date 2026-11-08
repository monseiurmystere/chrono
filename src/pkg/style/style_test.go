package style

import (
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		hex      string
		expected Color
		hasError bool
	}{
		{"#fff", Color{R: 255, G: 255, B: 255, A: 1.0}, false},
		{"#000000", Color{R: 0, G: 0, B: 0, A: 1.0}, false},
		{"#58a6ff", Color{R: 88, G: 166, B: 255, A: 1.0}, false},
		{"#ff000080", Color{R: 255, G: 0, B: 0, A: 128.0 / 255.0}, false},
		{"invalid", Color{}, true},
		{"#12", Color{}, true},
	}

	for _, tc := range tests {
		c, err := ParseHexColor(tc.hex)
		if tc.hasError && err == nil {
			t.Errorf("expected error for '%s', got nil", tc.hex)
		}
		if !tc.hasError && err != nil {
			t.Errorf("unexpected error for '%s': %v", tc.hex, err)
		}
		if !tc.hasError && (c.R != tc.expected.R || c.G != tc.expected.G || c.B != tc.expected.B) {
			t.Errorf("expected %v for '%s', got %v", tc.expected, tc.hex, c)
		}
	}
}

func TestSpecificityComparison(t *testing.T) {
	idRule := Specificity{IDs: 1, Classes: 0, Elements: 0}
	classRule := Specificity{IDs: 0, Classes: 2, Elements: 5}
	elemRule := Specificity{IDs: 0, Classes: 0, Elements: 1}

	if !idRule.GreaterThan(classRule) {
		t.Errorf("expected ID specificity to outweigh classes")
	}
	if !classRule.GreaterThan(elemRule) {
		t.Errorf("expected class specificity to outweigh element")
	}
	if elemRule.GreaterThan(classRule) {
		t.Errorf("expected element specificity not to outweigh classes")
	}
}
