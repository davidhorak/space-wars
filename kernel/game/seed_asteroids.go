package game

import (
	"math/rand"

	"github.com/davidhorak/space-wars/kernel/physics"
)

func SeedAsteroids(random *rand.Rand, width, height float64, maxAttempts int) []GameObject {
	asteroids := make([]GameObject, 0)
	count := random.Intn(MaxAsteroids-MinAsteroids) + MinAsteroids
	for i := 0; i <= count && maxAttempts > 0; i++ {
		maxAttempts--
		radius := random.Float64()*(MaxAsteroidSize-MinAsteroidSize) + MinAsteroidSize
		x := radius + (random.Float64() * (width - 2*radius))
		y := radius + (random.Float64() * (height - 2*radius))

		regenerate := false
		for _, asteroid := range asteroids {
			position := asteroid.Position()
			if position.Distance(physics.Vector2{X: x, Y: y}) < radius+asteroid.(*Asteroid).radius+MinAsteroidSeparation {
				regenerate = true
				break
			}
		}

		if regenerate {
			i--
			continue
		}

		asteroids = append(asteroids, NewAsteroid(NewUUID(), physics.Vector2{X: x, Y: y}, radius))
	}

	return asteroids
}
