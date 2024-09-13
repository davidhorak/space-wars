package game

import (
	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type Asteroid struct {
	id       int64
	enabled  bool
	position physics.Vector2
	radius   float64
	collider collider.CircleCollider
}

func NewAsteroid(id int64, position physics.Vector2, radius float64) *Asteroid {
	return &Asteroid{
		id:       id,
		enabled:  true,
		position: position,
		radius:   radius,
		collider: *collider.NewCircleCollider(position, radius),
	}
}

func (asteroid *Asteroid) ID() int64 {
	return asteroid.id
}

func (asteroid *Asteroid) Enabled() bool {
	return asteroid.enabled
}

func (asteroid *Asteroid) SetEnabled(enabled bool) {
	asteroid.enabled = enabled
}

func (asteroid *Asteroid) Position() physics.Vector2 {
	return asteroid.position
}

func (asteroid *Asteroid) SetPosition(position physics.Vector2) {
	asteroid.position = position
}

func (asteroid *Asteroid) Update(deltaTimeMs float64, gameManager *GameManager) {}

func (asteroid *Asteroid) Collider() collider.Collider {
	return &asteroid.collider
}

func (asteroid *Asteroid) OnCollision(other GameObject, gameManager *GameManager, order int) {}

func (asteroid *Asteroid) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "asteroid",
		"id":      asteroid.id,
		"enabled": asteroid.enabled,
		"position": map[string]interface{}{
			"x": asteroid.position.X,
			"y": asteroid.position.Y,
		},
		"radius":   asteroid.radius,
		"collider": asteroid.collider.Serialize(),
	}
}
