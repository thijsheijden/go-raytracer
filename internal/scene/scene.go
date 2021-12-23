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
	Camera  camera           `json:"camera"`
	Spheres []*object.Sphere `json:"spheres"`
}

type camera struct {
	Position vector.Vector `json:"position"`
}

// ThreeBalls returns a scene with three balls
func ThreeBalls() Scene {
	s := Scene{
		Camera: camera{
			Position: vector.New(0, 0, 0),
		},
		Spheres: make([]*object.Sphere, 0),
	}

	// Center metal sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, 0, -1), 0.5, object.Metal(color.New(0.8, 0.8, 0.8))))

	// Right purple sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(1, 0, -1), 0.5, object.Lambertian(color.New(0.6, 0, 1))))

	// Left sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(-1, 0, -1), 0.5, object.FuzzyMetal(color.New(0.8, 0.8, 0.8), 0.2)))

	// Ground plane sphere
	s.Spheres = append(s.Spheres, object.NewSphere(vector.New(0, -100.5, -1), 100, object.Lambertian(color.New(0, 1, 0))))

	return s
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
