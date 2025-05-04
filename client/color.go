package main

import (
	"fmt"
	"math"
)

// Gets a color code from a string.
// Used for coloring in chat messages in the color of the username.
func calculateColorCode(s string) string {
	hash := 0
	for i := range len(s) {
		hash = int([]rune(s)[i]) + ((hash << 5) - hash)
	}

	hue := float64(hash % 360)

	// We calculate in hsl first in order to make the colors not too dim.
	sat := 1.0
	light := 0.5

	r, g, b := hslToRgb(hue, sat, light)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func hslToRgb(h, s, l float64) (int, int, int) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r1, g1, b1 float64

	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	r := int((r1 + m) * 255)
	g := int((g1 + m) * 255)
	b := int((b1 + m) * 255)

	return r, g, b
}
