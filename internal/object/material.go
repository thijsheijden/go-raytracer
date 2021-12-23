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

// NewMaterial creates a new Material, if the material does not exist, the basic lambertian material is returned
// Albedo is the color of the material
func NewMaterial(name string, albedo color.RGB) Material {
	switch name {
	case "lambertian":
		return lambertian{
			albedo: albedo,
		}
	case "metal":
		return metal{
			albedo: albedo,
		}
	default:
		return lambertian{
			albedo: albedo,
		}
	}
}

// Basic diffuse material
type lambertian struct {
	albedo color.RGB
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

// Metal material
type metal struct {
	albedo color.RGB
}

func (m metal) Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool {
	reflected := r.Direction().Normalise().Reflect(hit.Normal)
	*scattered = ray.New(hit.Point, reflected)
	*attenuation = m.albedo
	return scattered.Direction().Dot(hit.Normal) > 0
}
