package game

import (
	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type Explosion struct {
	id          int64
	enabled     bool
	position    physics.Vector2
	radius      float64
	durationSec float64
	lifespanSec float64
}

func NewExplosion(id int64, position physics.Vector2, radius float64, lifespanSec float64) *Explosion {
	return &Explosion{
		id:          id,
		enabled:     true,
		position:    position,
		radius:      radius,
		durationSec: lifespanSec,
		lifespanSec: lifespanSec,
	}
}

func (explosion *Explosion) ID() int64 {
	return explosion.id
}

func (explosion *Explosion) Enabled() bool {
	return explosion.enabled
}

func (explosion *Explosion) SetEnabled(enabled bool) {
	explosion.enabled = enabled
}

func (explosion *Explosion) Position() physics.Vector2 {
	return explosion.position
}

func (explosion *Explosion) SetPosition(position physics.Vector2) {
	explosion.position = position
}

func (explosion *Explosion) Update(deltaTimeMs float64, gameManager *GameManager) {
	explosion.lifespanSec -= deltaTimeMs / 1000
	if explosion.lifespanSec <= 0 {
		explosion.lifespanSec = 0
		explosion.enabled = false
		gameManager.RemoveGameObject(explosion)
	}
}

func (explosion *Explosion) Collider() collider.Collider {
	return nil
}

func (explosion *Explosion) OnCollision(other GameObject, gameManager *GameManager, order int) {}

func (explosion *Explosion) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "explosion",
		"id":      explosion.id,
		"enabled": explosion.enabled,
		"position": map[string]interface{}{
			"x": explosion.position.X,
			"y": explosion.position.Y,
		},
		"radius":      explosion.radius,
		"durationSec": explosion.durationSec,
		"lifespanSec": explosion.lifespanSec,
	}
}
