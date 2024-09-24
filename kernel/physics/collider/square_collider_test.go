package collider

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/utils"
)

func TestSquareCollider_Position(t *testing.T) {
	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
	expected := physics.Vector2{X: 0, Y: 0}
	if square.Position() != expected {
		t.Errorf("Expected position to be %v, but got %v", expected, square.Position())
	}
}

func TestSquareCollider_SetPosition(t *testing.T) {
	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
	square.SetPosition(physics.Vector2{X: 1, Y: 1})
	expected := physics.Vector2{X: 1, Y: 1}
	if square.Position() != expected {
		t.Errorf("Expected position to be %v, but got %v", expected, square.Position())
	}
}

func TestSquareCollider_Rotation(t *testing.T) {
	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
	if square.Rotation() != 0 {
		t.Errorf("Expected rotation to be %v, but got %v", 0, square.Rotation())
	}
}

func TestSquareCollider_SetRotation(t *testing.T) {
	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
	square.SetRotation(math.Pi / 2)
	if square.Rotation() != math.Pi/2 {
		t.Errorf("Expected rotation to be %v, but got %v", math.Pi/2, square.Rotation())
	}
}

// func TestSquareCollider_Size(t *testing.T) {
// 	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
// 	if square.Size() != physics.Size{Width: 2, Height: 2} {
// 		t.Errorf("Expected size to be %v, but got %v", physics.Size{Width: 2, Height: 2}, square.Size())
// 	}
// }

// func TestSquareCollider_SetSize(t *testing.T) {
// 	square := NewSquareCollider(physics.Vector2{X: 0, Y: 0}, 0, physics.Size{Width: 2, Height: 2})
// 	square.SetSize(physics.Size{Width: 3, Height: 3})
// 	if square.Size() != physics.Size{Width: 3, Height: 3} {
// 		t.Errorf("Expected size to be %v, but got %v", physics.Size{Width: 3, Height: 3}, square.Size())
// 	}
// }

func TestSquareCollider_IsRotated(t *testing.T) {
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
		square := SquareCollider{
			position: physics.Vector2{X: 0, Y: 0},
			rotation: test.rotation,
			size:     physics.Size{Width: 6, Height: 4},
		}
		result := square.IsRotated()
		if result != test.expected {
			t.Errorf("Expected IsRotated to be %v, but got %v", test.expected, result)
		}
	}
}

func TestSquareCollider_Absolute(t *testing.T) {
	var tests = []struct {
		square   SquareCollider
		expected physics.Polygon
	}{
		{
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 6, Height: 4},
			},
			expected: physics.Polygon{
				Vertices: []physics.Vector2{
					{X: -3, Y: -2},
					{X: 3, Y: -2},
					{X: 3, Y: 2},
					{X: -3, Y: 2},
				},
			},
		},
		{
			square: SquareCollider{
				position: physics.Vector2{X: -3, Y: 2},
				rotation: 0,
				size:     physics.Size{Width: 6, Height: 4},
			},
			expected: physics.Polygon{
				Vertices: []physics.Vector2{
					{X: -6, Y: 0},
					{X: 0, Y: 0},
					{X: 0, Y: 4},
					{X: -6, Y: 4},
				},
			},
		},
		{
			square: SquareCollider{
				position: physics.Vector2{X: -3, Y: 2},
				rotation: math.Pi / 2, // 90 degrees
				size:     physics.Size{Width: 6, Height: 4},
			},
			expected: physics.Polygon{
				Vertices: []physics.Vector2{
					{X: -1, Y: -1},
					{X: -1, Y: 5},
					{X: -5, Y: 5},
					{X: -5, Y: -1},
				},
			},
		},
	}

	for _, test := range tests {
		polygon := test.square.Absolute()
		for i, vertex := range polygon.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected.Vertices[i].X) || !utils.AlmostEqual(vertex.Y, test.expected.Vertices[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected.Vertices[i], vertex)
			}
		}
	}
}

func TestSquareCollider_Polygon(t *testing.T) {
	var tests = []struct {
		square   SquareCollider
		expected []physics.Vector2
	}{
		{
			square: SquareCollider{
				rotation: 0,
				size:     physics.Size{Width: 6, Height: 4},
			},
			expected: []physics.Vector2{
				{X: -3, Y: -2},
				{X: 3, Y: -2},
				{X: 3, Y: 2},
				{X: -3, Y: 2},
			},
		},
		{
			square: SquareCollider{
				rotation: math.Pi / 2, // 90 degrees
				size:     physics.Size{Width: 6, Height: 4},
			},
			expected: []physics.Vector2{
				{X: 2, Y: -3},
				{X: 2, Y: 3},
				{X: -2, Y: 3},
				{X: -2, Y: -3},
			},
		},
		{
			square: SquareCollider{
				rotation: math.Pi / 4, // 45 degrees
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: []physics.Vector2{
				{X: 0, Y: -1.414213562373095},
				{X: 1.414213562373095, Y: 0},
				{X: 0, Y: 1.414213562373095},
				{X: -1.414213562373095, Y: 0},
			},
		},
	}

	for _, test := range tests {
		polygon := test.square.Polygon()

		if len(polygon.Vertices) != 4 {
			t.Errorf("Expected 4 corners, got %d", len(polygon.Vertices))
		}

		for i, vertex := range polygon.Vertices {
			if !utils.AlmostEqual(vertex.X, test.expected[i].X) || !utils.AlmostEqual(vertex.Y, test.expected[i].Y) {
				t.Errorf("Vertex %d: expected %v, got %v", i, test.expected[i], vertex)
			}
		}
	}
}

func TestSquareCollider_CollidesWithSquare(t *testing.T) {
	var tests = []struct {
		description string
		square1     SquareCollider
		square2     SquareCollider
		expected    bool
	}{
		{
			description: "Non-rotated squares colliding #1",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: 1, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Non-rotated squares colliding #2",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: -1, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Non-rotated squares colliding #3",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: 0, Y: 1},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Non-rotated squares colliding #4",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: 0, Y: -1},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Non-rotated non-inclusive collision (touching edges)",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: 2, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: false,
		},
		{
			description: "Rotated squares colliding #1",
			square1: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: math.Pi / 4, // 45 degrees
				size:     physics.Size{Width: 2, Height: 2},
			},
			square2: SquareCollider{
				position: physics.Vector2{X: 1, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result := test.square1.CollidesWith(&test.square2)
		if result != test.expected {
			t.Logf("Test case: %s", test.description)
			t.Errorf("Expected CollidesWith to be %v, but got %v", test.expected, result)
		}
	}
}

func TestSquareCollider_CollidesWithCircle(t *testing.T) {
	var tests = []struct {
		description string
		square      SquareCollider
		circle      CircleCollider
		expected    bool
	}{
		{
			description: "Overlapping circle",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			circle: CircleCollider{
				position: physics.Vector2{X: 0.000000000000001, Y: 0},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "Left of the square",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: 0},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "Left of the square (touching edge)",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -2, Y: 0},
				radius:   1,
			},
			expected: false,
		},
		{
			description: "top left corner of the square",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: -1},
				radius:   1,
			},
			expected: true,
		},
		{
			description: "top left corner of the square (touching edge)",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			circle: CircleCollider{
				position: physics.Vector2{X: -1.70711, Y: -1.70711},
				radius:   1,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.square.CollidesWith(&test.circle)
		if result != test.expected {
			t.Logf("Test case: %s", test.description)
			t.Errorf("Expected CollidesWith to be %v, but got %v", test.expected, result)
		}
	}
}

func TestSquareCollider_CollidesWithPolygon(t *testing.T) {
	var tests = []struct {
		description string
		square      SquareCollider
		polygon     PolygonCollider
		expected    bool
	}{
		{
			description: "Square colliding with a polygon",
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
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
