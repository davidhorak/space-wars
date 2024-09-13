package collider

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
)

type CircleCollider struct {
	enabled  bool
	position physics.Vector2
	radius   float64
}

func (circle *CircleCollider) Enabled() bool {
	return circle.enabled
}

func (circle *CircleCollider) SetEnabled(enabled bool) {
	circle.enabled = enabled
}

func (circle *CircleCollider) Position() physics.Vector2 {
	return circle.position
}

func (circle *CircleCollider) SetPosition(position physics.Vector2) {
	circle.position = position
}

func (circle *CircleCollider) Rotation() float64 {
	return 0
}

func (circle *CircleCollider) SetRotation(rotation float64) {}

func NewCircleCollider(position physics.Vector2, radius float64) *CircleCollider {
	return &CircleCollider{
		enabled:  true,
		position: position,
		radius:   radius,
	}
}

func (circle *CircleCollider) CollidesWith(other Collider) bool {
	switch other := other.(type) {
	case *SquareCollider:
		return other.CollidesWith(circle)
	case *CircleCollider:
		return circleCollidesWithCircle(*circle, *other)
	case *PolygonCollider:
		return other.CollidesWith(circle)
	default:
		return false
	}
}

func (circle *CircleCollider) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "circle",
		"enabled": circle.enabled,
		"position": map[string]interface{}{
			"x": circle.position.X,
			"y": circle.position.Y,
		},
		"radius": circle.radius,
	}
}

func circleCollidesWithCircle(circle CircleCollider, other CircleCollider) bool {
	distance := math.Sqrt(math.Pow(circle.position.X-other.position.X, 2) + math.Pow(circle.position.Y-other.position.Y, 2))
	return distance <= circle.radius+other.radius
}
