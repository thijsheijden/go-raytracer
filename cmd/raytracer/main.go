package main

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"raytracer/internal/color"
	"raytracer/internal/object"
	"raytracer/internal/ray"
	"raytracer/internal/scene"
	"raytracer/internal/vector"

	_ "image/jpeg" // Needed for JPEG decoder
)

var loadedScene scene.Scene
var nPixelSamples = 100

func main() {
	// Load in a scene file
	f, err := os.Open("scenes/singleSphere.json")
	if err != nil {
		panic(err)
	}

	// Decode the scene
	err = json.NewDecoder(f).Decode(&loadedScene)
	if err != nil {
		panic(err)
	}
	f.Close()

	// Image config
	const aspectRatio = 16.0 / 9.0
	const imageWidth = 1280
	const imageHeight = int(imageWidth / aspectRatio)

	// Camera config
	const viewportHeight = 2.0
	const viewportWidth = aspectRatio * viewportHeight
	const focalLength = 1

	var origin = vector.Vector{0, 0, 0}
	var horizontal = vector.Vector{viewportWidth, 0, 0}
	var vertical = vector.Vector{0, viewportHeight, 0}
	var lowerLeftCorner = origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(vector.Vector{0, 0, focalLength})

	image := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Cast rays for each of the image pixels
	for y := imageHeight - 1; y >= 0; y-- {
		for x := 0; x < imageWidth; x++ {
			// Define a new color for this pixel, which we will average later
			var pixelColor color.RGB = color.New(0, 0, 0)

			// Anti-aliasing
			for i := 0; i < nPixelSamples; i++ {
				var u float64 = (float64(x) + rand.Float64()) / (imageWidth + 1)
				var v float64 = (float64(y) + rand.Float64()) / float64(imageHeight+1)
				r := ray.New(origin, lowerLeftCorner.Add(horizontal.Scale(u)).Add(vertical.Scale(v)).Sub(origin))
				pixelColor = pixelColor.Add(colorRay(r))
			}

			// Set pixel in the image
			image.Set(x, imageHeight-y, pixelColor.Average(nPixelSamples).RGBA()) // Flipping y coordinates
		}
	}

	// Save image
	f, err = os.Create("outimage.jpg")
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
func colorRay(r ray.Ray) color.RGB {
	var hit object.Hit

	// Go over all spheres
	closest := 15.0 // tMax
	hitSomething := false
	var normal vector.Vector
	for _, s := range loadedScene.Spheres {
		if s.Intersect(&r, 0, 15, &hit) {
			if hit.T < closest {
				closest = hit.T
				normal = hit.Normal
				hitSomething = true
			}
		}
	}

	if hitSomething {
		return color.New(float32(normal.X)+1, float32(normal.Y)+1, float32(normal.Z)+1).Mul(0.5, 0.5, 0.5)
	}

	t := float32(0.5 * (r.Direction().Normalise().Y + 1.0))
	return color.New(1, 1, 1).Mul(1-t, 1-t, 1-t).Add(color.New(0.5, 0.7, 1).Mul(t, t, t))
}
