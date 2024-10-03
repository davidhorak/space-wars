package collider

import (
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/stretchr/testify/assert"
)

func TestCircleCollider_Enabled(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	assert.True(t, circle_collider.Enabled())
}

func TestCircleCollider_SetEnabled(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	circle_collider.SetEnabled(false)
	assert.False(t, circle_collider.Enabled())
}

func TestCircleCollider_Position(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	expected := physics.Vector2{X: 0, Y: 0}
	assert.Equal(t, expected, circle_collider.Position())
}

func TestCircleCollider_SetRotation(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	circle_collider.SetPosition(physics.Vector2{X: 1, Y: 1})
	expected := physics.Vector2{X: 1, Y: 1}
	assert.Equal(t, expected, circle_collider.Position())
}

func TestCircleCollider_SetPosition(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	circle_collider.SetPosition(physics.Vector2{X: 1, Y: 1})
	expected := physics.Vector2{X: 1, Y: 1}
	assert.Equal(t, expected, circle_collider.Position())
}

func TestCircleCollider_Rotation(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	expected := 0.0
	assert.Equal(t, expected, circle_collider.Rotation())
}

func TestCircleCollider_Radius(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 0, Y: 0}, 1)
	expected := 1.0
	assert.Equal(t, expected, circle_collider.Radius())
}

func TestCircleCollider_CollidesWithSquare(t *testing.T) {
	var tests = []struct {
		description string
		circle      CircleCollider
		square      SquareCollider
		expected    bool
	}{
		{
			description: "Overlapping circle",
			circle: CircleCollider{
				position: physics.Vector2{X: 0.000000000000001, Y: 0},
				radius:   1,
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Left of the square",
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: 0},
				radius:   1,
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "Left of the square (touching edge)",
			circle: CircleCollider{
				position: physics.Vector2{X: -2, Y: 0},
				radius:   1,
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: false,
		},
		{
			description: "top left corner of the square",
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: -1},
				radius:   1,
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: true,
		},
		{
			description: "top left corner of the square (touching edge)",
			circle: CircleCollider{
				position: physics.Vector2{X: -1.70711, Y: -1.70711},
				radius:   1,
			},
			square: SquareCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				size:     physics.Size{Width: 2, Height: 2},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.circle.CollidesWith(&test.square)
		assert.Equal(t, test.expected, result)
	}
}

func TestCircleCollider_CollidesWithCircle(t *testing.T) {
	tests := []struct {
		name     string
		circle1  CircleCollider
		circle2  CircleCollider
		expected bool
	}{
		{
			name:     "Identical circles",
			circle1:  CircleCollider{position: physics.Vector2{X: 0, Y: 0}, radius: 1},
			circle2:  CircleCollider{position: physics.Vector2{X: 0, Y: 0}, radius: 1},
			expected: true,
		},
		{
			name:     "Overlapping circles",
			circle1:  CircleCollider{position: physics.Vector2{X: 0, Y: 0}, radius: 1},
			circle2:  CircleCollider{position: physics.Vector2{X: 1, Y: 0}, radius: 1},
			expected: true,
		},
		{
			name:     "Touching circles",
			circle1:  CircleCollider{position: physics.Vector2{X: 0, Y: 0}, radius: 1},
			circle2:  CircleCollider{position: physics.Vector2{X: 2, Y: 0}, radius: 1},
			expected: true,
		},
		{
			name:     "Non-colliding circles",
			circle1:  CircleCollider{position: physics.Vector2{X: 0, Y: 0}, radius: 1},
			circle2:  CircleCollider{position: physics.Vector2{X: 3, Y: 0}, radius: 1},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.circle1.CollidesWith(&test.circle2)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCircleCollider_CollidesWithPolygon(t *testing.T) {
	polygonSquare := physics.Polygon{Vertices: []physics.Vector2{
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
	}}

	var tests = []struct {
		description string
		circle      CircleCollider
		polygon     PolygonCollider
		expected    bool
	}{
		{
			description: "Overlapping circle",
			circle: CircleCollider{
				position: physics.Vector2{X: 0.000000000000001, Y: 0},
				radius:   1,
			},
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: true,
		},
		{
			description: "Left of the square",
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: 0},
				radius:   1,
			},
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: true,
		},
		{
			description: "Left of the square (touching edge)",
			circle: CircleCollider{
				position: physics.Vector2{X: -2.000000000000001, Y: 0},
				radius:   1,
			},
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: false,
		},
		{
			description: "top left corner of the square",
			circle: CircleCollider{
				position: physics.Vector2{X: -1, Y: -1},
				radius:   1,
			},
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: true,
		},
		{
			description: "top left corner of the square (touching edge)",
			circle: CircleCollider{
				position: physics.Vector2{X: -1.70711, Y: -1.70711},
				radius:   1,
			},
			polygon: PolygonCollider{
				position: physics.Vector2{X: 0, Y: 0},
				rotation: 0,
				polygon:  polygonSquare,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.circle.CollidesWith(&test.polygon)
		assert.Equal(t, test.expected, result)
	}
}

func TestCircleCollider_CollidesWithOther(t *testing.T) {
	collider_mocked := new(MockCollider)

	circle_collider := NewCircleCollider(
		physics.Vector2{X: 0, Y: 0},
		1,
	)

	result := circle_collider.CollidesWith(collider_mocked)
	assert.False(t, result)
}

func TestCircleCollider_Serialize(t *testing.T) {
	circle_collider := NewCircleCollider(physics.Vector2{X: 5, Y: 10}, 1)

	assert.Equal(t, map[string]interface{}{
		"type":    "circle",
		"enabled": true,
		"position": map[string]interface{}{
			"x": 5.0,
			"y": 10.0,
		},
		"radius": 1.0,
	}, circle_collider.Serialize())
}
