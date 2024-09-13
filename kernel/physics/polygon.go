package physics

import (
	"math"
)

type Polygon struct {
	Vertices []Vector2
	edges    []Edge
	minX     float64
	minY     float64
	maxX     float64
	maxY     float64
}

// Rotate rotates the polygon by the given angle (in radians).
// Rotates the polygon around the origin (0, 0).
func (polygon *Polygon) Rotate(angle float64) Polygon {
	rotatedVertices := make([]Vector2, len(polygon.Vertices))
	for i, vertex := range polygon.Vertices {
		rotatedVertices[i].X = vertex.X*math.Cos(angle) - vertex.Y*math.Sin(angle)
		rotatedVertices[i].Y = vertex.X*math.Sin(angle) + vertex.Y*math.Cos(angle)
	}
	return Polygon{Vertices: rotatedVertices}
}

// Translate translates the polygon by the given offset.
func (polygon *Polygon) Translate(offset Vector2) Polygon {
	translatedVertices := make([]Vector2, len(polygon.Vertices))
	for i, vertex := range polygon.Vertices {
		translatedVertices[i].X = vertex.X + offset.X
		translatedVertices[i].Y = vertex.Y + offset.Y
	}
	return Polygon{Vertices: translatedVertices}
}

// Intersects checks if two polygons intersect.
func (polygon *Polygon) Intersects(other Polygon) bool {
	for _, edge := range other.Edges() {
		if polygon.Contains(edge.Start) || polygon.Contains(edge.End) {
			return true
		}
	}
	return false
}

// Contains checks if a point is inside a polygon using ray casting.
// Resources:
//   - https://en.wikipedia.org/wiki/Point_in_polygon#Ray_casting_algorithm
//   - https://people.utm.my/shahabuddin/?p=6277
//
// Simplification:
//   - It assumes that the collision is detected when a vertex enters the polygon,
//     i.e. it does not check for closest point on the polygon edge.
//     Example: ◇ placed directly over □.
//     It could be optimized by first finding the closest point on the polygon edge
//     and then checking if the point is inside the polygon.
//   - It does handle polygons with holes.
func (polygon *Polygon) Contains(point Vector2) bool {
	if len(polygon.Vertices) < 3 {
		return false
	}

	// Check if the point is inside the polygon bounding box
	minX, minY, maxX, maxY := polygon.Bounds()
	if point.X < minX || point.X > maxX || point.Y < minY || point.Y > maxY {
		return false
	}

	intersections := 0
	for _, edge := range polygon.Edges() {
		if ((edge.Start.Y > point.Y) != (edge.End.Y > point.Y)) && (point.X < (edge.End.X-edge.Start.X)*(point.Y-edge.Start.Y)/(edge.End.Y-edge.Start.Y)+edge.Start.X) {
			intersections++
		}
	}

	// odd-in / even-out
	return intersections%2 != 0
}

// Edges returns the edges of the polygon.
func (p *Polygon) Edges() []Edge {
	if p.edges != nil {
		return p.edges
	}

	p.edges = []Edge{}
	for i := 0; i < len(p.Vertices); i++ {
		p.edges = append(p.edges, Edge{Start: p.Vertices[i], End: p.Vertices[(i+1)%len(p.Vertices)]})
	}

	return p.edges
}

// Bounds returns the bounding box of the polygon.
func (p *Polygon) Bounds() (minX, minY, maxX, maxY float64) {
	if p.minX != 0 && p.minY != 0 && p.maxX != 0 && p.maxY != 0 {
		return p.minX, p.minY, p.maxX, p.maxY
	}

	p.minX, p.minY, p.maxX, p.maxY = math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, vertex := range p.Vertices {
		if vertex.X < p.minX {
			p.minX = vertex.X
		}
		if vertex.X > p.maxX {
			p.maxX = vertex.X
		}
		if vertex.Y < p.minY {
			p.minY = vertex.Y
		}
		if vertex.Y > p.maxY {
			p.maxY = vertex.Y
		}
	}

	return p.minX, p.minY, p.maxX, p.maxY
}
