package main

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
	"raytracer/internal/color"
	"raytracer/internal/object"
	"raytracer/internal/ray"
	"raytracer/internal/scene"
	"raytracer/internal/vector"
	"runtime/pprof"
	"sync"
	"time"

	_ "image/jpeg" // Needed for JPEG decoder
)

var loadedScene scene.Scene
var nPixelSamples = 50
var maxDepth = 10
var infinity = math.Inf(1)

// Number of threads to split the work up
const nThreads = 8

type imagePart struct {
	index, height int
	image         *image.RGBA
}

func main() {
	// Load in a scene file
	f, err := os.Open("scenes/threeSpheres.json")
	if err != nil {
		panic(err)
	}

	// Decode the scene
	err = json.NewDecoder(f).Decode(&loadedScene)
	if err != nil {
		panic(err)
	}
	f.Close()
	loadedScene.InitMaterials()

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

	const rowsPerThread = imageHeight / nThreads
	const rest = imageHeight - (nThreads * rowsPerThread)

	cpuProfile, err := os.Create("profile.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	var wg sync.WaitGroup
	imagePartChan := make(chan imagePart, nThreads)
	for thread := 0; thread < nThreads; thread++ {
		// If this is the last thread, it gets the remaining rows from incomplete division
		startY := thread * rowsPerThread
		endY := (thread + 1) * rowsPerThread
		if thread == nThreads-1 {
			endY += rest
		}

		// Generate random generator
		random := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Join waitgroup
		wg.Add(1)

		go func(thread, startY, endY int, random *rand.Rand, imagePartChan chan imagePart) {
			defer wg.Done()
			// Create slice for image part
			img := image.NewRGBA(image.Rect(0, 0, imageWidth, (endY - startY)))

			// Cast rays for each of the image pixels
			for y := startY; y <= endY; y++ {
				for x := 0; x < imageWidth; x++ {
					// Define a new color for this pixel, which we will average later
					var pixelColor color.RGB = color.New(0, 0, 0)

					// Anti-aliasing
					for i := 0; i < nPixelSamples; i++ {
						var u float64 = (float64(x) + random.Float64()) / (imageWidth + 1)
						var v float64 = (float64(y) + random.Float64()) / float64(imageHeight+1)
						r := ray.New(origin, lowerLeftCorner.Add(horizontal.Scale(u)).Add(vertical.Scale(v)).Sub(origin))
						pixelColor = pixelColor.Add(colorRay(r, maxDepth, random))
					}

					// Set pixel in slice
					img.Set(x, endY-y, pixelColor.Average(nPixelSamples).RGBA())
				}
			}

			// Thread is done
			imagePartChan <- imagePart{index: nThreads - thread - 1, height: endY - startY, image: img}

		}(thread, startY, endY, random, imagePartChan)
	}

	// Wait for all threads to finish
	wg.Wait()
	close(imagePartChan)

	// Stitch all image parts together
	fullImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for part := range imagePartChan {
		r := image.Rect(0, part.index*rowsPerThread, imageWidth, part.index*rowsPerThread+part.height)
		draw.Draw(fullImage, r, part.image, image.ZP, draw.Src)
	}

	// Save image
	output, err := os.Create("outimage.jpg")
	if err != nil {
		// Handle error
	}
	defer output.Close()

	// Specify the quality, between 0-100
	// Higher is better
	opt := jpeg.Options{
		Quality: 90,
	}
	err = jpeg.Encode(output, fullImage, &opt)
	if err != nil {
		// Handle error
	}
}

// Red to blue gradient based on y coord
func colorRay(r ray.Ray, depth int, random *rand.Rand) color.RGB {
	// Reached max recursion depth
	if depth <= 0 {
		return color.New(0, 0, 0)
	}

	var hit object.Hit

	if loadedScene.Hit(&r, 0.001, infinity, &hit) {
		var scattered ray.Ray
		var attenuation color.RGB

		if hit.Material.Scatter(&r, &hit, &attenuation, &scattered, random) {
			return colorRay(scattered, depth-1, random).Mul(attenuation.R, attenuation.G, attenuation.B)
		}
		return color.New(0, 0, 0)
	}

	t := float32(0.5 * (r.Direction().Normalise().Y + 1.0))
	return color.New(1, 1, 1).Mul(1-t, 1-t, 1-t).Add(color.New(0.5, 0.7, 1).Mul(t, t, t))
}
