package game

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeedAsteroids(t *testing.T) {
	random := rand.New(rand.NewSource(1))
	width, height := 1000.0, 1000.0

	asteroids := SeedAsteroids(random, width, height, 1000)

	assert.GreaterOrEqual(t, len(asteroids), MinAsteroids)
	assert.LessOrEqual(t, len(asteroids), MaxAsteroids)

	for i := 0; i < len(asteroids); i++ {
		pos := asteroids[i].Position()
		radius := asteroids[i].(*Asteroid).radius

		assert.GreaterOrEqual(t, pos.X, radius)
		assert.LessOrEqual(t, pos.X, width-radius)
		assert.GreaterOrEqual(t, pos.Y, radius)
		assert.LessOrEqual(t, pos.Y, height-radius)

		assert.GreaterOrEqual(t, radius, float64(MinAsteroidSize))
		assert.LessOrEqual(t, radius, float64(MaxAsteroidSize))

		for j := i + 1; j < len(asteroids); j++ {
			radius := asteroids[i].(*Asteroid).radius
			assert.GreaterOrEqual(t, radius, float64(MinAsteroidSize))
			assert.LessOrEqual(t, radius, float64(MaxAsteroidSize))

			for j := i + 1; j < len(asteroids); j++ {
				pos1 := asteroids[i].Position()
				pos2 := asteroids[j].Position()
				radius1 := asteroids[i].(*Asteroid).radius
				radius2 := asteroids[j].(*Asteroid).radius

				distance := pos1.Distance(pos2)
				minDistance := radius1 + radius2 + MinAsteroidSeparation

				assert.GreaterOrEqual(t, distance, minDistance)
			}
		}
	}
}

func TestSeedAsteroids_Regenerate(t *testing.T) {
	random := rand.New(rand.NewSource(1))
	width, height := 100.0, 100.0

	asteroids := SeedAsteroids(random, width, height, 100)

	assert.GreaterOrEqual(t, len(asteroids), MinAsteroids)
	assert.LessOrEqual(t, len(asteroids), MaxAsteroids)
}
