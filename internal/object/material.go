package object

import (
	"math/rand"
	"raytracer/internal/color"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// Material describes a material
type Material interface {
	Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool
}

// Basic diffuse material
type lambertian struct {
	albedo color.RGB
}

// NewMaterial creates a new Material, if the material does not exist, the basic lambertian material is returned
// Albedo is the color of the material
func NewMaterial(name string, albedo color.RGB) Material {
	switch name {
	case "lambertian":
		return lambertian{
			albedo: albedo,
		}
	default:
		return lambertian{
			albedo: albedo,
		}
	}
}

func (m lambertian) Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool {
	scatterDirection := hit.Normal.Add(vector.RandomInUnitSphere(rand))

	// Catch near zero scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = hit.Normal
	}

	scatteredRay := ray.New(hit.Point, scatterDirection)
	*scattered = scatteredRay
	*attenuation = m.albedo
	return true
}
