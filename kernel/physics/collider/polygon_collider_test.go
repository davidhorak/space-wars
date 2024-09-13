package collider

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/utils"
)

func TestPolygonCollider_Position(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	expected := physics.Vector2{X: 0, Y: 0}
	if polygon.Position() != expected {
		t.Errorf("Expected position to be %v, but got %v", expected, polygon.Position())
	}
}

func TestPolygonCollider_SetPosition(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	polygon.SetPosition(physics.Vector2{X: 1, Y: 1})
	expected := physics.Vector2{X: 1, Y: 1}
	if polygon.Position() != expected {
		t.Errorf("Expected position to be %v, but got %v", expected, polygon.Position())
	}
}

func TestPolygonCollider_Rotation(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	if polygon.Rotation() != 0 {
		t.Errorf("Expected rotation to be %v, but got %v", 0, polygon.Rotation())
	}
}

func TestPolygonCollider_SetRotation(t *testing.T) {
	polygon := NewPolygonCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Polygon{Vertices: []physics.Vector2{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}})
	polygon.SetRotation(math.Pi / 2)
	if polygon.Rotation() != math.Pi/2 {
		t.Errorf("Expected rotation to be %v, but got %v", math.Pi/2, polygon.Rotation())
	}
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
		if result != test.expected {
			t.Errorf("Expected IsRotated to be %v, but got %v", test.expected, result)
		}
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
		result := test.square.CollidesWith(&test.polygon)
		if result != test.expected {
			t.Logf("Test case: %s", test.description)
			t.Errorf("Expected CollidesWith to be %v, but got %v", test.expected, result)
		}
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
		if result != test.expected {
			t.Logf("Test case: %s", test.description)
			t.Errorf("Expected CollidesWith to be %v, but got %v", test.expected, result)
		}
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
		if result != test.expected {
			t.Logf("Test case: %s", test.description)
			t.Errorf("Expected CollidesWith to be %v, but got %v", test.expected, result)
		}
	}
}
