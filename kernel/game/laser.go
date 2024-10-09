package game

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

func NewLaserProjectile(id int64, position physics.Vector2, rotation float64, owner *Spaceship) *Projectile {
	direction := physics.Vector2{X: math.Cos(rotation), Y: math.Sin(rotation)}

	return &Projectile{
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
	}
}
