package ray

import "raytracer/internal/vector"

// Ray is a ray
type Ray struct {
	origin    vector.Vector
	direction vector.Vector
}

func New(origin, direction vector.Vector) Ray {
	return Ray{
		origin:    origin,
		direction: direction,
	}
}

// Origin returns the ray origin
func (r *Ray) Origin() vector.Vector {
	return r.origin
}

// Direction returns the ray direction
func (r *Ray) Direction() vector.Vector {
	return r.direction
}

// At returns the location on ray r after t steps
func (r *Ray) At(t float64) vector.Vector {
	return r.origin.Add(r.direction.Scale(t))
}
