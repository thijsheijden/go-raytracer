package vector

import (
	"math"
	"math/rand"
)

// Vector is a 3 dimensional vector
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// New creates a new Vector
func New(X, Y, Z float64) Vector {
	return Vector{
		X: X,
		Y: Y,
		Z: Z,
	}
}

// Random generates a random vector with the provided min and max values
func Random(min, max float64) Vector {
	return New(min+rand.Float64()*(max-min), min+rand.Float64()*(max-min), min+rand.Float64()*(max-min))
}

// RandomInUnitSphere generates a random vector within a unit sphere
func RandomInUnitSphere() Vector {
	for {
		p := Random(-1, 1)
		if math.Pow(p.Length(), 2) >= 1 {
			// Outside unit sphere
			continue
		}
		return p
	}
}

// Add adds two vectors
func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// Sub subtracts vector b from vector a
func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// Scale multiplies vector a by a scalar
func (a Vector) Scale(s float64) Vector {
	return Vector{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

// Dot calculates the dot product of two vectors a and b
func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Length gets the length of a vector
func (a Vector) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

// Cross calculates the cross product between two vectors a and b
func (a Vector) Cross(b Vector) Vector {
	return Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

// Normalise normalises a vector
func (a Vector) Normalise() Vector {
	return a.Scale(1 / a.Length())
}
