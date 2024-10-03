package physics

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/utils"
	"github.com/stretchr/testify/assert"
)

func TestPolygon_Rotate(t *testing.T) {
	tests := []struct {
		polygon  Polygon
		angle    float64
		expected Polygon
	}{
		{
			polygon:  Polygon{Vertices: []Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			angle:    0,
			expected: Polygon{Vertices: []Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
		},
		{
			polygon:  Polygon{Vertices: []Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			angle:    math.Pi / 2, // 90 degrees
			expected: Polygon{Vertices: []Vector2{{X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: -1}}},
		},
		{
			polygon:  Polygon{Vertices: []Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			angle:    math.Pi / 4, // 45 degrees
			expected: Polygon{Vertices: []Vector2{{X: 0, Y: -1.414213562373095}, {X: 1.414213562373095, Y: 0}, {X: 0, Y: 1.414213562373095}, {X: -1.414213562373095, Y: 0}}},
		},
	}

	for _, test := range tests {
		polygon := test.polygon.Rotate(test.angle)

		for i, vertex := range polygon.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected.Vertices[i].X) || !utils.AlmostEqual(vertex.Y, test.expected.Vertices[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected.Vertices[i], vertex)
			}
		}
	}
}

func TestPolygon_Translate(t *testing.T) {
	tests := []struct {
		polygon  Polygon
		offset   Vector2
		expected Polygon
	}{
		{
			polygon:  Polygon{Vertices: []Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			offset:   Vector2{X: 2, Y: 2},
			expected: Polygon{Vertices: []Vector2{{X: 1, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 3}, {X: 1, Y: 3}}},
		},
	}

	for _, test := range tests {
		polygon := test.polygon.Translate(test.offset)

		for i, vertex := range polygon.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected.Vertices[i].X) || !utils.AlmostEqual(vertex.Y, test.expected.Vertices[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected.Vertices[i], vertex)
			}
		}
	}
}

func TestPolygon_Intersects(t *testing.T) {
	polygonSquare := Polygon{Vertices: []Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}
	polygonSquareRotated45 := polygonSquare.Rotate(math.Pi / 4)

	tests := []struct {
		description string
		polygon     Polygon
		other       Polygon
		expected    bool
	}{
		{
			description: "Same polygon",
			polygon:     polygonSquare,
			other:       polygonSquare,
			expected:    true,
		},
		{
			description: "Horizontally translated polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: -1, Y: 0}),
			expected:    true,
		},
		{
			description: "Vertically translated polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: 0, Y: -1}),
			expected:    true,
		},
		{
			description: "left of polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: -2.000000000000001, Y: 0}),
			expected:    false,
		},
		{
			description: "right of polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: 2, Y: 0}),
			expected:    false,
		},
		{
			description: "above polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: 0, Y: 2}),
			expected:    false,
		},
		{
			description: "below polygon",
			polygon:     polygonSquare,
			other:       polygonSquare.Translate(Vector2{X: 0, Y: -2.000000000000001}),
			expected:    false,
		},
		{
			description: "45 degrees rotated polygon, translated left",
			polygon:     polygonSquare,
			other:       polygonSquareRotated45.Translate(Vector2{X: -0.5, Y: 0}),
			expected:    true,
		},
		{
			description: "45 degrees rotated polygon, translated right",
			polygon:     polygonSquare,
			other:       polygonSquareRotated45.Translate(Vector2{X: 0.5, Y: 0}),
			expected:    true,
		},
		{
			description: "45 degrees rotated polygon, translated up",
			polygon:     polygonSquare,
			other:       polygonSquareRotated45.Translate(Vector2{X: 0, Y: 0.5}),
			expected:    true,
		},
		{
			description: "45 degrees rotated polygon, translated down",
			polygon:     polygonSquare,
			other:       polygonSquareRotated45.Translate(Vector2{X: 0, Y: -0.5}),
			expected:    true,
		},
	}

	for _, test := range tests {
		got := test.polygon.Intersects(test.other)
		assert.Equal(t, test.expected, got)
	}
}

func TestPolygon_Contains(t *testing.T) {
	polygonSquare := Polygon{Vertices: []Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}
	polygonSquareRotated45 := polygonSquare.Rotate(math.Pi / 4)
	polygonWithInset := Polygon{Vertices: []Vector2{
		{X: -1.5, Y: -1},
		{X: 1.5, Y: -1},
		{X: 1.5, Y: 1},
		{X: 0.5, Y: 1},
		{X: 0.5, Y: 0},
		{X: -0.5, Y: 0},
		{X: -0.5, Y: 1},
		{X: -1.5, Y: 1},
	}}

	tests := []struct {
		description string
		polygon     Polygon
		point       Vector2
		expected    bool
	}{
		// Invalid polygons
		{
			description: "Polygon with less than 3 vertices",
			polygon:     Polygon{Vertices: []Vector2{{X: 0, Y: 0}}},
			point:       Vector2{X: 0, Y: 0},
			expected:    false,
		},
		// Square
		{
			description: "Point inside the polygon #1",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the polygon #2",
			polygon:     polygonSquare,
			point:       Vector2{X: -0.5, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the polygon #3",
			polygon:     polygonSquare,
			point:       Vector2{X: 0.5, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the polygon #4",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: -0.5},
			expected:    true,
		},
		{
			description: "Point inside the polygon #5",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: 0.5},
			expected:    true,
		},
		{
			description: "Top left corner",
			polygon:     polygonSquare,
			// Ray-casting left to right, overlapping left vertex is counted as inside
			point:    Vector2{X: -1, Y: -1},
			expected: true,
		},
		{
			description: "Top left corner (outside)",
			polygon:     polygonSquare,
			point:       Vector2{X: -1.000000000000001, Y: -1},
			expected:    false,
		},
		{
			description: "Top right corner",
			polygon:     polygonSquare,
			point:       Vector2{X: 1, Y: -1},
			expected:    false,
		},
		{
			description: "Bottom right corner",
			polygon:     polygonSquare,
			point:       Vector2{X: 1, Y: 1},
			expected:    false,
		},
		{
			description: "Bottom left corner",
			polygon:     polygonSquare,
			point:       Vector2{X: -1, Y: 1},
			expected:    false,
		},
		{
			description: "left edge",
			polygon:     polygonSquare,
			point:       Vector2{X: -1, Y: 0},
			// Ray-casting left to right, overlapping left vertex is counted as inside
			expected: true,
		},
		{
			description: "left edge (outside)",
			polygon:     polygonSquare,
			point:       Vector2{X: -1.000000000000001, Y: 0},
			expected:    false,
		},
		{
			description: "right edge",
			polygon:     polygonSquare,
			point:       Vector2{X: 1, Y: 0},
			expected:    false,
		},
		{
			description: "top edge",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: 1},
			expected:    false,
		},
		{
			description: "bottom edge",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: -1},
			// Ray-casting left to right, overlapping left vertex is counted as inside
			expected: true,
		},
		{
			description: "bottom edge (outside)",
			polygon:     polygonSquare,
			point:       Vector2{X: 0, Y: -1.000000000000001},
			expected:    false,
		},
		{
			description: "far away",
			polygon:     polygonSquare,
			point:       Vector2{X: -10, Y: -10},
			expected:    false,
		},
		// Rotated square
		{
			description: "Point inside the rotated polygon #1",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 0, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the rotated polygon #2",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: -0.5, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the rotated polygon #3",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 0.5, Y: 0},
			expected:    true,
		},
		{
			description: "Point inside the rotated polygon #4",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 0, Y: -0.5},
			expected:    true,
		},
		{
			description: "Point inside the rotated polygon #5",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 0, Y: 0.5},
			expected:    true,
		},
		{
			description: "Top left corner of the rotated polygon",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: -1, Y: -1},
			expected:    false,
		},
		{
			description: "Top right corner of the rotated polygon",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 1, Y: -1},
			expected:    false,
		},
		{
			description: "Bottom right corner of the rotated polygon",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: 1, Y: 1},
			expected:    false,
		},
		{
			description: "Bottom left corner of the rotated polygon",
			polygon:     polygonSquareRotated45,
			point:       Vector2{X: -1, Y: 1},
			expected:    false,
		},
		// Polygon with inset
		{
			description: "Point inside the polygon with inset",
			polygon:     polygonWithInset,
			point:       Vector2{X: 0, Y: -0.5},
			expected:    true,
		},
		{
			description: "Point in the inset of the polygon with inset",
			polygon:     polygonWithInset,
			point:       Vector2{X: 0, Y: 0.5},
			expected:    false,
		},
	}

	for _, test := range tests {
		got := test.polygon.Contains(test.point)
		assert.Equal(t, test.expected, got)
	}
}

func TestPolygon_Edges(t *testing.T) {
	polygon := Polygon{Vertices: []Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}

	expectedEdges := []Edge{
		{Start: Vector2{X: -1, Y: -1}, End: Vector2{X: 1, Y: -1}},
		{Start: Vector2{X: 1, Y: -1}, End: Vector2{X: 1, Y: 1}},
		{Start: Vector2{X: 1, Y: 1}, End: Vector2{X: -1, Y: 1}},
		{Start: Vector2{X: -1, Y: 1}, End: Vector2{X: -1, Y: -1}},
	}

	for i, edge := range polygon.Edges() {
		assert.Equal(t, expectedEdges[i], edge)
	}

	// Cache edges
	polygon.Edges()
	assert.Equal(t, expectedEdges, polygon.edges)
}

func TestPolygon_Bounds(t *testing.T) {
	polygon := Polygon{Vertices: []Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}

	minX, minY, maxX, maxY := polygon.Bounds()
	assert.Equal(t, -1.0, minX)
	assert.Equal(t, -1.0, minY)
	assert.Equal(t, 1.0, maxX)
	assert.Equal(t, 1.0, maxY)

	// Cached bounds
	polygon.Bounds()
	assert.Equal(t, -1.0, polygon.minX)
	assert.Equal(t, -1.0, polygon.minY)
	assert.Equal(t, 1.0, polygon.maxX)
	assert.Equal(t, 1.0, polygon.maxY)
}
