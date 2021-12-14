package object

import (
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// A Hittable object is an object that can be hit
type Hittable interface {
	Intersect(r *ray.Ray, tMin, tMax float64, hit *Hit) bool
}

// A Hit contains information returned upon a hit
type Hit struct {
	Point  vector.Vector // Hit point
	Normal vector.Vector // Normal for the hit
	T      float64       // The distance at which the hit occurred
}
