package scene

import (
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

// InitMaterials initialises all materials
func (s *Scene) InitMaterials() {
	for _, s := range s.Spheres {
		s.Material = object.NewMaterial(s.MaterialName, s.Albedo)
	}
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
