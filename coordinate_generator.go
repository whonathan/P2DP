// coordinate_config.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const (
	defaultCenterLat  = -7.093408
	defaultCenterLong = 109.280774
	defaultRadius     = 7.00

	// Custom coordinates for j03
	j03CenterLat  = -7.139337
	j03CenterLong = 109.252766
	j03Radius     = 0.01

	// Custom coordinates for j07
	j07CenterLat  = -7.139305
	j07CenterLong = 109.241971
	j07Radius     = 1.00

	// Custom coordinates for j08
	j08CenterLat  = -7.09902
	j08CenterLong = 109.32261
	j08Radius     = 0.01

	pi180       = math.Pi / 180
	_180pi      = 180 / math.Pi
	earthRadius = 6371.0
)

// CoordinateConfig holds the configuration for generating coordinates
type CoordinateConfig struct {
	CenterLat  float64
	CenterLong float64
	Radius     float64
}

// CoordinateGenerator handles generating random coordinates within specified bounds
type CoordinateGenerator struct {
	rng    *rand.Rand
	mu     sync.Mutex
	config CoordinateConfig
}

// NewCoordinateGenerator creates a new coordinate generator with the specified username
func NewCoordinateGenerator(seed int64, username string) *CoordinateGenerator {
	config := getCoordinateConfig(username)
	return &CoordinateGenerator{
		rng:    rand.New(rand.NewSource(seed)),
		config: config,
	}
}

// getCoordinateConfig returns the appropriate coordinate configuration based on username
func getCoordinateConfig(username string) CoordinateConfig {
	switch username {
	case "52260.j03":
		return CoordinateConfig{
			CenterLat:  j03CenterLat,
			CenterLong: j03CenterLong,
			Radius:     j03Radius,
		}
	case "52260.j07":
		return CoordinateConfig{
			CenterLat:  j07CenterLat,
			CenterLong: j07CenterLong,
			Radius:     j07Radius,
		}
	case "52260.j08":
		return CoordinateConfig{
			CenterLat:  j08CenterLat,
			CenterLong: j08CenterLong,
			Radius:     j08Radius,
		}
	default:
		return CoordinateConfig{
			CenterLat:  defaultCenterLat,
			CenterLong: defaultCenterLong,
			Radius:     defaultRadius,
		}
	}
}

// GenerateCoordinates generates a random pair of coordinates within the configured bounds
func (g *CoordinateGenerator) GenerateCoordinates() (string, string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	angle := g.rng.Float64() * 2 * math.Pi
	distance := math.Sqrt(g.rng.Float64()) * g.config.Radius

	centerLatRad := g.config.CenterLat * pi180
	centerLonRad := g.config.CenterLong * pi180

	distRatio := distance / earthRadius
	sinDist := math.Sin(distRatio)
	cosDist := math.Cos(distRatio)
	sinCenterLat := math.Sin(centerLatRad)
	cosCenterLat := math.Cos(centerLatRad)
	sinAngle := math.Sin(angle)
	cosAngle := math.Cos(angle)

	newLatRad := math.Asin(sinCenterLat*cosDist + cosCenterLat*sinDist*cosAngle)
	newLonRad := centerLonRad + math.Atan2(sinAngle*sinDist*cosCenterLat,
		cosDist-sinCenterLat*math.Sin(newLatRad))

	return fmt.Sprintf("%.6f", newLatRad*_180pi), fmt.Sprintf("%.6f", newLonRad*_180pi)
}

// GetConfig returns the current coordinate configuration
func (g *CoordinateGenerator) GetConfig() CoordinateConfig {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.config
}
