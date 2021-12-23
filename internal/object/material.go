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

// Lambertian returns a lambertian material
func Lambertian(albedo color.RGB) Material {
	return lambertian{
		albedo: albedo,
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

// Metal material
type metal struct {
	albedo color.RGB
}

// Metal returns a metal material
func Metal(albedo color.RGB) Material {
	return metal{
		albedo: albedo,
	}
}

func (m metal) Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool {
	reflected := r.Direction().Normalise().Reflect(hit.Normal)
	*scattered = ray.New(hit.Point, reflected)
	*attenuation = m.albedo
	return scattered.Direction().Dot(hit.Normal) > 0
}

type fuzzyMetal struct {
	albedo    color.RGB
	fuzziness float64
}

// FuzzyMetal returns a fuzzy metal material
func FuzzyMetal(albedo color.RGB, fuzziness float64) Material {
	return fuzzyMetal{
		albedo:    albedo,
		fuzziness: fuzziness,
	}
}

func (m fuzzyMetal) Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool {
	reflected := r.Direction().Normalise().Reflect(hit.Normal)
	*scattered = ray.New(hit.Point, reflected.Add(vector.RandomInUnitSphere(rand).Scale(m.fuzziness)))
	*attenuation = m.albedo
	return scattered.Direction().Dot(hit.Normal) > 0
}
