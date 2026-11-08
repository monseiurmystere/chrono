package style

import (
	"fmt"
	"strconv"
	"strings"
)

// DisplayMode defines how an element is formatted in the layout.
type DisplayMode string

const (
	DisplayFlex  DisplayMode = "flex"
	DisplayBlock DisplayMode = "block"
	DisplayNone  DisplayMode = "none"
)

// FlexDirection defines the main axis direction of a flex container.
type FlexDirection string

const (
	FlexRow            FlexDirection = "row"
	FlexColumn         FlexDirection = "column"
	FlexRowReverse     FlexDirection = "row-reverse"
	FlexColumnReverse  FlexDirection = "column-reverse"
)

// JustifyContent aligns items along the main axis.
type JustifyContent string

const (
	JustifyFlexStart    JustifyContent = "flex-start"
	JustifyFlexEnd      JustifyContent = "flex-end"
	JustifyCenter       JustifyContent = "center"
	JustifySpaceBetween JustifyContent = "space-between"
	JustifySpaceAround  JustifyContent = "space-around"
	JustifySpaceEvenly  JustifyContent = "space-evenly"
)

// AlignItems aligns items along the cross axis.
type AlignItems string

const (
	AlignFlexStart AlignItems = "flex-start"
	AlignFlexEnd   AlignItems = "flex-end"
	AlignCenter    AlignItems = "center"
	AlignStretch   AlignItems = "stretch"
)

// Color represents an RGBA color with channels from 0 to 255 and alpha 0.0 to 1.0.
type Color struct {
	R uint8
	G uint8
	B uint8
	A float64
}

// Hex creates a Color from a hex string (#RGB, #RRGGBB, #RRGGBBAA).
func ParseHexColor(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	switch len(hex) {
	case 3:
		r, err := strconv.ParseUint(string([]byte{hex[0], hex[0]}), 16, 8)
		if err != nil {
			return Color{}, err
		}
		g, err := strconv.ParseUint(string([]byte{hex[1], hex[1]}), 16, 8)
		if err != nil {
			return Color{}, err
		}
		b, err := strconv.ParseUint(string([]byte{hex[2], hex[2]}), 16, 8)
		if err != nil {
			return Color{}, err
		}
		return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 1.0}, nil
	case 6:
		rgb, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return Color{}, err
		}
		return Color{
			R: uint8((rgb >> 16) & 0xFF),
			G: uint8((rgb >> 8) & 0xFF),
			B: uint8(rgb & 0xFF),
			A: 1.0,
		}, nil
	case 8:
		rgba, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return Color{}, err
		}
		return Color{
			R: uint8((rgba >> 24) & 0xFF),
			G: uint8((rgba >> 16) & 0xFF),
			B: uint8((rgba >> 8) & 0xFF),
			A: float64(rgba&0xFF) / 255.0,
		}, nil
	default:
		return Color{}, fmt.Errorf("invalid hex color length: %s", hex)
	}
}

// Transparent returns a fully transparent color.
func Transparent() Color {
	return Color{R: 0, G: 0, B: 0, A: 0.0}
}

// White returns solid white.
func White() Color {
	return Color{R: 255, G: 255, B: 255, A: 1.0}
}

// Black returns solid black.
func Black() Color {
	return Color{R: 0, G: 0, B: 0, A: 1.0}
}

// ComputedStyle represents the fully resolved styling attributes for a DOM node.
type ComputedStyle struct {
	// Layout & Flexbox
	Display        DisplayMode
	FlexDirection  FlexDirection
	JustifyContent JustifyContent
	AlignItems     AlignItems
	Gap            float64
	FlexGrow       float64
	FlexShrink     float64
	FlexBasis      float64

	// Box dimensions (0 = auto / unconstrained)
	Width     float64
	Height    float64
	MinWidth  float64
	MaxWidth  float64
	MinHeight float64
	MaxHeight float64

	// Box Spacing
	MarginTop    float64
	MarginRight  float64
	MarginBottom float64
	MarginLeft   float64

	PaddingTop    float64
	PaddingRight  float64
	PaddingBottom float64
	PaddingLeft   float64

	// Visual styling
	Color           Color
	BackgroundColor Color
	BorderColor     Color
	BorderWidth     float64
	BorderRadius    float64
	Opacity         float64

	// Typography
	FontFamily string
	FontSize   float64
	FontWeight int
	TextAlign  string
}

// NewComputedStyle returns a ComputedStyle initialized with sensible default values.
func NewComputedStyle() *ComputedStyle {
	return &ComputedStyle{
		Display:         DisplayBlock,
		FlexDirection:   FlexRow,
		JustifyContent:  JustifyFlexStart,
		AlignItems:      AlignStretch,
		FlexShrink:      1.0,
		Opacity:         1.0,
		Color:           White(),
		BackgroundColor: Transparent(),
		BorderColor:     Transparent(),
		FontFamily:      "sans-serif",
		FontSize:        16.0,
		FontWeight:      400,
		TextAlign:       "left",
	}
}

// Specificity represents the specificity weight of a CSS selector.
type Specificity struct {
	IDs        int
	Classes    int
	Elements   int
}

// GreaterThan checks if this specificity is greater than another.
func (s Specificity) GreaterThan(other Specificity) bool {
	if s.IDs != other.IDs {
		return s.IDs > other.IDs
	}
	if s.Classes != other.Classes {
		return s.Classes > other.Classes
	}
	return s.Elements > other.Elements
}

// StyleRule represents a single parsed CSS style rule.
type StyleRule struct {
	Selector    string
	Specificity Specificity
	Properties  map[string]string
}
