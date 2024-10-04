package game

import (
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/stretchr/testify/assert"
)

func TestNewExplosion(t *testing.T) {
	position := physics.Vector2{X: 10, Y: 20}
	explosion := NewExplosion(1, position, 5, 2)

	assert.Equal(t, int64(1), explosion.ID())
	assert.True(t, explosion.Enabled())
	assert.Equal(t, position, explosion.Position())
	assert.Equal(t, 5.0, explosion.radius)
	assert.Equal(t, 2.0, explosion.lifespanSec)
}

func TestExplosion_Enabled(t *testing.T) {
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)
	assert.True(t, explosion.Enabled())
}

func TestExplosion_SetEnabled(t *testing.T) {
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)
	explosion.SetEnabled(false)
	assert.False(t, explosion.Enabled())
}

func TestExplosion_SetPosition(t *testing.T) {
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)
	newPosition := physics.Vector2{X: 15, Y: 25}
	explosion.SetPosition(newPosition)
	assert.Equal(t, newPosition, explosion.Position())
}

func TestExplosionUpdate(t *testing.T) {
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)
	gameManager := NewGameManager()

	// Test partial update
	explosion.Update(500, &gameManager)
	assert.Equal(t, 1.5, explosion.lifespanSec)
	assert.True(t, explosion.Enabled())

	// Test update that should disable the explosion
	explosion.Update(1600, &gameManager)
	assert.Equal(t, 0.0, explosion.lifespanSec)
	assert.False(t, explosion.Enabled())
}

func TestExplosion_Collider(t *testing.T) {
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)

	assert.Nil(t, explosion.Collider())
}

func TestExplosionSerialize(t *testing.T) {
	position := physics.Vector2{X: 10, Y: 20}
	explosion := NewExplosion(1, position, 5, 2)

	assert.Equal(t, map[string]interface{}{
		"type":    "explosion",
		"id":      int64(1),
		"enabled": true,
		"position": map[string]interface{}{
			"x": 10.0,
			"y": 20.0,
		},
		"radius":      5.0,
		"durationSec": 2.0,
		"lifespanSec": 2.0,
	}, explosion.Serialize())
}
