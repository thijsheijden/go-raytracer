package scene

import (
	"raytracer/internal/object"
	"raytracer/internal/vector"
)

// A Scene contains all objects
// Can be loaded in from a JSON file
type Scene struct {
	Camera  camera          `json:"camera"`
	Spheres []object.Sphere `json:"spheres"`
}

type camera struct {
	Position vector.Vector `json:"position"`
}
