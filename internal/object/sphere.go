package object

import (
	"math"
	"raytracer/internal/color"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// Sphere is a sphere, duh
type Sphere struct {
	Center   vector.Vector
	Radius   float64
	Material Material
	Albedo   color.RGB
}

// NewSphere creates a new Sphere
func NewSphere(center vector.Vector, radius float64, material Material) *Sphere {
	return &Sphere{
		Center:   center,
		Radius:   radius,
		Material: material,
	}
}

// Intersect calculates the intersection of a ray r with this sphere, between tMin and tMax
func (s *Sphere) Intersect(r *ray.Ray, tMin, tMax float64, hit *Hit) bool {
	var oc = r.Origin().Sub(s.Center)
	l := r.Direction().Length()
	var a = l * l
	if a == 0 {
		return false
	}
	var halfB = oc.Dot(r.Direction())
	ocl := oc.Length()
	var c = ocl*ocl - s.Radius*s.Radius
	var discriminant = halfB*halfB - a*c

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
	hit.Material = s.Material

	return true
}
