package collider

import (
	"math"

	"github.com/davidhorak/space-wars/kernel/physics"
)

type PolygonCollider struct {
	enabled  bool
	position physics.Vector2
	// Rotation in radians
	rotation float64
	// Polygon is clockwise ordered vertices, relative to the position, origin (0, 0).
	polygon physics.Polygon
}

func NewPolygonCollider(position physics.Vector2, rotation float64, polygon physics.Polygon) *PolygonCollider {
	return &PolygonCollider{
		enabled:  true,
		position: position,
		rotation: rotation,
		polygon:  polygon,
	}
}

func (polygon *PolygonCollider) Enabled() bool {
	return polygon.enabled
}

func (polygon *PolygonCollider) SetEnabled(enabled bool) {
	polygon.enabled = enabled
}

func (polygon *PolygonCollider) Position() physics.Vector2 {
	return polygon.position
}

func (polygon *PolygonCollider) SetPosition(position physics.Vector2) {
	polygon.position = position
}

func (polygon *PolygonCollider) Rotation() float64 {
	return polygon.rotation
}

func (polygon *PolygonCollider) SetRotation(rotation float64) {
	polygon.rotation = rotation
}

func (polygon *PolygonCollider) CollidesWith(other Collider) bool {
	switch other := other.(type) {
	case *SquareCollider:
		return other.CollidesWith(polygon)
	case *CircleCollider:
		return polygonCollidesWithCircle(*polygon, *other)
	case *PolygonCollider:
		return polygonCollidesWithPolygon(*polygon, *other)
	default:
		return false
	}
}

func (polygon *PolygonCollider) IsRotated() bool {
	return polygon.rotation != 0
}

// Absolute returns the polygon with all the transformations applied.
func (polygon *PolygonCollider) Absolute() physics.Polygon {
	res := polygon.polygon
	if polygon.IsRotated() {
		res = polygon.Rotated()
	}
	return res.Translate(polygon.position)
}

// Rotated returns the the polygon rotated by the rotation angle.
// The rotation is relative to the origin (0, 0).
func (polygon *PolygonCollider) Rotated() physics.Polygon {
	rotatedVertices := make([]physics.Vector2, len(polygon.polygon.Vertices))
	cos := math.Cos(polygon.rotation)
	sin := math.Sin(polygon.rotation)

	for i, vertex := range polygon.polygon.Vertices {
		rotatedVertices[i].X = vertex.X*cos - vertex.Y*sin
		rotatedVertices[i].Y = vertex.X*sin + vertex.Y*cos
	}

	return physics.Polygon{Vertices: rotatedVertices}
}

func (polygon *PolygonCollider) Serialize() map[string]interface{} {
	vertices := make([]map[string]interface{}, len(polygon.polygon.Vertices))
	for i, vertex := range polygon.polygon.Vertices {
		vertices[i] = map[string]interface{}{
			"x": vertex.X,
			"y": vertex.Y,
		}
	}

	return map[string]interface{}{
		"type":    "polygon",
		"enabled": polygon.enabled,
		"position": map[string]interface{}{
			"x": polygon.position.X,
			"y": polygon.position.Y,
		},
		"rotation": polygon.rotation,
		"vertices": vertices,
	}
}

func polygonCollidesWithCircle(polygon PolygonCollider, circle CircleCollider) bool {
	circleCenter := circle.position
	circleRadius := circle.radius
	abs := polygon.Absolute()

	for _, edge := range abs.Edges() {
		closestPoint := edge.ClosestPoint(circleCenter)
		distance := math.Sqrt(math.Pow(closestPoint.X-circleCenter.X, 2) + math.Pow(closestPoint.Y-circleCenter.Y, 2))
		if distance <= circleRadius {
			return true
		}
	}
	return false
}

func polygonCollidesWithPolygon(polygon PolygonCollider, other PolygonCollider) bool {
	a := polygon.Absolute()
	b := other.Absolute()
	return a.Intersects(b)
}
