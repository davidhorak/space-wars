package game

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

func NewRocketProjectile(id int64, position physics.Vector2, rotation float64, owner *Spaceship) *Projectile {
	direction := physics.Vector2{X: math.Cos(rotation), Y: math.Sin(rotation)}

	return &Projectile{
		id:                   id,
		damageType:           DamageTypeRocket,
		enabled:              true,
		position:             position,
		rotation:             rotation,
		velocity:             direction.Multiply(RocketSpeedSec),
		lifespanSec:          RocketLifespanSec,
		damage:               RocketDamage,
		owner:                owner,
		explosionRadius:      float64(RocketExplosionRadius),
		explosionDurationSec: float64(RocketExplosionDurationSec),
		collider: collider.NewCircleCollider(
			position,
			RocketDetonateRadius,
		),
	}
}
