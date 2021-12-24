package scene

import (
	"math"
	"math/rand"
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
	Position       vector.Vector // The position of the camera
	LookAt         vector.Vector // The point the camera is looking at
	VUp            vector.Vector // The vector pointing up from the camera, used for rotation
	VerticalFOV    float64
	FocalLength    float64
	ViewportWidth  float64
	ViewportHeight float64
}

// New creates a new scene
func New(camera Camera, aspectRatio float64, imageWidth int) Scene {
	// Camera direction and placement
	w := camera.Position.Sub(camera.LookAt).Normalise()
	u := camera.VUp.Cross(w).Normalise()
	v := w.Cross(u)

	// Scene constants
	origin := camera.Position
	horizontal := u.Scale(camera.ViewportWidth)
	vertical := v.Scale(camera.ViewportHeight)
	imageHeight := int(float64(imageWidth) / aspectRatio)

	return Scene{
		Camera:           camera,
		Origin:           origin,
		Horizontal:       horizontal,
		Vertical:         vertical,
		LowerLeftCorner:  origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(w),
		ImageWidth:       imageWidth,
		ImageHeight:      imageHeight,
		FloatImageWidth:  float64(imageWidth),
		FloatImageHeight: float64(imageHeight),
		Spheres:          make([]*object.Sphere, 0),
	}
}

// NewCamera creates a new camera
func NewCamera(position, lookAt, vup vector.Vector, verticalFOV, aspectRatio, focalLength float64) Camera {
	// FOV calculations
	theta := verticalFOV * (math.Pi / 180)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h

	return Camera{
		Position:       position,
		LookAt:         lookAt,
		VUp:            vup,
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
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(1, 0, -1), 0.5, object.Metal(color.New(0.8, 0.8, 0.8))))

	// Ground plane sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, -100.5, -1), 100, object.Lambertian(color.New(0.6, 0.6, 0.6))))
}

func (s *Scene) LotsOfSpheres() {
	groundMaterial := object.Lambertian(color.New(0.5, 0.5, 0.5))
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, -1000, 0), 1000, groundMaterial))

	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			chooseMat := rand.Float64()
			center := vector.New(a+0.9*rand.Float64(), 0.2, b+0.9*rand.Float64())
			if center.Sub(vector.New(4, 0.2, 0)).Length() > 0.9 {

				if chooseMat < 0.8 {
					s.Spheres = append(s.Spheres, object.NewSphere(center, 0.2, object.Lambertian(color.Random())))
					continue
				}

				if chooseMat < 0.95 {
					s.Spheres = append(s.Spheres, object.NewSphere(center, 0.2, object.FuzzyMetal(color.RandomInRange(0.5, 1), randomInRange(0, 0.5))))
					continue
				}

				s.Spheres = append(s.Spheres, object.NewSphere(center, 0.2, object.Dielectric(1.5)))
			}
		}
	}

	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, 1, 0), 1.0, object.Dielectric(1.5)))
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(-4, 1, 0), 1.0, object.Lambertian(color.New(0.4, 0.2, 0.1))))
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(4, 1, 0), 1.0, object.Metal(color.New(0.7, 0.6, 0.5))))
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

func randomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
