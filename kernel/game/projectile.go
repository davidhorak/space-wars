package game

import (
	"time"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type Projectile struct {
	id                   int64
	damageType           DamageType
	enabled              bool
	position             physics.Vector2
	velocity             physics.Vector2
	rotation             float64
	lifespanSec          float64
	damage               float64
	owner                *Spaceship
	collider             collider.Collider
	explosionRadius      float64
	explosionDurationSec float64
}

func NewProjectile(position physics.Vector2, velocity physics.Vector2, rotation float64, lifespanSec float64, damage float64, owner *Spaceship) *Projectile {
	return &Projectile{
		id:          NewUUID(),
		enabled:     true,
		damageType:  DamageTypeUnknown,
		position:    position,
		velocity:    velocity,
		rotation:    rotation,
		lifespanSec: lifespanSec,
		damage:      damage,
		owner:       owner,
		collider:    collider.NewSquareCollider(position, 1, physics.Size{Width: 1, Height: 1}),
	}
}

func (projectile *Projectile) ID() int64 {
	return projectile.id
}

func (projectile *Projectile) DamageType() DamageType {
	return projectile.damageType
}

func (projectile *Projectile) Enabled() bool {
	return projectile.enabled
}

func (projectile *Projectile) SetEnabled(enabled bool) {
	projectile.enabled = enabled
}

func (projectile *Projectile) Damage() float64 {
	return projectile.damage
}

func (projectile *Projectile) Position() physics.Vector2 {
	return projectile.position
}

func (projectile *Projectile) SetPosition(position physics.Vector2) {
	projectile.position = position
}

func (projectile *Projectile) Update(deltaTimeMs float64, gameManager *GameManager) {
	deltaTimeSec := deltaTimeMs / 1000
	projectile.lifespanSec -= deltaTimeSec
	if projectile.lifespanSec <= 0 {
		projectile.Destroy(gameManager, false)
		return
	}

	projectile.position = projectile.position.Add(projectile.velocity.Multiply(deltaTimeSec))
	projectile.collider.SetPosition(projectile.position)
}

func (projectile *Projectile) Collider() collider.Collider {
	return projectile.collider
}

func (projectile *Projectile) Destroy(gameManager *GameManager, createExplosion bool) {
	projectile.lifespanSec = 0
	projectile.enabled = false
	gameManager.RemoveGameObject(projectile)

	if createExplosion {
		gameManager.AddGameObject(NewExplosion(
			NewUUID(),
			physics.Vector2{
				X: projectile.position.X - float64(projectile.explosionRadius),
				Y: projectile.position.Y - float64(projectile.explosionRadius),
			},
			projectile.explosionRadius,
			projectile.explosionDurationSec,
		))
	}
}

func (projectile *Projectile) OnCollision(other GameObject, gameManager *GameManager, order int) {
	// Do not collide with the owner
	if spaceship, ok := other.(*Spaceship); ok {
		if spaceship.ID() == projectile.owner.ID() {
			return
		}

		gameManager.Logger().Damage(time.Now(), projectile.Damage(), projectile.owner.name, spaceship.name, projectile.damageType)
		spaceship.TakeDamage(projectile.damage, gameManager, projectile.owner)
		projectile.owner.AddScore(projectile.damage * ScorePerDamageCoefficient)
	}

	projectile.Destroy(gameManager, true)
}

func (projectile *Projectile) Serialize() map[string]interface{} {
	return map[string]interface{}{}
}
