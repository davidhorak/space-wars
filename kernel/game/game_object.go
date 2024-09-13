package game

import (
	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type GameObject interface {
	ID() int64
	Enabled() bool
	SetEnabled(enabled bool)
	Position() physics.Vector2
	SetPosition(position physics.Vector2)
	Update(deltaTimeMs float64, gameManager *GameManager)
	Collider() collider.Collider
	OnCollision(other GameObject, gameManager *GameManager, order int)
	Serialize() map[string]interface{}
}
