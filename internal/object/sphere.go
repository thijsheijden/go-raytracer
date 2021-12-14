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

// Intersect checks whether a ray r intersects the sphere
// Returns the t at which the intersection occurs, if there is no intersection -1.0 is returned
func (s *Sphere) Intersect(r ray.Ray) float64 {
	var oc = r.Origin().Sub(s.center)
	var a = math.Pow(r.Direction().Length(), 2)
	var halfB = oc.Dot(r.Direction())
	var c = math.Pow(oc.Length(), 2) - math.Pow(s.radius, 2)
	var discriminant = math.Pow(halfB, 2) - a*c

	if discriminant > 0 {
		// Calculate t
		return (-halfB - math.Sqrt(discriminant)) / a
	}

	// No intersection
	return -1.0
}
