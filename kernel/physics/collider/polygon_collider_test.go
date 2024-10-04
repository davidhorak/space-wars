package collider

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/utils"
	"github.com/stretchr/testify/assert"
)

func TestPolygonCollider_Enabled(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	assert.True(t, polygon.Enabled())
}

func TestPolygonCollider_SetEnabled(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	polygon.SetEnabled(false)
	assert.False(t, polygon.Enabled())
}

func TestPolygonCollider_Position(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	expected := physics.Vector2{X: 0, Y: 0}
	assert.Equal(t, expected, polygon.Position())
}

func TestPolygonCollider_SetPosition(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	polygon.SetPosition(physics.Vector2{X: 1, Y: 1})
	expected := physics.Vector2{X: 1, Y: 1}
	assert.Equal(t, expected, polygon.Position())
}

func TestPolygonCollider_Rotation(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	assert.Equal(t, 0.0, polygon.Rotation())
}

func TestPolygonCollider_SetRotation(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	polygon.SetRotation(math.Pi / 2)
	assert.Equal(t, math.Pi/2, polygon.Rotation())
}

func TestPolygonCollider_IsRotated(t *testing.T) {
	var tests = []struct {
		rotation float64
		expected bool
	}{
		{
			rotation: 0,
			expected: false,
		},
		{
			rotation: math.Pi / 2, // 90 degrees,
			expected: true,
		},
	}

	for _, test := range tests {
		polygon := PolygonCollider{
			position: physics.Vector2{X: 0, Y: 0},
			rotation: test.rotation,
		}
		result := polygon.IsRotated()
		assert.Equal(t, test.expected, result)
	}
}

func TestPolygonColliderAbsolute(t *testing.T) {
	tests := []struct {
		position physics.Vector2
		polygon  physics.Polygon
		rotation float64
		expected physics.Polygon
	}{
		{
			position: physics.Vector2{X: 0, Y: 0},
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			rotation: 0,
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
		},
		{
			position: physics.Vector2{X: 0, Y: 0},
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			rotation: math.Pi / 4, // 45 degrees
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: 0, Y: -1.414213562373095}, {X: 1.414213562373095, Y: 0}, {X: 0, Y: 1.414213562373095}, {X: -1.414213562373095, Y: 0}}},
		},
		{
			position: physics.Vector2{X: 2, Y: -4},
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			rotation: math.Pi / 4, // 45 degrees
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: 2, Y: -5.414213562373095}, {X: 3.414213562373095, Y: -4}, {X: 2, Y: -2.58578643763}, {X: 0.58578643762, Y: -4}}},
		},
	}

	for _, test := range tests {
		polygon := PolygonCollider{
			position: test.position,
			polygon:  test.polygon,
			rotation: test.rotation,
		}
		result := polygon.Absolute()
		for i, vertex := range result.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected.Vertices[i].X) || !utils.AlmostEqual(vertex.Y, test.expected.Vertices[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected.Vertices[i], vertex)
			}
		}
	}
}

func TestPolygonCollider_Rotated(t *testing.T) {
	tests := []struct {
		position physics.Vector2
		polygon  physics.Polygon
		rotation float64
		expected physics.Polygon
	}{
		{
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: 1, Y: 1}}},
			rotation: 0,
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: 1, Y: 1}}},
		},
		{
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -2, Y: -1}, {X: 2, Y: -1}, {X: 2, Y: 1}, {X: -2, Y: 1}}},
			rotation: math.Pi / 2, // 90 degrees
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: 1, Y: -2}, {X: 1, Y: 2}, {X: -1, Y: 2}, {X: -1, Y: -2}}},
		},
		{
			polygon:  physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}},
			rotation: math.Pi / 4, // 45
			expected: physics.Polygon{Vertices: []physics.Vector2{{X: 0, Y: -1.414213562373095}, {X: 1.414213562373095, Y: 0}, {X: 0, Y: 1.414213562373095}, {X: -1.414213562373095, Y: 0}}},
		},
	}

	for _, test := range tests {
		polygon := PolygonCollider{
			position: test.position,
			polygon:  test.polygon,
			rotation: test.rotation,
		}
		result := polygon.Rotated()

		for i, vertex := range result.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected.Vertices[i].X) || !utils.AlmostEqual(vertex.Y, test.expected.Vertices[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected.Vertices[i], vertex)
			}
		}
	}
}

func TestPolygonCollider_CollidesWithSquare(t *testing.T) {
	var tests = []struct {
		description string
		polygon     PolygonCollider
		square      SquareCollider
		expected    bool
	}{
		{
			description: "Polygon colliding with a square",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon: physics.Polygon{
					Vertices: []physics.Vector2{
						{X: -1, Y: -1},
						{X: 1, Y: -1},
						{X: 1, Y: 1},
						{X: -1, Y: 1},
					},
				},
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result := test.polygon.CollidesWith(&test.square)
		assert.Equal(t, test.expected, result)
	}
}

func TestPolygonCollider_CollidesWithCircle(t *testing.T) {
	polygonSquare := physics.Polygon{Vertices: []physics.Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}

	var tests = []struct {
		description string
		polygon     PolygonCollider
		circle      CircleCollider
		expected    bool
	}{
		{
			description: "Overlapping circle",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			circle: CircleCollider{
				position: physics.Vector2{X: 0.000000000000001, Y: 0},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "Left of the square",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: 0},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "Left of the square (touching edge)",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -2.000000000000001, Y: 0},
				radius:   1,
			},
			expected: false,
		},
		{
			description: "top left corner of the square",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: -1},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "top left corner of the square (touching edge)",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1.70711, Y: -1.70711},
				radius:   1,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.polygon.CollidesWith(&test.circle)
		assert.Equal(t, test.expected, result)
	}
}

func TestPolygonCollider_CollidesWithPolygon(t *testing.T) {
	polygonSquare := physics.Polygon{Vertices: []physics.Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}
	polygonSquareRotated45 := polygonSquare.Rotate(math.Pi / 4)

	tests := []struct {
		description string
		polygon     PolygonCollider
		other       PolygonCollider
		expected    bool
	}{
		{
			description: "Overlapping polygons",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			other: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: true,
		},
		{
			description: "Overlapping with rotated polygon, left side",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			other: PolygonCollider{
				position: physics.Vector2{X: -1, Y: 0},
				rotation: 0,
				polygon:  polygonSquareRotated45,
			},
			expected: true,
		},
		{
			description: "Overlapping with rotated polygon, left side, outside",
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			other: PolygonCollider{
				position: physics.Vector2{X: -2.82842712475, Y: 0},
				rotation: 0,
				polygon:  polygonSquareRotated45,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.polygon.CollidesWith(&test.other)
		assert.Equal(t, test.expected, result)
	}
}

func TestPolygonCollider_CollidesWithOther(t *testing.T) {
	collider_mocked := new(MockCollider)

	polygon_collider := NewPolygonCollider(
		physics.Vector2{X: 0, Y: 0},
		0,
		physics.Polygon{Vertices: []physics.Vector2{
			{X: -1, Y: -1},
			{X: 1, Y: -1},
			{X: 1, Y: 1},
			{X: -1, Y: 1},
		}},
	)

	result := polygon_collider.CollidesWith(collider_mocked)
	assert.False(t, result)
}

func TestPolygonCollider_Serialize(t *testing.T) {
	polygon_collider := NewPolygonCollider(
		physics.Vector2{X: 0, Y: 0},
		0,
		physics.Polygon{Vertices: []physics.Vector2{
			{X: -1, Y: -1},
			{X: 1, Y: -1},
			{X: 1, Y: 1},
			{X: -1, Y: 1},
		}},
	)

	assert.Equal(t, map[string]interface{}{
		"type":    "polygon",
		"enabled": true,
		"position": map[string]interface{}{
			"x": 0.0,
			"y": 0.0,
		},
		"rotation": 0.0,
		"vertices": []map[string]interface{}{
			{
				"x": -1.0,
				"y": -1.0,
			},
			{
				"x": 1.0,
				"y": -1.0,
			},
			{
				"x": 1.0,
				"y": 1.0,
			},
			{
				"x": -1.0,
				"y": 1.0,
			},
		},
	}, polygon_collider.Serialize())
}
