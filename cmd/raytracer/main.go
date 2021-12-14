package main

import (
	"image"
	gocolor "image/color"
	"image/jpeg"
	"os"
	"raytracer/internal/color"
	"raytracer/internal/ray"
	"raytracer/internal/vector"

	_ "image/jpeg" // Needed for JPEG decoder
)

func main() {
	// Image config
	const aspectRatio = 16.0 / 9.0
	const imageWidth = 1280
	const imageHeight = int(imageWidth / aspectRatio)

	// Camera config
	const viewportHeight = 2.0
	const viewportWidth = aspectRatio * viewportHeight
	const focalLength = 1.0

	var origin = vector.Vector{0, 0, 0}
	var horizontal = vector.Vector{viewportWidth, 0, 0}
	var vertical = vector.Vector{0, viewportHeight, 0}
	var lowerLeftCorner = origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(vector.Vector{0, 0, focalLength})

	image := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Cast rays for each of the image pixels
	for y := imageHeight - 1; y >= 0; y-- {
		for x := 0; x < imageWidth; x++ {
			var u float64 = float64(x) / (imageWidth + 1)
			var v float64 = float64(y) / float64(imageHeight+1)
			r := ray.New(origin, lowerLeftCorner.Add(horizontal.Scale(u)).Add(vertical.Scale(v)).Sub(origin))
			image.Set(x, y, colorRay(r))
		}
	}

	// Save image
	f, err := os.Create("outimage.jpg")
	if err != nil {
		// Handle error
	}
	defer f.Close()

	// Specify the quality, between 0-100
	// Higher is better
	opt := jpeg.Options{
		Quality: 90,
	}
	err = jpeg.Encode(f, image, &opt)
	if err != nil {
		// Handle error
	}
}

// Red to blue gradient based on y coord
func colorRay(r ray.Ray) gocolor.RGBA {
	t := float32(0.5 * (r.Direction().Normalise().Y + 1.0))
	return color.New(1.0, 0, 0).Mul(1-t, 1-t, 1-t).Add(color.New(0, 0, 1).Mul(t, t, t)).RGBA()
}
