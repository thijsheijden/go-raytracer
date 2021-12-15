package color

import "image/color"

// RGB describes an RGB color using float values between 0 and 1
type RGB struct {
	R, G, B float32
}

// New creates a new RGB color
func New(R, G, B float32) RGB {
	return RGB{
		R: R,
		G: G,
		B: B,
	}
}

// RGBA converts the RGB color to RGBA color
func (c RGB) RGBA() color.RGBA {
	return color.RGBA{
		R: uint8(clamp(c.R, 0, 1) * 255),
		G: uint8(clamp(c.G, 0, 1) * 255),
		B: uint8(clamp(c.B, 0, 1) * 255),
		A: 255,
	}
}

// Mul multiplies the RGB values by v1, v2 and v3 respectively
func (c RGB) Mul(v1, v2, v3 float32) RGB {
	return RGB{
		R: c.R * v1,
		G: c.G * v2,
		B: c.B * v3,
	}
}

// Add adds two colors together
func (c RGB) Add(c2 RGB) RGB {
	return RGB{
		R: c.R + c2.R,
		G: c.G + c2.G,
		B: c.B + c2.B,
	}
}

// Average averages the color over n samples
func (c RGB) Average(nSamples int) RGB {
	return RGB{
		R: c.R / float32(nSamples),
		G: c.G / float32(nSamples),
		B: c.B / float32(nSamples),
	}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
