package main

import (
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
var nPixelSamples = 500
var maxDepth = 50
var infinity = math.Inf(1)
var imageWidth = 1080

// Number of threads to split the work up
const nThreads = 16

type imagePart struct {
	index, height int
	image         *image.RGBA
}

func main() {
	const aspectRatio = 16.0 / 9.0
	loadedScene = scene.New(scene.NewCamera(vector.New(13, 2, 3), vector.New(0, 0, 0), vector.New(0, 1, 0), 20, aspectRatio, 1.0), aspectRatio, imageWidth)
	loadedScene.LotsOfSpheres()

	var rowsPerThread = loadedScene.ImageHeight / nThreads
	var rest = loadedScene.ImageHeight - (nThreads * rowsPerThread)

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
			img := image.NewRGBA(image.Rect(0, 0, loadedScene.ImageWidth, (endY - startY)))

			// Cast rays for each of the image pixels
			for y := startY; y <= endY; y++ {
				for x := 0; x < loadedScene.ImageWidth; x++ {
					// Define a new color for this pixel, which we will average later
					var pixelColor color.RGB = color.New(0, 0, 0)

					// Anti-aliasing
					for i := 0; i < nPixelSamples; i++ {
						var u float64 = (float64(x) + random.Float64()) / (loadedScene.FloatImageWidth + 1.0)
						var v float64 = (float64(y) + random.Float64()) / float64(loadedScene.FloatImageHeight+1)
						r := ray.New(loadedScene.Origin, loadedScene.LowerLeftCorner.Add(loadedScene.Horizontal.Scale(u)).Add(loadedScene.Vertical.Scale(v)).Sub(loadedScene.Origin))
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
	fullImage := image.NewRGBA(image.Rect(0, 0, loadedScene.ImageWidth, loadedScene.ImageHeight))
	for part := range imagePartChan {
		r := image.Rect(0, part.index*rowsPerThread, loadedScene.ImageWidth, part.index*rowsPerThread+part.height)
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
