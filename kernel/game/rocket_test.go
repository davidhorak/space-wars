package game

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"

	"github.com/stretchr/testify/assert"
)

func TestNewRocketProjectile(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewRocketProjectile(1, physics.Vector2{X: 15, Y: 30}, math.Pi, owner)

	assert.Equal(t, int64(1), projectile.ID())
	assert.Equal(t, DamageTypeRocket, projectile.damageType)
	assert.Equal(t, true, projectile.Enabled())
	assert.Equal(t, physics.Vector2{X: 15, Y: 30}, projectile.Position())
	assert.Equal(t, math.Pi, projectile.rotation)
	assert.Equal(t, physics.Vector2{X: -274, Y: math.Sin(math.Pi) * 274}, projectile.velocity)
	assert.Equal(t, 10.0, projectile.lifespanSec)
	assert.Equal(t, 60.0, projectile.damage)
	assert.Equal(t, owner, projectile.owner)
	assert.Equal(t, 30.0, projectile.explosionRadius)
	assert.Equal(t, 1.0, projectile.explosionDurationSec)
	assert.Equal(t, 15.0, projectile.collider.Position().X)
	assert.Equal(t, 30.0, projectile.collider.Position().Y)
	assert.Equal(t, 20.0, projectile.collider.(*collider.CircleCollider).Radius())
}

func TestRocket_Serialize(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewRocketProjectile(1, physics.Vector2{X: 15, Y: 30}, math.Pi, owner)

	assert.Equal(t, map[string]interface{}{
		"type":    "rocket",
		"id":      projectile.ID(),
		"enabled": true,
		"position": map[string]interface{}{
			"x": 15.0,
			"y": 30.0,
		},
		"rotation": math.Pi,
		"velocity": map[string]interface{}{
			"x": -274.0,
			"y": math.Sin(math.Pi) * 274.0,
		},
		"lifespanSec": 10.0,
		"damage":      60.0,
		"owner":       projectile.owner.ID(),
		"collider":    projectile.collider.Serialize(),
	}, projectile.Serialize())
}
