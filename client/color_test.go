package main

import (
	"fmt"
	"regexp"
	"testing"
)

var hexColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func TestCalculateColorCode_Format(t *testing.T) {
	color := calculateColorCode("example")
	if !hexColorRe.MatchString(color) {
		t.Errorf("Color format invalid: got %s", color)
	}
}

func TestCalculateColorCode_Deterministic(t *testing.T) {
	color1 := calculateColorCode("testuser")
	color2 := calculateColorCode("testuser")
	if color1 != color2 {
		t.Errorf("Expected deterministic color, got %s and %s", color1, color2)
	}
}

func TestCalculateColorCode_Uniqueness(t *testing.T) {
	color1 := calculateColorCode("alice")
	color2 := calculateColorCode("bob")
	if color1 == color2 {
		t.Errorf("Expected different colors for different inputs: %s", color1)
	}
}

func TestCalculateColorCode_Brightness(t *testing.T) {
	color := calculateColorCode("vibrant")
	r, g, b := parseHexColor(color)

	// Compute luminance
	brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	if brightness < 90 || brightness > 200 {
		t.Errorf("Color too dark or too bright: %s (brightness %.2f)", color, brightness)
	}
}

// helper: convert hex color string to RGB ints
func parseHexColor(s string) (int, int, int) {
	var r, g, b int
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		panic("invalid color string: " + s)
	}
	return r, g, b
}
