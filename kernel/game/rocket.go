package game

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type RocketProjectile struct {
	Projectile
}

func NewRocketProjectile(id int64, position physics.Vector2, rotation float64, owner *Spaceship) *RocketProjectile {
	direction := physics.Vector2{X: math.Cos(rotation), Y: math.Sin(rotation)}

	return &RocketProjectile{
		Projectile: Projectile{
			id:                   id,
			enabled:              true,
			damageType:           DamageTypeRocket,
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
		},
	}
}

func (rocket *RocketProjectile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "rocket",
		"id":      rocket.id,
		"enabled": rocket.enabled,
		"position": map[string]interface{}{
			"x": rocket.position.X,
			"y": rocket.position.Y,
		},
		"rotation": rocket.rotation,
		"velocity": map[string]interface{}{
			"x": rocket.velocity.X,
			"y": rocket.velocity.Y,
		},
		"lifespanSec": rocket.lifespanSec,
		"damage":      rocket.damage,
		"owner":       rocket.owner.ID(),
		"collider":    rocket.collider.Serialize(),
	}
}
