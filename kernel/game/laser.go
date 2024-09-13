package game

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type LaserProjectile struct {
	Projectile
}

func NewLaserProjectile(id int64, position physics.Vector2, rotation float64, owner *Spaceship) *LaserProjectile {
	direction := physics.Vector2{X: math.Cos(rotation), Y: math.Sin(rotation)}

	return &LaserProjectile{
		Projectile: Projectile{
			id:                   id,
			damageType:           DamageTypeLaser,
			enabled:              true,
			position:             position,
			rotation:             rotation,
			velocity:             direction.Multiply(LaserVelocitySec),
			lifespanSec:          LaserLifespanSec,
			damage:               LaserDamage,
			owner:                owner,
			explosionRadius:      float64(LaserExplosionRadius),
			explosionDurationSec: float64(LaserExplosionDurationSec),
			collider: collider.NewSquareCollider(
				position,
				rotation,
				physics.Size{Width: LaserWidth, Height: LaserLength},
			),
		},
	}
}

func (laser *LaserProjectile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "laser",
		"id":      laser.id,
		"enabled": laser.enabled,
		"position": map[string]interface{}{
			"x": laser.position.X,
			"y": laser.position.Y,
		},
		"rotation": laser.rotation,
		"velocity": map[string]interface{}{
			"x": laser.velocity.X,
			"y": laser.velocity.Y,
		},
		"lifespanSec": laser.lifespanSec,
		"damage":      laser.damage,
		"owner":       laser.owner.ID(),
		"collider":    laser.collider.Serialize(),
	}
}
