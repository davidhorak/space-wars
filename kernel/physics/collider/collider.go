package collider

import "github.com/davidhorak/space-wars/kernel/physics"

type Collider interface {
	Enabled() bool
	SetEnabled(enabled bool)
	Position() physics.Vector2
	SetPosition(position physics.Vector2)
	Rotation() float64
	SetRotation(rotation float64)
	CollidesWith(other Collider) bool
	Serialize() map[string]interface{}
}
