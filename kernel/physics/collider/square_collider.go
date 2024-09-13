package collider

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
)

type SquareCollider struct {
	enabled  bool
	position physics.Vector2
	// Rotation in radians
	rotation float64
	size     physics.Size
}

func NewSquareCollider(position physics.Vector2, rotation float64, size physics.Size) *SquareCollider {
	return &SquareCollider{
		enabled:  true,
		position: position,
		rotation: rotation,
		size:     size,
	}
}

func (square *SquareCollider) Enabled() bool {
	return square.enabled
}

func (square *SquareCollider) SetEnabled(enabled bool) {
	square.enabled = enabled
}

func (square *SquareCollider) Position() physics.Vector2 {
	return square.position
}

func (square *SquareCollider) SetPosition(position physics.Vector2) {
	square.position = position
}

func (square *SquareCollider) Rotation() float64 {
	return square.rotation
}

func (square *SquareCollider) SetRotation(rotation float64) {
	square.rotation = rotation
}

func (square *SquareCollider) CollidesWith(other Collider) bool {
	switch other := other.(type) {
	case *SquareCollider:
		return squareCollidesWithSquare(*square, *other)
	case *CircleCollider:
		return squareCollidesWithCircle(*square, *other)
	case *PolygonCollider:
		return squareCollidesWithPolygon(*square, *other)
	default:
		return false
	}
}

func (square *SquareCollider) IsRotated() bool {
	return square.rotation != 0
}

// Absolute returns the polygon with all the transformations applied.
func (square *SquareCollider) Absolute() physics.Polygon {
	res := square.Polygon()
	return res.Translate(square.position)
}

// Polygon returns the polygon of the square.
func (square *SquareCollider) Polygon() physics.Polygon {
	boxWidthHalf := square.size.Width / 2
	boxHeightHalf := square.size.Height / 2

	vertices := []physics.Vector2{
		{X: -boxWidthHalf, Y: -boxHeightHalf},
		{X: boxWidthHalf, Y: -boxHeightHalf},
		{X: boxWidthHalf, Y: boxHeightHalf},
		{X: -boxWidthHalf, Y: boxHeightHalf},
	}

	polygon := physics.Polygon{Vertices: vertices}

	if square.IsRotated() {
		return polygon.Rotate(square.rotation)
	}

	return polygon
}

func (square *SquareCollider) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":    "square",
		"enabled": square.enabled,
		"position": map[string]interface{}{
			"x": square.position.X,
			"y": square.position.Y,
		},
		"rotation": square.rotation,
		"size": map[string]interface{}{
			"width":  square.size.Width,
			"height": square.size.Height,
		},
	}
}

func squareCollidesWithSquare(box1 SquareCollider, box2 SquareCollider) bool {
	// If both boxes are not rotated, use the simpler collision detection
	if !box1.IsRotated() && !box2.IsRotated() {
		boxWidth := box1.size.Width / 2
		boxHeight := box1.size.Height / 2
		otherWidth := box2.size.Width / 2
		otherHeight := box2.size.Height / 2

		if box1.position.X+boxWidth <= box2.position.X-otherWidth ||
			box1.position.X-boxWidth >= box2.position.X+otherWidth ||
			box1.position.Y+boxHeight <= box2.position.Y-otherHeight ||
			box1.position.Y-boxHeight >= box2.position.Y+otherHeight {
			return false
		}
		return true
	}

	// If one of the boxes is rotated, use polygon intersection
	polygon1 := box1.Absolute()
	polygon2 := box2.Absolute()
	return polygon1.Intersects(polygon2)
}

func squareCollidesWithCircle(square SquareCollider, circle CircleCollider) bool {
	circleCenter := circle.position
	circleRadius := circle.radius
	polygon := square.Absolute()
	for _, edge := range polygon.Edges() {
		closestPoint := edge.ClosestPoint(circleCenter)
		distance := math.Sqrt(math.Pow(closestPoint.X-circleCenter.X, 2) + math.Pow(closestPoint.Y-circleCenter.Y, 2))
		if distance < circleRadius {
			return true
		}
	}
	return false
}

func squareCollidesWithPolygon(square SquareCollider, polygon PolygonCollider) bool {
	polygon1 := square.Absolute()
	return polygon1.Intersects(polygon.Absolute())
}
