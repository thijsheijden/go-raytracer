package object

import (
	"math"
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

type dielectric struct {
	refractionIndex float64
}

// Dielectric creates a new dieletric material with the given refraction index
func Dielectric(refractionIndex float64) Material {
	return dielectric{
		refractionIndex: refractionIndex,
	}
}

func (m dielectric) Scatter(r *ray.Ray, hit *Hit, attenuation *color.RGB, scattered *ray.Ray, rand *rand.Rand) bool {
	*attenuation = color.New(1, 1, 1)

	var refractionRatio float64
	if hit.FrontFace {
		refractionRatio = 1.0 / m.refractionIndex
	} else {
		refractionRatio = m.refractionIndex
	}

	unitDirection := r.Direction().Normalise()
	cosTheta := min(unitDirection.Scale(-1).Dot(hit.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0
	var direction vector.Vector

	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand.Float64() {
		direction = unitDirection.Reflect(hit.Normal)
	} else {
		direction = unitDirection.Refract(hit.Normal, refractionRatio)
	}

	*scattered = ray.New(hit.Point, direction)
	return true
}

func reflectance(cos, refractionRatio float64) float64 {
	// Use Schlick's approximation for reflectance.
	r0 := (1 - refractionRatio) / (1 + refractionRatio)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
