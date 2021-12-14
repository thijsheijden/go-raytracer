package object

import (
	"math"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// Sphere is a sphere, duh
type Sphere struct {
	center vector.Vector
	radius float64
}

// NewSphere creates a new Sphere
func NewSphere(center vector.Vector, radius float64) Sphere {
	return Sphere{
		center: center,
		radius: radius,
	}
}

// Center returns the sphere center
func (s *Sphere) Center() vector.Vector {
	return s.center
}

// Radius returns the sphere radius
func (s *Sphere) Radius() float64 {
	return s.radius
}

// Intersect calculates the intersection of a ray r with this sphere, between tMin and tMax
func (s *Sphere) Intersect(r *ray.Ray, tMin, tMax float64, hit *Hit) bool {
	var oc = r.Origin().Sub(s.center)
	var a = math.Pow(r.Direction().Length(), 2)
	var halfB = oc.Dot(r.Direction())
	var c = math.Pow(oc.Length(), 2) - math.Pow(s.radius, 2)
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
	hit.Normal = hit.Point.Sub(s.center).Scale(1 / s.radius)

	return true
}
