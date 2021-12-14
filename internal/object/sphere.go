package object

import (
	"math"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// Sphere is a sphere, duh
type Sphere struct {
	Center vector.Vector `json:"center"`
	Radius float64       `json:"radius"`
}

// NewSphere creates a new Sphere
func NewSphere(center vector.Vector, radius float64) Sphere {
	return Sphere{
		Center: center,
		Radius: radius,
	}
}

// Intersect calculates the intersection of a ray r with this sphere, between tMin and tMax
func (s *Sphere) Intersect(r *ray.Ray, tMin, tMax float64, hit *Hit) bool {
	var oc = r.Origin().Sub(s.Center)
	var a = math.Pow(r.Direction().Length(), 2)
	var halfB = oc.Dot(r.Direction())
	var c = math.Pow(oc.Length(), 2) - math.Pow(s.Radius, 2)
	var discriminant = math.Pow(halfB, 2) - a*c

	// No hit
	if discriminant < 0 {
		return false
	}

	var sqrtD = math.Sqrt(discriminant)

	// Find the nearest t
	var root = (-halfB - sqrtD) / a
	if root < tMin || tMax < root {
		root = (-halfB + sqrtD) / a
		if root < tMin || tMax < root {
			return false
		}
	}

	hit.T = root
	hit.Point = r.At(root)
	outwardNormal := hit.Point.Sub(s.Center).Scale(1 / s.Radius)
	hit.SetFaceNormal(r, &outwardNormal)

	return true
}
