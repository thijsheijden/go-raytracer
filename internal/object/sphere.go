package object

import (
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// Sphere is a sphere, duh
type Sphere struct {
	center vector.Vector
	radius float64
}

// New creates a new Sphere
func New(center vector.Vector, radius float64) Sphere {
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
func (s *Sphere) Intersect(r ray.Ray) bool {
	var a = r.Direction().Dot(r.Direction())
	var oc = r.Origin().Sub(s.center)
	var b = 2.0 * oc.Dot(r.Direction())
	var c = oc.Dot(oc) - s.radius*s.radius
	var discriminant = b*b - 4*a*c
	return discriminant > 0
}
