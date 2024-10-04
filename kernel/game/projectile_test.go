package game

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectile(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	assert.GreaterOrEqual(t, projectile.ID(), int64(1))
	assert.Equal(t, DamageTypeUnknown, projectile.DamageType())
	assert.True(t, projectile.enabled)
	assert.Equal(t, physics.Vector2{X: 15, Y: 30}, projectile.Position())
	assert.Equal(t, 10.0, projectile.velocity.X)
	assert.Equal(t, 20.0, projectile.velocity.Y)
	assert.Equal(t, math.Pi, projectile.rotation)
	assert.Equal(t, 5.0, projectile.lifespanSec)
	assert.Equal(t, 20.0, projectile.damage)
	assert.Equal(t, owner, projectile.owner)
	assert.Equal(t, 1.0, projectile.collider.(*collider.SquareCollider).Size().Width)
	assert.Equal(t, 1.0, projectile.collider.(*collider.SquareCollider).Size().Height)
}

func TestProjectile_ID(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	assert.GreaterOrEqual(t, projectile.ID(), int64(1))
}

func TestProjectile_DamageType(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)
	projectile.damageType = DamageTypeRocket

	assert.Equal(t, DamageTypeRocket, projectile.DamageType())
}

func TestProjectile_Enabled(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)
	projectile.enabled = false

	assert.False(t, projectile.Enabled())
}

func TestProjectile_SetEnabled(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	projectile.SetEnabled(false)
	assert.False(t, projectile.Enabled())

	projectile.SetEnabled(true)
	assert.True(t, projectile.Enabled())
}

func TestProjectile_Damage(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)
	projectile.damage = 30.0

	assert.Equal(t, 30.0, projectile.Damage())
}

func TestProjectile_Position(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	assert.Equal(t, physics.Vector2{X: 15, Y: 30}, projectile.Position())
}

func TestProjectile_SetPosition(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	projectile.SetPosition(physics.Vector2{X: 20, Y: 40})
	assert.Equal(t, physics.Vector2{X: 20, Y: 40}, projectile.Position())
}

func TestProjectile_Update(t *testing.T) {
	gameManager := NewGameManager()
	lifespanSec := 5.0
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, lifespanSec, 20.0, owner)

	projectile.Update(1000, &gameManager)
	assert.Equal(t, lifespanSec-1.0, projectile.lifespanSec)
	assert.Equal(t, physics.Vector2{X: 25, Y: 50}, projectile.Position())

	projectile.Update(4000, &gameManager)
	assert.Equal(t, 0.0, projectile.lifespanSec)
	assert.False(t, projectile.enabled)
}

func TestProjectile_Collider(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	assert.Equal(t, projectile.collider, projectile.Collider())
}

func TestProjectile_Destroy(t *testing.T) {
	gameManager := NewGameManager()
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)
	projectile.explosionRadius = 1.0
	projectile.explosionDurationSec = 1.0

	// Without explosion
	gameManager.AddGameObject(projectile)
	projectile.Destroy(&gameManager, false)

	assert.Equal(t, 0.0, projectile.lifespanSec)
	assert.False(t, projectile.enabled)
	assert.Equal(t, 0, len(gameManager.GameObjects()))

	// With explosion
	gameManager.AddGameObject(projectile)
	projectile.Destroy(&gameManager, true)

	assert.Equal(t, 0.0, projectile.lifespanSec)
	assert.False(t, projectile.enabled)
	assert.Equal(t, 1, len(gameManager.GameObjects()))
	assert.IsType(t, &Explosion{}, gameManager.GameObjects()[0])
	explosion := gameManager.GameObjects()[0].(*Explosion)
	assert.Equal(t, physics.Vector2{X: 14, Y: 29}, explosion.Position())
	assert.Equal(t, 1.0, explosion.radius)
	assert.Equal(t, 1.0, explosion.lifespanSec)
}

func TestProjectile_OnCollision(t *testing.T) {
	gameManager := NewGameManager()
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	other := NewSpaceship(2, "other", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	gameManager.AddGameObject(projectile)
	gameManager.AddGameObject(owner)

	// Collide with owner
	projectile.OnCollision(owner, &gameManager, 0)

	assert.Equal(t, 2, len(gameManager.GameObjects()))

	// Collide with other
	gameManager.AddGameObject(other)
	projectile.OnCollision(other, &gameManager, 0)

	assert.Equal(t, "\"owner\" did 20.00 damage to \"other\" with unknown", gameManager.Logger().Logs()[0].message)
	assert.Equal(t, 3, len(gameManager.GameObjects()))
	assert.Equal(t, 80.0, other.health)
	assert.Equal(t, 10.0, owner.score)
}

func TestProjectile_Serialize(t *testing.T) {
	owner := NewSpaceship(1, "owner", physics.Vector2{X: 15, Y: 30}, 100)
	projectile := NewProjectile(physics.Vector2{X: 15, Y: 30}, physics.Vector2{X: 10, Y: 20}, math.Pi, 5.0, 20.0, owner)

	assert.Equal(t, map[string]interface{}{}, projectile.Serialize())
}
