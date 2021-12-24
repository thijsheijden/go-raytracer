package scene

import (
	"raytracer/internal/color"
	"raytracer/internal/object"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// A Scene contains all objects
// Can be loaded in from a JSON file
type Scene struct {
	Camera                            Camera
	Origin                            vector.Vector
	Horizontal                        vector.Vector
	Vertical                          vector.Vector
	LowerLeftCorner                   vector.Vector
	Spheres                           []*object.Sphere
	ImageHeight, ImageWidth           int
	FloatImageHeight, FloatImageWidth float64
}

// Camera is the camera
type Camera struct {
	Position       vector.Vector
	LookDirection  vector.Vector
	VerticalFOV    float64
	FocalLength    float64
	ViewportWidth  float64
	ViewportHeight float64
}

// New creates a new scene
func New(camera Camera, aspectRatio float64, imageWidth int) Scene {
	origin := vector.New(0, 0, 0)
	horizontal := vector.New(camera.ViewportWidth, 0, 0)
	vertical := vector.New(0, camera.ViewportHeight, 0)
	imageHeight := int(float64(imageWidth) / aspectRatio)
	return Scene{
		Camera:           camera,
		Origin:           origin,
		Horizontal:       horizontal,
		Vertical:         vertical,
		LowerLeftCorner:  origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(vector.New(0, 0, camera.FocalLength)),
		ImageWidth:       imageWidth,
		ImageHeight:      imageHeight,
		FloatImageWidth:  float64(imageWidth),
		FloatImageHeight: float64(imageHeight),
		Spheres:          make([]*object.Sphere, 0),
	}
}

// NewCamera creates a new camera
func NewCamera(position, lookDirection vector.Vector, verticalFOV, viewportHeight, aspectRatio, focalLength float64) Camera {
	return Camera{
		Position:       position,
		LookDirection:  lookDirection,
		VerticalFOV:    verticalFOV,
		ViewportWidth:  aspectRatio * viewportHeight,
		ViewportHeight: viewportHeight,
		FocalLength:    focalLength,
	}
}

// ThreeBalls returns a scene with three balls
func (s *Scene) ThreeBalls() {

	// Center metal sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, 0, -1), 0.5, object.Metal(color.New(0.8, 0.8, 0.8))))

	// Right purple sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(1, 0, -1), 0.5, object.Dielectric(1.5)))

	// Left sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(-1, 0, -1), 0.5, object.FuzzyMetal(color.New(0.8, 0.8, 0.8), 0.2)))

	// Ground plane sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, -100.5, -1), 100, object.Lambertian(color.New(0, 1, 0))))
}

// GlassBalls places balls of glass
func (s *Scene) GlassBalls() {
	// Center glass sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, 0, -1), 0.5, object.Dielectric(1.6)))

	// Teal sphere left
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(-1, 0, -1), 0.5, object.Lambertian(color.New(0, 0.6, 0.6))))

	// Light green sphere behind center
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, 0, -6), 0.5, object.Lambertian(color.New(0.2, 0.8, 0.2))))
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(1.5, 0, -6), 0.5, object.Lambertian(color.New(0.5, 0, 0.5))))

	// Light red fuzzy metal sphere to the right
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(1, 0, -1), 0.5, object.FuzzyMetal(color.New(0.5, 0, 0), 0.2)))

	// Ground plane sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, -100.5, -1), 100, object.Lambertian(color.New(0.6, 0.6, 0.6))))
}

// Hit checks for hits in the scene
func (s *Scene) Hit(r *ray.Ray, tMin, tMax float64, hit *object.Hit) bool {
	var tempHit object.Hit
	hitAnything := false
	closestSoFar := tMax

	for _, sphere := range s.Spheres {
		if sphere.Intersect(r, tMin, closestSoFar, &tempHit) {
			hitAnything = true
			closestSoFar = tempHit.T
			*hit = tempHit
		}
	}

	return hitAnything
}
