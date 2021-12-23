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
	Point     vector.Vector // Hit point
	Normal    vector.Vector // Normal for the hit
	T         float64       // The distance at which the hit occurred
	FrontFace bool          // Whether the normal faces outwards
	Material  Material      // A pointer to the material that was hit
}

// SetFaceNormal sets the normal based on the dot product between the ray direction and the outward normal
func (h *Hit) SetFaceNormal(r *ray.Ray, outwardNormal *vector.Vector) {
	h.FrontFace = r.Direction().Dot(*outwardNormal) < 0
	if h.FrontFace {
		h.Normal = *outwardNormal
	} else {
		h.Normal = outwardNormal.Scale(-1)
	}
}
