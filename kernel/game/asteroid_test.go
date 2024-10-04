package game

import (
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/stretchr/testify/assert"
)

func TestNewAsteroid(t *testing.T) {
	id := int64(1)
	position := physics.Vector2{X: 10, Y: 20}
	radius := 5.0

	asteroid := NewAsteroid(id, position, radius)

	assert.Equal(t, id, asteroid.ID())
	assert.True(t, asteroid.Enabled())
	assert.Equal(t, position, asteroid.Position())
	assert.Equal(t, radius, asteroid.radius)
}

func TestAsteroid_Enabled(t *testing.T) {
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	assert.True(t, asteroid.Enabled())
}

func TestAsteroid_SetEnabled(t *testing.T) {
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	asteroid.SetEnabled(false)
	assert.False(t, asteroid.Enabled())
}

func TestAsteroid_SetPosition(t *testing.T) {
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)

	newPosition := physics.Vector2{X: 15, Y: 25}
	asteroid.SetPosition(newPosition)

	assert.Equal(t, newPosition, asteroid.Position())
}

func TestAsteroid_Serialize(t *testing.T) {
	id := int64(1)
	position := physics.Vector2{X: 10, Y: 20}
	radius := 5.0

	asteroid := NewAsteroid(id, position, radius)

	assert.Equal(t, map[string]interface{}{
		"type":    "asteroid",
		"id":      id,
		"enabled": true,
		"position": map[string]interface{}{
			"x": 10.0,
			"y": 20.0,
		},
		"radius":   radius,
		"collider": asteroid.collider.Serialize(),
	}, asteroid.Serialize())
}
